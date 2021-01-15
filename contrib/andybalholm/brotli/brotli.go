package brotli

import (
	"fmt"
	"io"
	"sync"

	"github.com/andybalholm/brotli"
)

const (
	Encoding     = "br"
	DefaultLevel = brotli.DefaultCompression
)

type Options = brotli.WriterOptions

type compressor struct {
	pool sync.Pool
	opts Options
}

func New(opts Options) (c *compressor, err error) {
	defer func() {
		if r := recover(); r != nil {
			c, err = nil, fmt.Errorf("panic: %v", r)
		}
	}()
	gw := brotli.NewWriterOptions(nil, opts)
	c = &compressor{opts: opts}
	c.pool.Put(gw)
	return c, nil
}

func (c *compressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*writer); ok {
		gw.Reset(w)
		return gw
	}
	gw := brotli.NewWriterOptions(w, c.opts)
	return &writer{
		Writer: gw,
		c:      c,
	}
}

type writer struct {
	*brotli.Writer
	c *compressor
}

func (w *writer) Close() error {
	err := w.Writer.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}
