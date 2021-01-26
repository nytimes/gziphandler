package pgzip_test

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"testing"

	stdgzip "compress/gzip"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/klauspost/pgzip"
)

var _ httpcompression.CompressorProvider = &pgzip.Compressor{}

func TestPgzip(t *testing.T) {
	t.Parallel()

	s := []byte("hello world!")

	c, err := pgzip.New(pgzip.Options{BlockSize: 1 << 20, Blocks: runtime.GOMAXPROCS(0)})
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
