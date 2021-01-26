package brotli_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"

	_brotli "github.com/andybalholm/brotli"
)

var _ httpcompression.CompressorProvider = &brotli.Compressor{}

func TestBrotli(t *testing.T) {
	t.Parallel()

	s := []byte("hello world!")

	c, err := brotli.New(brotli.Options{})
	if err != nil {
		t.Fatal(err)
	}
	b := &bytes.Buffer{}
	w := c.Get(b)
	w.Write(s)
	w.Close()

	r := _brotli.NewReader(b)
	d, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(s, d) != 0 {
		t.Fatalf("decoded string mismatch\ngot: %q\nexp: %q", string(s), string(d))
	}
}
