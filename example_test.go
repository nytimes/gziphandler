package httpcompression_test

import (
	"log"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	"github.com/CAFxX/httpcompression/contrib/klauspost/gzip"
	"github.com/CAFxX/httpcompression/contrib/klauspost/zstd"
	kpzstd "github.com/klauspost/compress/zstd"
)

func Example() {
	brEnc, err := brotli.New(brotli.Options{})
	if err != nil {
		log.Fatal(err)
	}
	gzEnc, err := gzip.New(gzip.Options{})
	if err != nil {
		log.Fatal(err)
	}
	_, _ = httpcompression.Handler(
		httpcompression.Compressor(brotli.Encoding, 1, brEnc),
		httpcompression.Compressor(gzip.Encoding, 0, gzEnc),
		httpcompression.Prefer(httpcompression.PreferServer),
		httpcompression.MinSize(100),
		httpcompression.ContentTypes([]string{
			"image/jpeg",
			"image/gif",
			"image/png",
		}, true),
	)
}

func ExampleWithDictionary() {
	zEnc, err := zstd.New()
	if err != nil {
		log.Fatal(err)
	}
	dict := []byte("dictionary contents") // replace with dictionary contents
	zdEnc, err := zstd.New(kpzstd.WithEncoderDict(dict))
	if err != nil {
		log.Fatal(err)
	}
	_, _ = httpcompression.Handler(
		// Add the zstd compressor with the dictionary.
		// We need to pick a custom content-encoding name. It is recommended to:
		// - avoid names that contain standard names (e.g. "gzip", "deflate", "br" or "zstd")
		// - include the dictionary ID, so that multiple dictionaries can be used (including
		//   e.g. multiple versions of the same dictionary)
		httpcompression.Compressor("z00000000", 3, zdEnc),
		httpcompression.Compressor(zstd.Encoding, 2, zEnc),
		httpcompression.Prefer(httpcompression.PreferServer),
		httpcompression.MinSize(0),
		httpcompression.ContentTypes([]string{
			"image/jpeg",
			"image/gif",
			"image/png",
		}, true),
	)
}
