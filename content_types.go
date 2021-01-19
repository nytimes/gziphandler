package httpcompression

import "mime"

// text/html text/xml text/plain text/css text/javascript application/javascript application/json application/xml
// text/css text/xml application/javascript application/atom+xml application/rss+xml text/mathml text/plain text/x-component image/svg+xml application/json application/vnd.google-earth.kml+xml application/x-perl application/xhtml+xml application/xspf+xml
// ^text/
// [+](xml|json|cbor)$

// ContentTypes specifies a list of content types to compare
// the Content-Type header to before compressing. If none
// match, and blacklist is false, the response will be returned as-is.
//
// Content types are compared in a case-insensitive, whitespace-ignored
// manner.
//
// A MIME type without any other directive will match a content type
// that has the same MIME type, regardless of that content type's other
// directives. I.e., "text/html" will match both "text/html" and
// "text/html; charset=utf-8".
//
// A MIME type with any other directive will only match a content type
// that has the same MIME type and other directives. I.e.,
// "text/html; charset=utf-8" will only match "text/html; charset=utf-8".
//
// If blacklist is true then only content types that do not match the
// provided list of types are compressed. If blacklist is false, only
// content types that match the provided list are compressed.
//
// By default, responses are compressed regardless of Content-Type.
func ContentTypes(types []string, blacklist bool) Option {
	return func(c *config) error {
		c.contentTypes = []parsedContentType{}
		for _, v := range types {
			mediaType, params, err := mime.ParseMediaType(v)
			if err != nil {
				return err
			}
			c.contentTypes = append(c.contentTypes, parsedContentType{mediaType, params})
		}
		c.blacklist = blacklist
		return nil
	}
}

// Parsed representation of one of the inputs to ContentTypes.
// See https://golang.org/pkg/mime/#ParseMediaType
type parsedContentType struct {
	mediaType string
	params    map[string]string
}

// equals returns whether this content type matches another content type.
func (pct *parsedContentType) equals(mediaType string, params map[string]string) bool {
	if pct.mediaType != mediaType {
		return false
	}
	// if pct has no params, don't care about other's params
	if len(pct.params) == 0 {
		return true
	}
	// FIXME: the slow path is asinine, unnecessary, and should be eradicated.
	return pct.equalsSlow(mediaType, params)
}

func (pct *parsedContentType) equalsSlow(mediaType string, params map[string]string) bool {
	// if pct has any params, they must be identical to other's.
	if len(pct.params) != len(params) {
		return false
	}
	for k, v := range pct.params {
		if w, ok := params[k]; !ok || v != w {
			return false
		}
	}
	return true
}
