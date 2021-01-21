package zstd_test

import (
	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/klauspost/zstd"
)

var _ httpcompression.CompressorProvider = &zstd.Compressor{}
