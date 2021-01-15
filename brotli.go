package gziphandler

import (
	"github.com/CAFxX/gziphandler/contrib/andybalholm/brotli"
	_brotli "github.com/andybalholm/brotli"
)

func NewDefaultBrotliCompressor(quality int) (c CompressorProvider, err error) {
	return brotli.New(_brotli.WriterOptions{Quality: quality})
}
