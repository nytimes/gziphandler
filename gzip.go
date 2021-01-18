package httpcompression

import (
	"compress/gzip"
	"io"
	"sync"
)

type defaultGzipCompressor struct {
	pool  sync.Pool
	level int
}

func NewDefaultGzipCompressor(level int) (*defaultGzipCompressor, error) {
	gw, err := gzip.NewWriterLevel(nil, level)
	if err != nil {
		return nil, err
	}
	c := &defaultGzipCompressor{level: level}
	c.pool.Put(gw)
	return c, nil
}

func (c *defaultGzipCompressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*defaultGzipWriter); ok {
		gw.Reset(w)
		return gw
	}
	gw, _ := gzip.NewWriterLevel(w, c.level)
	return &defaultGzipWriter{
		Writer: gw,
		c:      c,
	}
}

type defaultGzipWriter struct {
	*gzip.Writer
	c *defaultGzipCompressor
}

func (w *defaultGzipWriter) Close() error {
	err := w.Writer.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}
