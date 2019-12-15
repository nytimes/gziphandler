package gziphandler

import "mime"

type handleType int

const (
	handleNone handleType = iota
	handleGzip
	handleBrotli
)

// returns how we should handle the request
func handleContentTypes(gzipContentTypes, brotliContentTypes []parsedContentType, blacklist bool, ct string, prefer preferType, accept acceptsType) handleType {
	switch prefer.priorityFor(accept) {
	case priorityGzip:
		if accept.gzip() && handleContentType(gzipContentTypes, blacklist, ct) {
			return handleGzip
		}
		if accept.brotli() && handleContentType(brotliContentTypes, blacklist, ct) {
			return handleBrotli
		}
	case priorityBrotli:
		if accept.brotli() && handleContentType(brotliContentTypes, blacklist, ct) {
			return handleBrotli
		}
		if accept.gzip() && handleContentType(gzipContentTypes, blacklist, ct) {
			return handleGzip
		}
	}
	return handleNone
}

// returns true if we've been configured to compress the specific content type.
func handleContentType(contentTypes []parsedContentType, blacklist bool, ct string) bool {
	// If contentTypes is empty we handle all content types.
	if len(contentTypes) == 0 {
		return true
	}
	return handleContentTypeSlow(contentTypes, blacklist, ct)
}

func handleContentTypeSlow(contentTypes []parsedContentType, blacklist bool, ct string) bool {
	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return false
	}

	for _, c := range contentTypes {
		if c.equals(mediaType, params) {
			return !blacklist
		}
	}

	return blacklist
}
