package gziphandler

// Prefer controls the behavior of the middleware in case both Gzip and Brotli
// can be used to compress a response (i.e. in case the client supports both
// encodings, and the MIME type of the response is allowed for both encodings).
// See the comments on the PreferType constants for the supported values.
func Prefer(prefer PreferType) Option {
	return func(c *config) {
		c.prefer = prefer
	}
}

// PreferType allows to control the choice of compression algorithm when
// multiple algorithms are allowed by both client and server.
type PreferType int

const (
	// PreferGzip uses Gzip if the client supports both Brotli and Gzip.
	PreferGzip PreferType = iota

	// PreferBrotli uses Brotli if the client supports both Brotli and Gzip.
	PreferBrotli

	// PreferClientThenGzip uses the client preference, or Gzip if no preference is specified.
	PreferClientThenGzip

	// PreferClientThenBrotli uses the client preference, or Brotli if no preference is specified.
	PreferClientThenBrotli
)

type priorityType int

const (
	priorityGzip priorityType = iota
	priorityBrotli
)

// returns which scheme we should try first
func (p PreferType) priorityFor(a acceptsType) priorityType {
	if p == PreferClientThenGzip || p == PreferClientThenBrotli {
		switch a {
		case acceptsBrotliThenGzip, acceptsBrotli:
			return priorityBrotli
		case acceptsGzipThenBrotli, acceptsGzip:
			return priorityGzip
		case acceptsGzipAndBrotli:
			if p == PreferClientThenGzip {
				return priorityGzip
			}
			return priorityBrotli
		}
	}
	if p == PreferGzip {
		return priorityGzip
	}
	return priorityBrotli
}
