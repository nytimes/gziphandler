package httpcompression_test

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	"github.com/CAFxX/httpcompression/contrib/klauspost/gzip"
	"github.com/CAFxX/httpcompression/contrib/klauspost/zstd"
	kpzstd "github.com/klauspost/compress/zstd"
)

func Example() {
	// Create a compression adapter with default configuration
	compress, err := httpcompression.DefaultAdapter()
	if err != nil {
		log.Fatal(err)
	}
	// Define your handler, and apply the compression adapter.
	http.Handle("/", compress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})))
	// ...
}

func ExampleCustom() {
	brEnc, err := brotli.New(brotli.Options{})
	if err != nil {
		log.Fatal(err)
	}
	gzEnc, err := gzip.New(gzip.Options{})
	if err != nil {
		log.Fatal(err)
	}
	_, _ = httpcompression.Adapter(
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
	// Default zstd compressor
	zEnc, err := zstd.New()
	if err != nil {
		log.Fatal(err)
	}
	// zstd compressor with custom dictionary
	dict, coding, err := readZstdDictionary("tests/dictionary")
	if err != nil {
		log.Fatal(err)
	}
	zdEnc, err := zstd.New(kpzstd.WithEncoderDict(dict))
	if err != nil {
		log.Fatal(err)
	}
	_, _ = httpcompression.DefaultAdapter(
		// Add the zstd compressor with the dictionary.
		// We need to pick a custom content-encoding name. It is recommended to:
		// - avoid names that contain standard names (e.g. "gzip", "deflate", "br" or "zstd")
		// - include the dictionary ID, so that multiple dictionaries can be used (including
		//   e.g. multiple versions of the same dictionary)
		httpcompression.Compressor(coding, 3, zdEnc),
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

func readZstdDictionary(file string) (dict []byte, coding string, err error) {
	dictFile, err := os.Open(file)
	if err != nil {
		return nil, "", err
	}
	dict, err = ioutil.ReadAll(dictFile)
	if err != nil {
		return nil, "", err
	}
	if len(dict) < 8 {
		return nil, "", fmt.Errorf("invalid dictionary")
	}
	dictID := binary.LittleEndian.Uint32(dict[4:8]) // read the dictionary ID
	coding = fmt.Sprintf("z_%08x", dictID)          // build the encoding name: z_XXXXXXXX (where XXXXXXXX is the dictionary ID in hex lowercase)
	return
}
