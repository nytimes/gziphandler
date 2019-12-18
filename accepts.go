package gziphandler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const (
	// defaultQValue is the default qvalue to assign to an encoding if no explicit qvalue is set.
	// This is actually kind of ambiguous in RFC 2616, so hopefully it's correct.
	// The examples seem to indicate that it is.
	defaultQValue = 1.0
)

type acceptsType int

const (
	acceptsNone acceptsType = iota
	acceptsGzip
	acceptsBrotli
	acceptsGzipAndBrotli
	acceptsGzipThenBrotli // client prefers Gzip
	acceptsBrotliThenGzip // client prefers Brotli
)

func (a acceptsType) gzip() bool {
	return a == acceptsGzip || a == acceptsGzipAndBrotli || a == acceptsGzipThenBrotli || a == acceptsBrotliThenGzip
}

func (a acceptsType) brotli() bool {
	return a == acceptsBrotli || a == acceptsGzipAndBrotli || a == acceptsGzipThenBrotli || a == acceptsBrotliThenGzip
}

// acceptsCompression returns which kind of compressions the client will accept.
func acceptsCompression(r *http.Request) acceptsType {
	acceptedEncodings, _ := parseEncodings(r.Header.Get(acceptEncoding))
	gzip, br := acceptedEncodings["gzip"], acceptedEncodings["br"]
	switch {
	case gzip > 0 && br > 0 && br > gzip:
		return acceptsBrotliThenGzip
	case gzip > 0 && br > 0 && br < gzip:
		return acceptsGzipThenBrotli
	case gzip > 0 && br > 0 && br == gzip:
		return acceptsGzipAndBrotli
	case gzip > 0 && br == 0:
		return acceptsGzip
	case gzip == 0 && br > 0:
		return acceptsBrotli
	default:
		return acceptsNone
	}
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
	for n, part := range strings.SplitN(s, ";", 2) {
		part = strings.TrimSpace(part)
		qvalue = defaultQValue

		if n == 0 {
			coding = strings.ToLower(part)
		} else if part := strings.Replace(part, " ", "", -1); strings.HasPrefix(part, "q=") {
			qvalue, _ = strconv.ParseFloat(strings.TrimPrefix(part, "q="), 64)

			if qvalue < 0.0 || math.IsNaN(qvalue) {
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
