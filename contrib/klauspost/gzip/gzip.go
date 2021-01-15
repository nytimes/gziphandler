package gzip

import (
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
)

const (
	Encoding     = "gzip"
	DefaultLevel = gzip.DefaultCompression
)

type compressor struct {
	pool sync.Pool
	opts Options
}

type Options struct {
	Level int
}

func New(opts Options) (c *compressor, err error) {
	gw, err := gzip.NewWriterLevel(nil, opts.Level)
	if err != nil {
		return nil, err
	}
	c = &compressor{opts: opts}
	c.pool.Put(gw)
	return c, nil
}

func (c *compressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*writer); ok {
		gw.Reset(w)
		return gw
	}
	gw, err := gzip.NewWriterLevel(w, c.opts.Level)
	if err != nil {
		panic(err)
	}
	return &writer{
		Writer: gw,
		c:      c,
	}
}

type writer struct {
	*gzip.Writer
	c *compressor
}

func (w *writer) Close() error {
	err := w.Writer.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}
