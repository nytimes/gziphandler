package brotli_test

import (
	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
)

var _ httpcompression.CompressorProvider = &brotli.Compressor{}
