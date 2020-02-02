package gziphandler

import "net/http"
import "testing"

func TestAccepts(t *testing.T) {
	t.Parallel()

	cases := []struct {
		AcceptEncoding string
		AcceptsType    acceptsType
	}{
		{"", acceptsNone},
		{"bazinga", acceptsNone},
		{"identity", acceptsNone},
		{"deflate", acceptsNone},
		// This is not really correct according to
		// https://tools.ietf.org/html/rfc7231#section-5.3.4
		// but it's a safe choice.
		{"*", acceptsNone},

		{"gzip", acceptsGzip},
		{"gzip,bazinga", acceptsGzip},
		{"gzip;q=1", acceptsGzip},
		{"gzip;q=0.5", acceptsGzip},
		{"gzip;q=0", acceptsNone},

		{"br", acceptsBrotli},
		{"br,identity", acceptsBrotli},
		{"br;q=1", acceptsBrotli},
		{"br;q=0.5", acceptsBrotli},
		{"br;q=0", acceptsNone},

		{"br,gzip", acceptsGzipAndBrotli},
		{"gzip,br", acceptsGzipAndBrotli},
		{"gzip,br,identity", acceptsGzipAndBrotli},
		{"gzip;q=1,br;q=1", acceptsGzipAndBrotli},
		{"gzip;q=0.5,br;q=0.5", acceptsGzipAndBrotli},
		{"gzip;q=0,br;q=0", acceptsNone},

		{"gzip;q=1,br;q=0.5", acceptsGzipThenBrotli},
		{"gzip;q=0.5,br;q=0", acceptsGzip},

		{"gzip;q=0.5,br;q=1", acceptsBrotliThenGzip},
		{"gzip;q=0,br;q=0.5", acceptsBrotli},
	}

	for _, c := range cases {
		t.Run(c.AcceptEncoding, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "/", nil)
			r.Header.Set("Accept-Encoding", c.AcceptEncoding)
			a := acceptsCompression(r)
			if a != c.AcceptsType {
				t.Fatalf("got %q, want %q", a, c.AcceptsType)
			}
		})
	}
}
