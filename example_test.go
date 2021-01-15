package gziphandler_test

import (
	"log"

	"github.com/CAFxX/gziphandler"
	"github.com/CAFxX/gziphandler/contrib/andybalholm/brotli"
	"github.com/CAFxX/gziphandler/contrib/klauspost/gzip"
	"github.com/CAFxX/gziphandler/contrib/klauspost/zstd"
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
	_, _ = gziphandler.Middleware(
		gziphandler.Compressor(brotli.Encoding, 1, brEnc),
		gziphandler.Compressor(gzip.Encoding, 0, gzEnc),
		gziphandler.Prefer(gziphandler.PreferServer),
		gziphandler.MinSize(100),
		gziphandler.ContentTypes([]string{
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
	dictID := "0000"                      // replace with dictionary ID
	zdEnc, err := zstd.New(kpzstd.WithEncoderDict(dict))
	if err != nil {
		log.Fatal(err)
	}
	_, _ = gziphandler.Middleware(
		gziphandler.Compressor(zstd.Encoding+"_"+dictID, 1, zdEnc), // non-standard, client support needed
		gziphandler.Compressor(zstd.Encoding, 0, zEnc),
		gziphandler.Prefer(gziphandler.PreferServer),
		gziphandler.MinSize(0),
		gziphandler.ContentTypes([]string{
			"image/jpeg",
			"image/gif",
			"image/png",
		}, true),
	)
}
