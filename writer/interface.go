package writer

import "io"

type GzipWriter interface {
	Close() error
	Flush() error
	Write(p []byte) (int, error)
}

type GzipWriterFactory = func(writer io.Writer, level int) GzipWriter
