package gzip_test

import (
	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/klauspost/gzip"
)

var _ httpcompression.CompressorProvider = &gzip.Compressor{}
