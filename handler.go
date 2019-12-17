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

// Middleware returns a wrapper function (often known as middleware)
// which can be used to wrap an HTTP handler to transparently compress the response
// body if the client supports it (via the Accept-Encoding header).
// It is possible to pass one or more options to modify the middleware configuration.
// An error will be returned if invalid options are given.
func Middleware(opts ...Option) (func(http.Handler) http.Handler, error) {
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
				gw := &gzipResponseWriter{
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
					w = gzipResponseWriterWithCloseNotify{gw}
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
	prefer       PreferType
}

func (c *config) validate() error {
	if c.gzLevel != gzip.DefaultCompression && (c.gzLevel < gzip.BestSpeed || c.gzLevel > gzip.BestCompression) {
		return fmt.Errorf("invalid gzip compression level requested: %d", c.gzLevel)
	}

	if c.brLevel < brotli.BestSpeed || c.brLevel > brotli.BestCompression {
		return fmt.Errorf("invalid brotli compression level requested: %d", c.brLevel)
	}

	if c.minSize < 0 {
		return fmt.Errorf("minimum size must be more than zero: %d", c.minSize)
	}

	switch c.prefer {
	case PreferBrotli, PreferClientThenBrotli, PreferClientThenGzip, PreferGzip:
	default:
		return fmt.Errorf("invalid prefer config: %v", c.prefer)
	}

	return nil
}

// Option can be passed to Middleware to control its configuration.
type Option func(c *config)

// MinSize is an option that controls the minimum size of payloads that
// should be compressed. The default is DefaultMinSize.
func MinSize(size int) Option {
	return func(c *config) {
		c.minSize = size
	}
}

// GzipCompressionLevel is an option that controls the Gzip compression
// level to be used when compressing payloads.
// The default is gzip.DefaultCompression.
func GzipCompressionLevel(level int) Option {
	return func(c *config) {
		c.gzLevel = level
	}
}

// BrotliCompressionLevel is an option that controls the Brotli compression
// level to be used when compressing payloads.
// The default is 3 (the same default used in the reference brotli C
// implementation).
func BrotliCompressionLevel(level int) Option {
	return func(c *config) {
		c.brLevel = level
	}
}
