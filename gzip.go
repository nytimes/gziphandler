package gziphandler

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	vary            = "Vary"
	acceptEncoding  = "Accept-Encoding"
	contentEncoding = "Content-Encoding"
	contentType     = "Content-Type"
	contentLength   = "Content-Length"
)

type codings map[string]float64

const (
	// DefaultQValue is the default qvalue to assign to an encoding if no explicit qvalue is set.
	// This is actually kind of ambiguous in RFC 2616, so hopefully it's correct.
	// The examples seem to indicate that it is.
	DefaultQValue = 1.0

	// DefaultMinSize defines the minimum size to reach to enable compression.
	// It's 512 bytes.
	DefaultMinSize = 512
)

// gzipWriterPools stores a sync.Pool for each compression level for reuse of
// gzip.Writers. Use poolIndex to covert a compression level to an index into
// gzipWriterPools.
var gzipWriterPools [gzip.BestCompression - gzip.BestSpeed + 2]*sync.Pool

func init() {
	for i := gzip.BestSpeed; i <= gzip.BestCompression; i++ {
		addLevelPool(i)
	}
	addLevelPool(gzip.DefaultCompression)
}

// poolIndex maps a compression level to its index into gzipWriterPools. It
// assumes that level is a valid gzip compression level.
func poolIndex(level int) int {
	// gzip.DefaultCompression == -1, so we need to treat it special.
	if level == gzip.DefaultCompression {
		return gzip.BestCompression - gzip.BestSpeed + 1
	}
	return level - gzip.BestSpeed
}

func addLevelPool(level int) {
	gzipWriterPools[poolIndex(level)] = &sync.Pool{
		New: func() interface{} {
			// NewWriterLevel only returns error on a bad level, we are guaranteeing
			// that this will be a valid level so it is okay to ignore the returned
			// error.
			w, _ := gzip.NewWriterLevel(nil, level)
			return w
		},
	}
}

// GzipResponseWriter provides an http.ResponseWriter interface, which gzips
// bytes before writing them to the underlying response. This doesn't close the
// writers, so don't forget to do that.
// It can be configured to skip response smaller than minSize.
type GzipResponseWriter struct {
	http.ResponseWriter
	index int // Index for gzipWriterPools.
	gw    *gzip.Writer

	code int // Saves the WriteHeader value.

	minSize      int    // Specifed the minimum response size to gzip. Only the first Write call is checked.
	buff         []byte // Holds the first part of the write before the reaching the minSize or the end of the write.
	bytesWritten int    // Keep trace of the numbers of bytes written.
}

// Write appends data to the gzip writer.
func (w *GzipResponseWriter) Write(b []byte) (int, error) {
	// Update the numbers of bytes written.
	w.bytesWritten += len(b)

	if _, ok := w.ResponseWriter.Header()[contentType]; !ok {
		// If content type is not set, infer it from the uncompressed body.
		w.ResponseWriter.Header().Set(contentType, http.DetectContentType(b))
	}

	// If the global write is bigger than the minSize, compression is used.
	// Otherwise it save the write into a buffer and send it with the regular
	// ResponseWriter at close.
	if w.bytesWritten > w.minSize {
		w.ResponseWriter.Header().Set(contentEncoding, "gzip")

		// if the Content-Length is already set, then calls to Write on gzip
		// will fail to set the Content-Length header since its already set
		// See: https://github.com/golang/go/issues/14975
		w.ResponseWriter.Header().Del(contentLength)

		// Write the header and contentEncoding to gzip.
		w.writeHeader()

		n := 0
		// If the buffer is not nil, it needs to be empty into the gzip writer.
		if w.buff != nil {
			if w.gw == nil {
				w.init()
			}

			n1, err := w.gw.Write(w.buff)
			if err != nil {
				return n, err
			}

			// Saved the buffed length.
			n = n1
		}

		w.buff = nil

		if w.gw == nil {
			w.init()
		}

		n1, err := w.gw.Write(b)
		if err != nil {
			return n, err
		}

		// Add the buffered length to the new write length.
		n += n1
		return n, err
	}

	// When arrived here, it's to save the writes into the buffer until it
	// goes over the minSize or it closes.
	w.buff = append(w.buff, b...)
	return 0, nil
}

// WriteHeader just saves the response code until close or GZIP effective writes.
func (w *GzipResponseWriter) WriteHeader(code int) {
	w.code = code
}

// writeHeader uses the saved code to send it to the ResponseWriter.
func (w *GzipResponseWriter) writeHeader() {
	// debug.PrintStack()
	if w.code == 0 {
		w.code = http.StatusOK
	}
	w.ResponseWriter.WriteHeader(w.code)
}

// init graps a new gzip writer from the gzipWriterPool and writes the correct
// content encoding header.
func (w *GzipResponseWriter) init() {
	// Bytes written during ServeHTTP are redirected to this gzip writer
	// before being written to the underlying response.
	gzw := gzipWriterPools[w.index].Get().(*gzip.Writer)
	gzw.Reset(w.ResponseWriter)
	w.gw = gzw
}

