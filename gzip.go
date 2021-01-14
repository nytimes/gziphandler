package gziphandler

import (
	"compress/gzip"
	"fmt"
	"io"
	"sync"
)

// gzipWriterPools stores a sync.Pool for each compression level for reuse of
// gzip.Writers. Use poolIndex to covert a compression level to an index into
// gzipWriterPools.
var gzipWriterPools [gzip.BestCompression - gzip.BestSpeed + 2]sync.Pool

func init() {
	for i := gzip.BestSpeed; i <= gzip.BestCompression; i++ {
		addGzipLevelPool(i)
	}
	addGzipLevelPool(gzip.DefaultCompression)
}

// gzipPoolIndex maps a compression level to its index into gzipWriterPools. It
// assumes that level is a valid gzip compression level.
func gzipPoolIndex(level int) int {
	// gzip.DefaultCompression == -1, so we need to treat it special.
	// TODO: handle gzip.HuffmanOnly
	if level == gzip.DefaultCompression {
		return gzip.BestCompression - gzip.BestSpeed + 1
	}
	return level - gzip.BestSpeed
}

func addGzipLevelPool(level int) {
	gzipWriterPools[gzipPoolIndex(level)] = sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(nil, level)
			if err != nil {
				// NewWriterLevel only returns error on a bad level, we are guaranteeing
				// that this will be a valid level so this panic should never fire.
				panic(fmt.Errorf("gzip writer initialization: %v", err))
			}
			return w
		},
	}
}

func getGzipWriter(w io.Writer, level int) *gzip.Writer {
	gw, _ := gzipWriterPools[gzipPoolIndex(level)].Get().(*gzip.Writer)
	gw.Reset(w)
	return gw
}

func putGzipWriter(gw *gzip.Writer, level int) {
	gw.Reset(nil)
	gzipWriterPools[gzipPoolIndex(level)].Put(gw)
}

type DefaultGzipCompressor struct{}

func (_ DefaultGzipCompressor) Get(w io.Writer, level int) io.WriteCloser {
	return DefaultGzipWriter{getGzipWriter(w, level), level}
}

type DefaultGzipWriter struct {
	*gzip.Writer
	level int
}

func (w DefaultGzipWriter) Close() error {
	err := w.Writer.Close()
	putGzipWriter(w.Writer, w.level)
	return err
}
