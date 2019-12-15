package gziphandler

type preferType int

const (
	// PreferGzip uses Gzip if the client supports both Brotli and Gzip.
	PreferGzip preferType = iota

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

// returns true if we should try gzip first
func (p preferType) priorityFor(a acceptsType) priorityType {
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
