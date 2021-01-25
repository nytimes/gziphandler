package httpcompression

import (
	"fmt"
	"sort"
)

// Prefer controls the behavior of the middleware in case both Gzip and Brotli
// can be used to compress a response (i.e. in case the client supports both
// encodings, and the MIME type of the response is allowed for both encodings).
// See the comments on the PreferType constants for the supported values.
func Prefer(prefer PreferType) Option {
	return func(c *config) error {
		switch prefer {
		case PreferServer, PreferClient:
			c.prefer = prefer
			return nil
		default:
			return fmt.Errorf("unknown prefer type: %v", prefer)
		}
	}
}

// PreferType allows to control the choice of compression algorithm when
// multiple algorithms are allowed by both client and server.
type PreferType byte

const (
	// PreferServer prefers compressors in the order specified on the server.
	// If two or more compressors have the same priority on the server, the client preference is taken into consideration.
	// If both server and client do no specify a preference between two or more compressors, the order is determined by the name of the encoding.
	// PreferServer is the default.
	PreferServer PreferType = iota

	// PreferClient prefers compressors in the order specified by the client.
	// If two or more compressors have the same priority according to the client, the server priority is taken into consideration.
	// If both server and client do no specify a preference between two or more compressors, the order is determined by the name of the encoding.
	PreferClient
)

func preferredEncoding(accept codings, comps comps, common []string, prefer PreferType) string {
	if len(common) == 0 {
		panic("no common encoding")
	}
	switch prefer {
	case PreferServer:
		sort.Slice(common, func(i, j int) bool {
			ci, cj := comps[common[i]].priority, comps[common[j]].priority
			if ci != cj {
				return ci > cj // desc
			}
			ai, aj := accept[common[i]], accept[common[j]]
			if ai != aj {
				return ai > aj // desc
			}
			return common[i] < common[j] // asc
		})
	case PreferClient:
		sort.Slice(common, func(i, j int) bool {
			ai, aj := accept[common[i]], accept[common[j]]
			if ai != aj {
				return ai > aj // desc
			}
			ci, cj := comps[common[i]].priority, comps[common[j]].priority
			if ci != cj {
				return ci > cj // desc
			}
			return common[i] < common[j] // asc
		})
	default:
		panic("unknown prefer type")
	}
	return common[0]
}
