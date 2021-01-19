package httpcompression

import "mime"

// returns true if we've been configured to compress the specific content type.
func handleContentType(ct string, contentTypes []parsedContentType, blacklist bool) bool {
	// If contentTypes is empty we handle all content types.
	if len(contentTypes) == 0 {
		return !blacklist
	}
	return handleContentTypeSlow(ct, contentTypes, blacklist)
}

func handleContentTypeSlow(ct string, contentTypes []parsedContentType, blacklist bool) bool {
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
