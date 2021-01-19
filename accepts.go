package httpcompression

import (
	"math"
	"strconv"
	"strings"
)

const (
	// defaultQValue is the default qvalue to assign to an encoding if no explicit qvalue is set.
	// This is actually kind of ambiguous in RFC 2616, so hopefully it's correct.
	// The examples seem to indicate that it is.
	defaultQValue = 1.0
)

// acceptedCompression returns the list of common compression scheme supported by client and server.
func acceptedCompression(accept codings, comps comps) []string {
	var s []string
	// pick smallest N to do O(N) iterations
	if len(accept) < len(comps) {
		for k, v := range accept {
			if v > 0 && comps[k].comp != nil {
				s = append(s, k)
			}
		}
	} else {
		for k, v := range comps {
			if v.comp != nil && accept[k] > 0 {
				s = append(s, k)
			}
		}
	}
	return s
}

// parseEncodings attempts to parse a list of codings, per RFC 2616, as might
// appear in an Accept-Encoding header. It returns a map of content-codings to
// quality values.
// Errors encountered during parsing the codings are ignored.
//
// See: http://tools.ietf.org/html/rfc2616#section-14.3.
func parseEncodings(s string) codings {
	c := make(codings)
	for _, ss := range strings.Split(s, ",") {
		coding, qvalue := parseCoding(ss)
		if coding == "" {
			continue
		}
		c[coding] = qvalue
	}
	return c
}

// parseCoding parses a single conding (content-coding with an optional qvalue),
// as might appear in an Accept-Encoding header. It attempts to forgive minor
// formatting errors.
func parseCoding(s string) (coding string, qvalue float64) {
	qvalue = defaultQValue

	p := strings.IndexRune(s, ';')
	if p != -1 {
		if part := strings.Replace(s[p+1:], " ", "", -1); strings.HasPrefix(part, "q=") {
			qvalue, _ = strconv.ParseFloat(part[2:], 64)
			if qvalue < 0.0 || math.IsNaN(qvalue) {
				qvalue = 0.0
			} else if qvalue > 1.0 {
				qvalue = 1.0
			}
		}
	} else {
		p = len(s)
	}
	coding = strings.ToLower(strings.TrimSpace(s[:p]))
	return
}
