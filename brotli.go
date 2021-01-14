package gziphandler

import (
	"io"
	"sync"

	"github.com/andybalholm/brotli"
)

var brotliWriterPools [brotli.BestCompression - brotli.BestSpeed + 1]sync.Pool

const brotliDefaultCompression = 3

func init() {
	for i := brotli.BestSpeed; i <= brotli.BestCompression; i++ {
		addBrotliLevelPool(i)
	}
}

func brotliPoolIndex(level int) int {
	return level - brotli.BestSpeed
}

func addBrotliLevelPool(level int) {
	brotliWriterPools[brotliPoolIndex(level)] = sync.Pool{
		New: func() interface{} {
			return brotli.NewWriterLevel(nil, level)
		},
	}
}

func getBrotliWriter(w io.Writer, level int) *brotli.Writer {
	bw, _ := brotliWriterPools[brotliPoolIndex(level)].Get().(*brotli.Writer)
	bw.Reset(w)
	// TODO: use BROTLI_MODE_TEXT and BROTLI_MODE_FONT
	return bw
}

func putBrotliWriter(bw *brotli.Writer, level int) {
	bw.Reset(nil) // avoid keeping writer alive
	brotliWriterPools[brotliPoolIndex(level)].Put(bw)
}

type DefaultBrotliCompressor struct{}

func (_ DefaultBrotliCompressor) Get(w io.Writer, level int) io.WriteCloser {
	return DefaultBrotliWriter{getBrotliWriter(w, level), level}
}

type DefaultBrotliWriter struct {
	*brotli.Writer
	level int
}

func (w DefaultBrotliWriter) Close() error {
	err := w.Writer.Close()
	putBrotliWriter(w.Writer, w.level)
	return err
}
