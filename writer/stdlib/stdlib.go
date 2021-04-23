package stdlib

import (
	"compress/gzip"
	"github.com/NYTimes/gziphandler/writer"
	"io"
	"sync"
)

// gzipWriterPools stores a sync.Pool for each compression level for reuse of
// gzip.Writers. Use poolIndex to covert a compression level to an index into
// gzipWriterPools.
var gzipWriterPools [gzip.BestCompression - gzip.BestSpeed + 2]*sync.Pool

func init() {
	for i := gzip.BestSpeed; i <= gzip.BestCompression; i++ {
		addLevelPool(i)
	}
	addLevelPool(gzip.DefaultCompression)
}

// poolIndex maps a compression level to its index into gzipWriterPools. It
// assumes that level is a valid gzip compression level.
func poolIndex(level int) int {
	// gzip.DefaultCompression == -1, so we need to treat it special.
	if level == gzip.DefaultCompression {
		return gzip.BestCompression - gzip.BestSpeed + 1
	}
	return level - gzip.BestSpeed
}

func addLevelPool(level int) {
	gzipWriterPools[poolIndex(level)] = &sync.Pool{
		New: func() interface{} {
			// NewWriterLevel only returns error on a bad level, we are guaranteeing
			// that this will be a valid level so it is okay to ignore the returned
			// error.
			w, _ := gzip.NewWriterLevel(nil, level)
			return w
		},
	}
}

type pooledWriter struct {
	*gzip.Writer
	index int
}

func (pw *pooledWriter) Close() error {
	err := pw.Writer.Close()
	gzipWriterPools[pw.index].Put(pw.Writer)
	pw.Writer = nil
	return err
}

func NewWriter(w io.Writer, level int) writer.GzipWriter {
	index := poolIndex(level)
	gzw := gzipWriterPools[index].Get().(*gzip.Writer)
	gzw.Reset(w)
	return &pooledWriter{
		Writer: gzw,
		index:  index,
	}
}

func ImplementationInfo() string {
	return "compress/zlib"
}
