package custom

import "io"

type Compressor interface {
	Get(parent io.Writer, level int) (compressor io.WriteCloser)
}

type Flusher interface {
	Flush() error
}
