package internal

import (
	"bytes"
	"runtime"
	"sync"

	"github.com/CAFxX/httpcompression"
)

func RaceTestCompressionProvider(c httpcompression.CompressorProvider, n int) {
	var wg sync.WaitGroup
	for i := runtime.GOMAXPROCS(0); i >= 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b := &bytes.Buffer{}
			for j := 0; j < n; j++ {
				b.Reset()
				w := c.Get(b)
				w.Write([]byte("hello world"))
				w.Close()
			}
		}()
	}
	wg.Wait()
}
