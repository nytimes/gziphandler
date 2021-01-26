// +build race

package pgzip_test

import (
	"runtime"
	"testing"

	"github.com/CAFxX/httpcompression/contrib/internal"
	"github.com/CAFxX/httpcompression/contrib/klauspost/pgzip"
)

func TestPgzipRace(t *testing.T) {
	t.Parallel()
	c, _ := pgzip.New(pgzip.Options{BlockSize: 1 << 20, Blocks: runtime.GOMAXPROCS(0)})
	internal.RaceTestCompressionProvider(c, 100)
}
