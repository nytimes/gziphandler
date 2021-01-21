package pgzip_test

import (
	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/klauspost/pgzip"
)

var _ httpcompression.CompressorProvider = &pgzip.Compressor{}
