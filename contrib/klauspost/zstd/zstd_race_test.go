// +build race

package zstd_test

import (
	"testing"

	"github.com/CAFxX/httpcompression/contrib/internal"
	"github.com/CAFxX/httpcompression/contrib/klauspost/zstd"
)

func TestZstdRace(t *testing.T) {
	t.Parallel()
	c, _ := zstd.New()
	internal.RaceTestCompressionProvider(c, 100)
}
