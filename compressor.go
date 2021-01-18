package httpcompression

import (
	"io"
)

type CompressorProvider interface {
	// Get returns a writer that writes compressed output to the supplied parent io.Writer.
	// When Close() is called on the returned io.WriteCloser, it is guaranteed that it will
	// not be used anymore so implementations can safely recycle the compressor (e.g. put the
	// WriteCloser in a pool to be reused by a later call to Get).
	// The returned io.WriteCloser can optionally implement the Flusher interface if it is
	// able to flush data buffered internally.
	Get(parent io.Writer) (compressor io.WriteCloser)
}

type Flusher interface {
	// Flush flushes the data buffered internally by the Writer.
	// Flush does not need to internally flush the parent Writer.
	Flush() error
}

type comps map[string]comp

type comp struct {
	comp     CompressorProvider
	priority int
}

func Compressor(contentEncoding string, priority int, compressor CompressorProvider) Option {
	return func(c *config) {
		if compressor == nil {
			delete(c.compressor, contentEncoding)
			return
		}
		c.compressor[contentEncoding] = comp{compressor, priority}
	}
}
