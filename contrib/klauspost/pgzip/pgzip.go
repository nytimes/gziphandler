package pgzip

import (
	"io"
	"sync"

	"github.com/klauspost/pgzip"
)

const (
	Encoding           = "gzip"
	DefaultCompression = pgzip.DefaultCompression
)

type compressor struct {
	pool sync.Pool
	opts Options
}

type Options struct {
	Level     int
	BlockSize int
	Blocks    int
}

func New(opts Options) (c *compressor, err error) {
	gw, err := pgzip.NewWriterLevel(nil, opts.Level)
	if err != nil {
		return nil, err
	}
	err = gw.SetConcurrency(opts.BlockSize, opts.Blocks)
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
	gw, err := pgzip.NewWriterLevel(w, c.opts.Level)
	if err != nil {
		panic(err)
	}
	err = gw.SetConcurrency(c.opts.BlockSize, c.opts.Blocks)
	if err != nil {
		panic(err)
	}
	return &writer{
		Writer: gw,
		c:      c,
	}
}

type writer struct {
	*pgzip.Writer
	c *compressor
}

func (w *writer) Close() error {
	err := w.Writer.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}
