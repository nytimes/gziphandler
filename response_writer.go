package gziphandler

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/CAFxX/gziphandler/custom"
)

// gzipResponseWriter provides an http.ResponseWriter interface, which gzips
// bytes before writing them to the underlying response. This doesn't close the
// writers, so don't forget to do that.
// It can be configured to skip response smaller than minSize.
type gzipResponseWriter struct {
	http.ResponseWriter

	config config
	accept acceptsType

	w          io.WriteCloser
	writerType writerType
	code       int    // Saves the WriteHeader value.
	buf        []byte // Holds the first part of the write before reaching the minSize or the end of the write.
}

var (
	_ io.WriteCloser = &gzipResponseWriter{}
	_ http.Flusher   = &gzipResponseWriter{}
	_ http.Hijacker  = &gzipResponseWriter{}
)

type writerType byte

const (
	noWriter writerType = iota
	passthroughWriter
	gzipWriter
	brotliWriter
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

// Write appends data to the gzip writer.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.w != nil {
		// The responseWriter is already initialized: use it.
		return w.w.Write(b)
	}

	// Save the write into a buffer for later use in GZIP responseWriter (if content is long enough) or at close with regular responseWriter.
	// On the first write, w.buf changes from nil to a valid slice
	w.buf = append(w.buf, b...)

	var (
		cl, _ = strconv.Atoi(w.Header().Get(contentLength))
		ct    = w.Header().Get(contentType)
		ce    = w.Header().Get(contentEncoding)
	)
	// Only continue if they didn't already choose an encoding or a known unhandled content length or type.
	if ce == "" && (cl == 0 || cl >= w.config.minSize) && (ct == "" || handleContentTypes(w.config.contentTypes, w.config.contentTypes, w.config.blacklist, ct, w.config.prefer, w.accept) != handleNone) {
		// If the current buffer is less than minSize and a Content-Length isn't set, then wait until we have more data.
		if len(w.buf) < w.config.minSize && cl == 0 {
			return len(b), nil
		}
		// If the Content-Length is larger than minSize or the current buffer is larger than minSize, then continue.
		if cl >= w.config.minSize || len(w.buf) >= w.config.minSize {
			// If a Content-Type wasn't specified, infer it from the current buffer.
			if ct == "" {
				ct = http.DetectContentType(w.buf)
				w.Header().Set(contentType, ct)
			}
			handle := handleContentTypes(w.config.contentTypes, w.config.contentTypes, w.config.blacklist, ct, w.config.prefer, w.accept)
			if handle == handleGzip || handle == handleBrotli {
				if err := w.startCompress(handle); err != nil {
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
func (w *gzipResponseWriter) startCompress(h handleType) error {
	if h == handleGzip {
		w.Header().Set(contentEncoding, "gzip")
	} else if h == handleBrotli {
		w.Header().Set(contentEncoding, "br")
	} else {
		panic("unknown handleType")
	}

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
		if h == handleGzip {
			w.w = w.config.gzipComp.Get(w.ResponseWriter, w.config.gzLevel)
			w.writerType = gzipWriter
		} else if h == handleBrotli {
			w.w = w.config.brotliComp.Get(w.ResponseWriter, w.config.brLevel)
			w.writerType = brotliWriter
		}

		n, err := w.w.Write(w.buf)

		// This should never happen (per io.Writer docs), but if the write didn't
		// accept the entire buffer but returned no specific error, we have no clue
		// what's going on, so abort just to be safe.
		if err == nil && n < len(w.buf) {
			err = io.ErrShortWrite
		}
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
	w.writerType = passthroughWriter
	// If Write was never called then don't call Write on the underlying ResponseWriter.
	if w.buf == nil {
		return nil
	}
	n, err := w.ResponseWriter.Write(w.buf)
	w.buf = nil
	// This should never happen (per io.Writer docs), but if the write didn't
	// accept the entire buffer but returned no specific error, we have no clue
	// what's going on, so abort just to be safe.
	if err == nil && n < len(w.buf) {
		err = io.ErrShortWrite
	}
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
	switch w.writerType {
	case passthroughWriter:
		return nil
	case gzipWriter, brotliWriter:
		if w.w == nil {
			return nil
		}
		gw := w.w
		w.w = nil
		return gw.Close()
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
		// Only flush once startGzip, startBrotli or startPlain has been called.
		//
		// Flush is thus a no-op until we're certain whether a plain
		// or compressed response will be served.
		return
	}

	switch w.writerType {
	case gzipWriter, brotliWriter:
		if fw, _ := w.w.(custom.Flusher); fw != nil {
			fw.Flush()
		}
	}

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
