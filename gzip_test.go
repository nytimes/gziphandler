package gziphandler

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEncodings(t *testing.T) {

	examples := map[string]codings{

		// Examples from RFC 2616
		"compress, gzip": codings{"compress": 1.0, "gzip": 1.0},
		"":               codings{},
		"*":              codings{"*": 1.0},
		"compress;q=0.5, gzip;q=1.0":         codings{"compress": 0.5, "gzip": 1.0},
		"gzip;q=1.0, identity; q=0.5, *;q=0": codings{"gzip": 1.0, "identity": 0.5, "*": 0.0},

		// More random stuff
		"AAA;q=1":     codings{"aaa": 1.0},
		"BBB ; q = 2": codings{"bbb": 1.0},
	}

	for eg, exp := range examples {
		act, _ := parseEncodings(eg)
		assert.Equal(t, exp, act)
	}
}

func TestGzipHandler(t *testing.T) {
	testBody := "aaabbbccc"

	// This just exists to provide something for GzipHandler to wrap.
	handler := newTestHandler(testBody)

	// requests without accept-encoding are passed along as-is

	req1, _ := http.NewRequest("GET", "/whatever", nil)
	res1 := httptest.NewRecorder()
	handler.ServeHTTP(res1, req1)

	assert.Equal(t, 200, res1.Code)
	assert.Equal(t, "", res1.Header().Get("Content-Encoding"))
	assert.Equal(t, "Accept-Encoding", res1.Header().Get("Vary"))
	assert.Equal(t, testBody, res1.Body.String())

	// but requests with accept-encoding:gzip are compressed if possible

	req2, _ := http.NewRequest("GET", "/whatever", nil)
	req2.Header.Set("Accept-Encoding", "gzip")
	res2 := httptest.NewRecorder()
	handler.ServeHTTP(res2, req2)

	assert.Equal(t, 200, res2.Code)
	assert.Equal(t, "gzip", res2.Header().Get("Content-Encoding"))
	assert.Equal(t, "Accept-Encoding", res2.Header().Get("Vary"))
	assert.Equal(t, gzipStrLevel(testBody, gzip.DefaultCompression), res2.Body.Bytes())

	// content-type header is correctly set based on uncompressed body

	req3, _ := http.NewRequest("GET", "/whatever", nil)
	req3.Header.Set("Accept-Encoding", "gzip")
	res3 := httptest.NewRecorder()
	handler.ServeHTTP(res3, req3)

	assert.Equal(t, http.DetectContentType([]byte(testBody)), res3.Header().Get("Content-Type"))
}

func TestNewGzipLevelHandler(t *testing.T) {
	testBody := "aaabbbccc"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, testBody)
	})

	for lvl := gzip.BestSpeed; lvl <= gzip.BestCompression; lvl++ {
		wrapper, err := NewGzipLevelHandler(lvl)
		if !assert.Nil(t, err, "NewGzipLevleHandler returned error for level:", lvl) {
			continue
		}

		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		res := httptest.NewRecorder()
		wrapper(handler).ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "gzip", res.Header().Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header().Get("Vary"))
		assert.Equal(t, gzipStrLevel(testBody, lvl), res.Body.Bytes())

	}
}

func TestNewGzipLevelHandlerReturnsErrorForInvalidLevels(t *testing.T) {
	var err error
	_, err = NewGzipLevelHandler(-42)
	assert.NotNil(t, err)

	_, err = NewGzipLevelHandler(42)
	assert.NotNil(t, err)
}

func TestMustNewGzipLevelHandlerWillPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic was not called")
		}
	}()

	_ = MustNewGzipLevelHandler(-42)
}

func TestGzipHandlerNoBody(t *testing.T) {
	tests := []struct {
		statusCode int
	}{
		{http.StatusOK}, // Can contain a body.
		// Body must be empty.
		{http.StatusNoContent},
		{http.StatusNotModified},
	}

	for _, test := range tests {
		handler := GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.statusCode)
		}))

		rec := httptest.NewRecorder()
		// TODO: in Go1.7 httptest.NewRequest was introduced this should be used
		// once 1.7 is not longer supported.
		req := &http.Request{
			Method:     "GET",
			URL:        &url.URL{Path: "/"},
			Proto:      "HTTP/1.1",
			ProtoMinor: 1,
			RemoteAddr: "192.0.2.1:1234",
			Header:     make(http.Header),
		}
		req.Header.Set("Accept-Encoding", "gzip")
		handler.ServeHTTP(rec, req)

		body, err := ioutil.ReadAll(rec.Body)
		if err != nil {
			t.Fatalf("Unexpected error reading response body: %v", err)
		}

		header := rec.Header()
		assert.Equal(t, "", header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", header.Get("Vary"))
		assert.Equal(t, 0, len(body))
	}
}

// --------------------------------------------------------------------

func BenchmarkGzipHandler_S2k(b *testing.B)   { benchmark(b, false, 2048) }
func BenchmarkGzipHandler_S20k(b *testing.B)  { benchmark(b, false, 20480) }
func BenchmarkGzipHandler_S100k(b *testing.B) { benchmark(b, false, 102400) }
func BenchmarkGzipHandler_P2k(b *testing.B)   { benchmark(b, true, 2048) }
func BenchmarkGzipHandler_P20k(b *testing.B)  { benchmark(b, true, 20480) }
func BenchmarkGzipHandler_P100k(b *testing.B) { benchmark(b, true, 102400) }

// --------------------------------------------------------------------

func gzipStrLevel(s string, lvl int) []byte {
	var b bytes.Buffer
	w, _ := gzip.NewWriterLevel(&b, lvl)
	io.WriteString(w, s)
	w.Close()
	return b.Bytes()
}

func benchmark(b *testing.B, parallel bool, size int) {
	bin, err := ioutil.ReadFile("testdata/benchmark.json")
	if err != nil {
		b.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	handler := newTestHandler(string(bin[:size]))

	if parallel {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				runBenchmark(b, req, handler)
			}
		})
	} else {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			runBenchmark(b, req, handler)
		}
	}
}

func runBenchmark(b *testing.B, req *http.Request, handler http.Handler) {
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if code := res.Code; code != 200 {
		b.Fatalf("Expected 200 but got %d", code)
	} else if blen := res.Body.Len(); blen < 500 {
		b.Fatalf("Expected complete response body, but got %d bytes", blen)
	}
}

func newTestHandler(body string) http.Handler {
	return GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, body)
	}))
}
