package zstd

import (
	"io"
	"sync"

	"github.com/klauspost/compress/zstd"
)

const Encoding = "zstd"

type compressor struct {
	pool sync.Pool
	opts []zstd.EOption
}

func New(opts ...zstd.EOption) (c *compressor, err error) {
	opts = append([]zstd.EOption(nil), opts...)
	gw, err := zstd.NewWriter(nil, opts...)
	if err != nil {
		return nil, err
	}
	c = &compressor{opts: opts}
	c.pool.Put(gw)
	return c, nil
}

func (c *compressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*zstdWriter); ok {
		gw.Reset(w)
		return gw
	}
	gw, err := zstd.NewWriter(w, c.opts...)
	if err != nil {
		panic(err)
	}
	return &zstdWriter{
		Encoder: gw,
		c:       c,
	}
}

type zstdWriter struct {
	*zstd.Encoder
	c *compressor
}

func (w *zstdWriter) Close() error {
	err := w.Encoder.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}