// Close will close the gzip.Writer and will put it back in the gzipWriterPool.
func (w *GzipResponseWriter) Close() error {
	// Buffer not nil means the regular response it returned.
	if w.buff != nil {
		w.writeHeader()
		w.ResponseWriter.Write(w.buff)
		return nil
	}

	if w.gw == nil {
		return nil
	}

	err := w.gw.Close()
	gzipWriterPools[w.index].Put(w.gw)
	return err
}

// Flush flushes the underlying *gzip.Writer and then the underlying
// http.ResponseWriter if it is an http.Flusher. This makes GzipResponseWriter
// an http.Flusher.
func (w *GzipResponseWriter) Flush() {
	if w.gw != nil {
		w.gw.Flush()
	}

	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}

// Hijack implements http.Hijacker. If the underlying ResponseWriter is a
// Hijacker, its Hijack method is returned. Otherwise an error is returned.
func (w *GzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("http.Hijacker interface is not supported")
}

// verify Hijacker interface implementation
var _ http.Hijacker = &GzipResponseWriter{}

// MustNewGzipLevelHandler behaves just like NewGzipLevelHandler except that in
// an error case it panics rather than returning an error.
func MustNewGzipLevelHandler(level int) func(http.Handler) http.Handler {
	wrap, err := NewGzipLevelHandler(level)
	if err != nil {
		panic(err)
	}
	return wrap
}

// NewGzipLevelHandler returns a wrapper function (often known as middleware)
// which can be used to wrap an HTTP handler to transparently gzip the response
// body if the client supports it (via the Accept-Encoding header). Responses will
// be encoded at the given gzip compression level. An error will be returned only
// if an invalid gzip compression level is given, so if one can ensure the level
// is valid, the returned error can be safely ignored.
func NewGzipLevelHandler(level int) (func(http.Handler) http.Handler, error) {
	return NewGzipLevelAndMinSize(level, DefaultMinSize)
}

// NewGzipLevelAndMinSize behave as NewGzipLevelHandler except it let the caller
// specify the minimum size before compression.
func NewGzipLevelAndMinSize(level, askedMinSize int) (func(http.Handler) http.Handler, error) {
	if level != gzip.DefaultCompression && (level < gzip.BestSpeed || level > gzip.BestCompression) {
		return nil, fmt.Errorf("invalid compression level requested: %d", level)
	}
	if askedMinSize < 0 {
		return nil, fmt.Errorf("minimum size must be more than zero")
	}
	return func(h http.Handler) http.Handler {
		index := poolIndex(level)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(vary, acceptEncoding)

			if acceptsGzip(r) {
				gw := &GzipResponseWriter{
					ResponseWriter: w,
					index:          index,
					minSize:        askedMinSize,

					buff: []byte{},
				}
				defer gw.Close()

				h.ServeHTTP(gw, r)
			} else {
				h.ServeHTTP(w, r)
			}
		})
	}, nil
}

// GzipHandler wraps an HTTP handler, to transparently gzip the response body if
// the client supports it (via the Accept-Encoding header). This will compress at
// the default compression level.
func GzipHandler(h http.Handler) http.Handler {
	wrapper, _ := NewGzipLevelHandler(gzip.DefaultCompression)
	return wrapper(h)
}

// acceptsGzip returns true if the given HTTP request indicates that it will
// accept a gzipped response.
func acceptsGzip(r *http.Request) bool {
	acceptedEncodings, _ := parseEncodings(r.Header.Get(acceptEncoding))
	return acceptedEncodings["gzip"] > 0.0
}

// parseEncodings attempts to parse a list of codings, per RFC 2616, as might
// appear in an Accept-Encoding header. It returns a map of content-codings to
// quality values, and an error containing the errors encountered. It's probably
// safe to ignore those, because silently ignoring errors is how the internet
// works.
//
// See: http://tools.ietf.org/html/rfc2616#section-14.3.
func parseEncodings(s string) (codings, error) {
	c := make(codings)
	var e []string

	for _, ss := range strings.Split(s, ",") {
		coding, qvalue, err := parseCoding(ss)

		if err != nil {
			e = append(e, err.Error())
		} else {
			c[coding] = qvalue
		}
	}

	// TODO (adammck): Use a proper multi-error struct, so the individual errors
	//                 can be extracted if anyone cares.
	if len(e) > 0 {
		return c, fmt.Errorf("errors while parsing encodings: %s", strings.Join(e, ", "))
	}

	return c, nil
}

// parseCoding parses a single conding (content-coding with an optional qvalue),
// as might appear in an Accept-Encoding header. It attempts to forgive minor
// formatting errors.
func parseCoding(s string) (coding string, qvalue float64, err error) {
	for n, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		qvalue = DefaultQValue

		if n == 0 {
			coding = strings.ToLower(part)
		} else if strings.HasPrefix(part, "q=") {
			qvalue, err = strconv.ParseFloat(strings.TrimPrefix(part, "q="), 64)

			if qvalue < 0.0 {
				qvalue = 0.0
			} else if qvalue > 1.0 {
				qvalue = 1.0
			}
		}
	}

	if coding == "" {
		err = fmt.Errorf("empty content-coding")
	}

	return
}
