package custom

import "io"

type Compressor interface {
	// Get returns a writer that writes compressed output to the supplied parent io.Writer.
	// The level is the compression level that should be used when compressing. The range
	// depends on the compression algorithm; it is in the range -1 to 9 for gzip and 0 to 11
	// for brotli.
	// When Close() is called on the returned io.WriteCloser, it is guaranteed that it will
	// not be used anymore so implementations can safely recycle the compressor (e.g. put the
	// WriteCloser in a pool to be reused by a later call to Get).
	// The returned io.WriteCloser can optionally implement the Flusher interface if it is
	// able to flush data buffered internally.
	Get(parent io.Writer, level int) (compressor io.WriteCloser)
}

type Flusher interface {
	// Flush flushes the data buffered internally by the Writer.
	// Flush does not need to internally flush the parent Writer.
	Flush() error
}
