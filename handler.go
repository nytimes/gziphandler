package gziphandler // import "github.com/CAFxX/gziphandler"

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/andybalholm/brotli"
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
	// DefaultMinSize is the default minimum size for which we enable compression.
	// 20 is a very conservative default borrowed from nginx: you will probably want
	// to measure if a higher minimum size improves performance for your workloads.
	DefaultMinSize = 20
)

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
func NewGzipLevelAndMinSize(level, minSize int) (func(http.Handler) http.Handler, error) {
	return GzipHandlerWithOpts(GzipCompressionLevel(level), MinSize(minSize))
}

// GzipHandlerWithOpts behaves like NewGzipLevelHandler except it allows the caller
// to specify any of the supported options.
func GzipHandlerWithOpts(opts ...option) (func(http.Handler) http.Handler, error) {
	c := &config{
		gzLevel: gzip.DefaultCompression,
		brLevel: brotliDefaultCompression,
		prefer:  PreferClientThenBrotli,
		minSize: DefaultMinSize,
	}

	for _, o := range opts {
		o(c)
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(vary, acceptEncoding)

			if ac := acceptsCompression(r); ac != acceptsNone {
				gw := &GzipResponseWriter{
					ResponseWriter: w,
					gzLevel:        c.gzLevel,
					brLevel:        c.brLevel,
					minSize:        c.minSize,
					contentTypes:   c.contentTypes,
					blacklist:      c.blacklist,
					accept:         ac,
					prefer:         c.prefer,
				}
				defer gw.Close()

				if _, ok := w.(http.CloseNotifier); ok {
					w = GzipResponseWriterWithCloseNotify{gw}
				} else {
					w = gw
				}
			}

			h.ServeHTTP(w, r)
		})
	}, nil
}

// Used for functional configuration.
type config struct {
	minSize      int
	gzLevel      int
	brLevel      int
	contentTypes []parsedContentType
	blacklist    bool
	prefer       preferType
}

func (c *config) validate() error {
	if c.gzLevel != gzip.DefaultCompression && (c.gzLevel < gzip.BestSpeed || c.gzLevel > gzip.BestCompression) {
		return fmt.Errorf("invalid gzip compression level requested: %d", c.gzLevel)
	}

	if c.brLevel < brotli.BestSpeed || c.brLevel > brotli.BestCompression {
		return fmt.Errorf("invalid brotli compression level requested: %d", c.brLevel)
	}

	if c.minSize < 0 {
		return fmt.Errorf("minimum size must be more than zero")
	}

	return nil
}

type option func(c *config)

// MinSize is an option that controls the minimum size of payloads that
// should be compressed. The default is DefaultMinSize.
func MinSize(size int) option {
	return func(c *config) {
		c.minSize = size
	}
}

// GzipCompressionLevel is an option that controls the Gzip compression
// level to be used when compressing payloads.
// The default is gzip.DefaultCompression.
func GzipCompressionLevel(level int) option {
	return func(c *config) {
		c.gzLevel = level
	}
}

// BrotliCompressionLevel is an option that controls the Brotli compression
// level to be used when compressing payloads.
// The default is 3 (the same default used in the reference brotli C
// implementation).
func BrotliCompressionLevel(level int) option {
	return func(c *config) {
		c.brLevel = level
	}
}

// Prefer controls the behavior of the middleware in case both Gzip and Brotli
// can be used to compress a response (i.e. in case the client supports both
// encodings, and the MIME type of the response is allowed for both encodings).
// See the comments on the PreferXxx constants for the supported values.
func Prefer(prefer preferType) option {
	return func(c *config) {
		c.prefer = prefer
	}
}

// GzipHandler wraps an HTTP handler, to transparently gzip the response body if
// the client supports it (via the Accept-Encoding header). This will compress at
// the default compression level.
func GzipHandler(h http.Handler) http.Handler {
	wrapper, _ := NewGzipLevelHandler(gzip.DefaultCompression)
	return wrapper(h)
}
