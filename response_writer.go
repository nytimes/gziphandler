package gziphandler

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// gzipResponseWriter provides an http.ResponseWriter interface, which gzips
// bytes before writing them to the underlying response. This doesn't close the
// writers, so don't forget to do that.
// It can be configured to skip response smaller than minSize.
type gzipResponseWriter struct {
	http.ResponseWriter

	config config
	accept codings
	common []string
	pool   *sync.Pool

	w    io.Writer
	enc  string
	code int    // Saves the WriteHeader value.
	buf  []byte // Holds the first part of the write before reaching the minSize or the end of the write.
}

var (
	_ io.WriteCloser = &gzipResponseWriter{}
	_ http.Flusher   = &gzipResponseWriter{}
	_ http.Hijacker  = &gzipResponseWriter{}
)

type gzipResponseWriterWithCloseNotify struct {
	*gzipResponseWriter
}

func (w gzipResponseWriterWithCloseNotify) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

var (
	_ io.WriteCloser = gzipResponseWriterWithCloseNotify{}
	_ http.Flusher   = gzipResponseWriterWithCloseNotify{}
	_ http.Hijacker  = gzipResponseWriterWithCloseNotify{}
)

const maxBuf = 1 << 16 // maximum size of recycled buffer

// Write appends data to the gzip writer.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.w != nil {
		// The responseWriter is already initialized: use it.
		return w.w.Write(b)
	}

	// Save the write into a buffer for later use in GZIP responseWriter (if content is long enough) or at close with regular responseWriter.
	// On the first write, w.buf changes from nil to a valid slice
	if w.buf == nil {
		w.buf, _ = w.pool.Get().([]byte)
	}
	w.buf = append(w.buf, b...)

	var (
		cl, _ = strconv.Atoi(w.Header().Get(contentLength))
		ct    = w.Header().Get(contentType)
		ce    = w.Header().Get(contentEncoding)
	)
	// Only continue if they didn't already choose an encoding or a known unhandled content length or type.
	if ce == "" && (cl == 0 || cl >= w.config.minSize) && (ct == "" || handleContentType(ct, w.config.contentTypes, w.config.blacklist)) {
		// If the current buffer is less than minSize and a Content-Length isn't set, then wait until we have more data.
		if len(w.buf) < w.config.minSize && cl == 0 {
			return len(b), nil
		}
		// If the Content-Length is larger than minSize or the current buffer is larger than minSize, then continue.
		if cl >= w.config.minSize || len(w.buf) >= w.config.minSize {
			// If a Content-Type wasn't specified, infer it from the current buffer.
			if ct == "" {
				ct = http.DetectContentType(w.buf)
				if ct != "" {
					// net/http by default performs content sniffing but this is disabled if content-encoding is set.
					// Since we set content-encoding, if content-type was not set and we successfully sniffed it,
					// set the content-type.
					w.Header().Set(contentType, ct)
				}
			}
			if handleContentType(ct, w.config.contentTypes, w.config.blacklist) {
				enc := preferredEncoding(w.accept, w.config.compressor, w.common, w.config.prefer)
				if err := w.startCompress(enc); err != nil {
					return 0, err
				}
				return len(b), nil
			}
		}
	}
	// If we got here, we should not GZIP this response.
	if err := w.startPlain(); err != nil {
		return 0, err
	}
	return len(b), nil
}

// startCompress initializes a compressing writer and writes the buffer.
func (w *gzipResponseWriter) startCompress(enc string) error {
	comp, ok := w.config.compressor[enc]
	if !ok {
		panic("unknown compressor")
	}

	w.Header().Set(contentEncoding, enc)

	// if the Content-Length is already set, then calls to Write on gzip
	// will fail to set the Content-Length header since its already set
	// See: https://github.com/golang/go/issues/14975.
	w.Header().Del(contentLength)

	// Write the header to gzip response.
	if w.code != 0 {
		w.ResponseWriter.WriteHeader(w.code)
		// Ensure that no other WriteHeader's happen
		w.code = 0
	}

	// Initialize and flush the buffer into the gzip response if there are any bytes.
	// If there aren't any, we shouldn't initialize it yet because on Close it will
	// write the gzip header even if nothing was ever written.
	if len(w.buf) > 0 {
		w.w = comp.comp.Get(w.ResponseWriter)
		w.enc = enc

		n, err := w.w.Write(w.buf)

		// This should never happen (per io.Writer docs), but if the write didn't
		// accept the entire buffer but returned no specific error, we have no clue
		// what's going on, so abort just to be safe.
		if err == nil && n < len(w.buf) {
			err = io.ErrShortWrite
		}
		w.recycleBuffer()
		return err
	}
	return nil
}

// startPlain writes to sent bytes and buffer the underlying ResponseWriter without gzip.
func (w *gzipResponseWriter) startPlain() error {
	if w.code != 0 {
		w.ResponseWriter.WriteHeader(w.code)
		// Ensure that no other WriteHeader's happen
		w.code = 0
	}
	w.w = w.ResponseWriter
	w.enc = ""
	// If Write was never called then don't call Write on the underlying ResponseWriter.
	if w.buf == nil {
		return nil
	}
	n, err := w.ResponseWriter.Write(w.buf)
	// This should never happen (per io.Writer docs), but if the write didn't
	// accept the entire buffer but returned no specific error, we have no clue
	// what's going on, so abort just to be safe.
	if err == nil && n < len(w.buf) {
		err = io.ErrShortWrite
	}
	w.recycleBuffer()
	return err
}

// WriteHeader just saves the response code until close or GZIP effective writes.
func (w *gzipResponseWriter) WriteHeader(code int) {
	if w.code == 0 {
		w.code = code
	}
}

// Close will close the gzip.Writer and will put it back in the gzipWriterPool.
func (w *gzipResponseWriter) Close() error {
	if cw, ok := w.w.(io.Closer); ok {
		w.w = nil
		return cw.Close()
	}

	// compression not triggered yet, write out regular response.
	err := w.startPlain()
	// Returns the error if any at write.
	if err != nil {
		err = fmt.Errorf("gziphandler: write to regular responseWriter at close gets error: %v", err)
	}
	return err
}

// Flush flushes the underlying *gzip.Writer and then the underlying
// http.ResponseWriter if it is an http.Flusher. This makes gzipResponseWriter
// an http.Flusher.
func (w *gzipResponseWriter) Flush() {
	if w.w == nil {
		// Flush is thus a no-op until we're certain whether a plain
		// or compressed response will be served.
		return
	}

	// Flush the compressor, if supported,
	// note: http.ResponseWriter does not implement Flusher, so we need to call ResponseWriter.Flush anyway.
	if fw, ok := w.w.(Flusher); ok {
		_ = fw.Flush()
	}

	// Flush the ResponseWriter (the previous Flusher is not expected to flush the parent writer).
	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}

// Hijack implements http.Hijacker. If the underlying ResponseWriter is a
// Hijacker, its Hijack method is returned. Otherwise an error is returned.
func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("http.Hijacker interface is not supported")
}

func (w *gzipResponseWriter) recycleBuffer() {
	if cap(w.buf) > 0 && cap(w.buf) <= maxBuf {
		w.pool.Put(w.buf[:0])
	}
	w.buf = nil
}
