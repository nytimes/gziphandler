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

// Compressor returns an Option that sets the CompressorProvider for a specific Content-Encoding.
// If multiple CompressorProviders are set for the same Content-Encoding, the last one is used.
// If compressor is nil, it disables the specified Content-Encoding.
// Priority is used to specify the priority of the Content-Encoding. A higher number means higher
// priority. See PreferType to understand how priority is used to select the Content-Encoding for
// a specific request.
func Compressor(contentEncoding string, priority int, compressor CompressorProvider) Option {
	return func(c *config) error {
		if compressor == nil {
			delete(c.compressor, contentEncoding)
			return nil
		}
		c.compressor[contentEncoding] = comp{compressor, priority}
		return nil
	}
}

func mustCompressor(contentEncoding string, priority int, compressor CompressorProvider, err error) Option {
	if err != nil {
		return errorOption(err)
	}
	return Compressor(contentEncoding, priority, compressor)
}
