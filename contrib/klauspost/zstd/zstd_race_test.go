// +build race

package zstd_test

import (
	"bytes"
	"runtime"
	"sync"
	"testing"

	"github.com/CAFxX/httpcompression/contrib/klauspost/zstd"
)

func TestZstdRace(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup
	c, _ := zstd.New()
	for i := runtime.GOMAXPROCS(0); i >= 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b := &bytes.Buffer{}
			for j := 0; j < 100; j++ {
				b.Reset()
				w := c.Get(b)
				w.Write([]byte("hello world"))
				w.Close()
			}
		}()
	}
	wg.Wait()
}
