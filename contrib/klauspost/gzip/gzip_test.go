package gzip_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	stdgzip "compress/gzip"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/klauspost/gzip"
)

var _ httpcompression.CompressorProvider = &gzip.Compressor{}

func TestGzip(t *testing.T) {
	t.Parallel()

	s := []byte("hello world!")

	c, err := gzip.New(gzip.Options{})
	if err != nil {
		t.Fatal(err)
	}
	b := &bytes.Buffer{}
	w := c.Get(b)
	w.Write(s)
	w.Close()

	r, err := stdgzip.NewReader(b)
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
