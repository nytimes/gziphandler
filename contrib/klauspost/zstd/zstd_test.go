package zstd_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/klauspost/zstd"
	kpzstd "github.com/klauspost/compress/zstd"
)

var _ httpcompression.CompressorProvider = &zstd.Compressor{}

func TestZstd(t *testing.T) {
	t.Parallel()

	s := []byte("hello world!")

	c, err := zstd.New()
	if err != nil {
		t.Fatal(err)
	}
	b := &bytes.Buffer{}
	w := c.Get(b)
	w.Write(s)
	w.Close()

	r, err := kpzstd.NewReader(b)
	if err != nil {
		t.Fatal(err)
	}
	d, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(s, d) != 0 {
		t.Fatalf("decoded string mismatch\ngot: %q\nexp: %q", string(s), string(d))
	}
}
