// +build race

package brotli_test

import (
	"testing"

	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	"github.com/CAFxX/httpcompression/contrib/internal"
)

func TestZstdRace(t *testing.T) {
	t.Parallel()
	c, _ := brotli.New(brotli.Options{})
	internal.RaceTestCompressionProvider(c, 100)
}
