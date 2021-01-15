package gziphandler

import (
	"fmt"
	"io"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAcceptedCompression(t *testing.T) {
	t.Parallel()
	cases := []struct {
		accept        codings
		onlyGzip      []string
		onlyBrotli    []string
		gzipAndBrotli []string
	}{
		{codings{}, nil, nil, nil},
		{codings{"identity": 1}, nil, nil, nil},
		{codings{"yadda": 1}, nil, nil, nil},
		{codings{"gzip": 1}, []string{"gzip"}, nil, []string{"gzip"}},
		{codings{"gzip": 0.5}, []string{"gzip"}, nil, []string{"gzip"}},
		{codings{"gzip_0": 1}, nil, nil, nil},
		{codings{"gzip": 1, "identity": 1}, []string{"gzip"}, nil, []string{"gzip"}},
		{codings{"gzip": 0}, nil, nil, nil},
		{codings{"br": 1}, nil, []string{"br"}, []string{"br"}},
		{codings{"gzip": 1, "br": 1}, []string{"gzip"}, []string{"br"}, []string{"gzip", "br"}},
		{codings{"gzip": 1, "br": 0.5}, []string{"gzip"}, []string{"br"}, []string{"gzip", "br"}},
		{codings{"gzip": 1, "br": 0}, []string{"gzip"}, nil, []string{"gzip"}},
	}
	onlyGzip := comps{"gzip": comp{comp: fakeCompressor{}}}
	onlyBrotli := comps{"br": comp{comp: fakeCompressor{}}}
	gzipAndBrotli := comps{"gzip": comp{comp: fakeCompressor{}}, "br": comp{comp: fakeCompressor{}}}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.accept), func(t *testing.T) {
			t.Run("onlyGzip", func(t *testing.T) {
				a := acceptedCompression(c.accept, onlyGzip)
				sort.Strings(a)
				sort.Strings(c.onlyGzip)
				assert.Equal(t, a, c.onlyGzip)
			})
			t.Run("onlyBrotli", func(t *testing.T) {
				a := acceptedCompression(c.accept, onlyBrotli)
				sort.Strings(a)
				sort.Strings(c.onlyBrotli)
				assert.Equal(t, a, c.onlyBrotli)
			})
			t.Run("gzipAndBrotli", func(t *testing.T) {
				a := acceptedCompression(c.accept, gzipAndBrotli)
				sort.Strings(a)
				sort.Strings(c.gzipAndBrotli)
				assert.Equal(t, a, c.gzipAndBrotli)
			})
		})
	}
}

type fakeCompressor struct{}

func (fakeCompressor) Get(_ io.Writer) io.WriteCloser { return nil }
