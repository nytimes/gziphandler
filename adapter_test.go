package httpcompression

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/CAFxX/httpcompression/contrib/klauspost/zstd"
	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
)

const (
	smallTestBody = "aaabbcaaabbbcccaaab"
	testBody      = "aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc"
)

func TestParseEncodings(t *testing.T) {
	t.Parallel()

	examples := map[string]codings{

		// Examples from RFC 2616
		"compress, gzip":                     {"compress": 1.0, "gzip": 1.0},
		"":                                   {},
		"*":                                  {"*": 1.0},
		"compress;q=0.5, gzip;q=1.0":         {"compress": 0.5, "gzip": 1.0},
		"gzip;q=1.0, identity; q=0.5, *;q=0": {"gzip": 1.0, "identity": 0.5, "*": 0.0},

		// More random stuff
		"AAA;q=1":     {"aaa": 1.0},
		"BBB ; q = 2": {"bbb": 1.0},
		"CCC; q = -1": {"ccc": 0.0},
		"DDD;":        {"ddd": 1.0},
		"EEE;;":       {"eee": 1.0},
		"FFF;q=;":     {"fff": 0.0},
		";":           {},
		";q=1":        {},
		";;":          {},
		";;q=1":       {},
	}

	for eg, exp := range examples {
		act := parseEncodings(eg)
		assert.Equal(t, exp, act)
	}
}

func TestGzipHandler(t *testing.T) {
	t.Parallel()

	// This just exists to provide something for GzipHandler to wrap.
	handler := newTestHandler(testBody)

	// requests without accept-encoding are passed along as-is
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, testBody, resp.Body.String())
	}

	// but requests with accept-encoding:gzip are compressed if possible
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, gzipStrLevel(testBody, gzip.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:br
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "br")
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip,br (br wins)
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip,br")
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip,br;q=0.5 (gzip wins)
	// because the server has no preference, we use the client preference (gzip)
	{
		c, _ := NewDefaultGzipCompressor(gzip.DefaultCompression)

		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip,br;q=0.5")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, Compressor(gzipEncoding, 1, c)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, gzipStrLevel(testBody, gzip.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip,br (br wins)
	// because the server has no preference, we use the client preference
	// becuase the client has no preference, we rely on the encoding name
	{
		c, _ := NewDefaultGzipCompressor(gzip.DefaultCompression)

		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip,br")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, Compressor(gzipEncoding, 1, c)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip,br and PreferClient (br wins)
	// because the client use q=1 for both, we rely on the server preference
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip,br")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, Prefer(PreferClient)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip,br and PreferClient (br wins)
	// because the client use q=1 for both, we rely on the server preference
	// becuase the server preference is the same, we rely on the encoding name
	{
		c, _ := NewDefaultGzipCompressor(gzip.DefaultCompression)

		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip,br")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, Prefer(PreferClient), Compressor(gzipEncoding, 1, c)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip,br;q=0.5 and PreferClient (gzip wins)
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip,br;q=0.5")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, Prefer(PreferClient)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, gzipStrLevel(testBody, gzip.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip;q=0.1,br;q=0.5 and PreferClient (br wins)
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip;q=0.1,br;q=0.5")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, Prefer(PreferClient)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip;q=0,br;q=0.5 and PreferClient (br wins)
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip;q=0,br;q=0.5")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, Prefer(PreferClient)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}

	// same, but with accept-encoding:gzip,br;q=0.5 and PreferServer (brotli wins)
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip,br;q=0.5")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, Prefer(PreferServer)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}

	// same, but disabling gzip (brotli wins)
	{
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip,br")
		resp := httptest.NewRecorder()
		newTestHandler(testBody, GzipCompressor(nil)).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, brotliStrLevel(testBody, brotli.DefaultCompression), resp.Body.Bytes())
	}
}

func TestGzipHandlerSmallBodyNoCompression(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(smallTestBody)

	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	res := resp.Result()

	// with less than 20 bytes the response should not be gzipped

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "", res.Header.Get("Content-Encoding"))
	assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
	assert.Equal(t, smallTestBody, resp.Body.String())

}

func TestGzipHandlerRepeatedCompressionGzip(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(testBody)
	testBodyBr := brotliStrLevel(testBody, brotli.DefaultCompression)

	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "br")
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		res := resp.Result()
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "br", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, testBodyBr, resp.Body.Bytes())
	}
}

func TestGzipHandlerRepeatedCompressionBrotli(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(testBody)
	testBodyGzip := gzipStrLevel(testBody, gzip.DefaultCompression)

	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		res := resp.Result()
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, testBodyGzip, resp.Body.Bytes())
	}
}

func TestGzipHandlerAlreadyCompressed(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(testBody)

	req, _ := http.NewRequest("GET", "/gzipped", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, testBody, res.Body.String())
}

func TestNewGzipLevelHandler(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, testBody)
	})

	for lvl := gzip.BestSpeed; lvl <= gzip.BestCompression; lvl++ {
		wrapper, err := DefaultAdapter(GzipCompressionLevel(lvl))
		if !assert.Nil(t, err, "NewGzipLevleHandler returned error for level:", lvl) {
			continue
		}

		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		resp := httptest.NewRecorder()
		wrapper(handler).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "Accept-Encoding", res.Header.Get("Vary"))
		assert.Equal(t, gzipStrLevel(testBody, lvl), resp.Body.Bytes())
	}
}

func TestNewGzipLevelHandlerReturnsErrorForInvalidLevels(t *testing.T) {
	t.Parallel()

	var err error
	_, err = DefaultAdapter(GzipCompressionLevel(-42))
	assert.NotNil(t, err)

	_, err = DefaultAdapter(GzipCompressionLevel(42))
	assert.NotNil(t, err)
}

func TestGzipHandlerNoBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode      int
		contentEncoding string
		emptyBody       bool
		body            []byte
	}{
		// Body must be empty.
		{http.StatusNoContent, "", true, nil},
		{http.StatusNotModified, "", true, nil},
		// Body is going to get gzip'd no matter what.
		{http.StatusOK, "", true, []byte{}},
		{http.StatusOK, "gzip", false, []byte(testBody)},
	}

	for num, test := range tests {
		mw, _ := DefaultAdapter()
		handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.statusCode)
			if test.body != nil {
				w.Write(test.body)
			}
		}))

		rec := httptest.NewRecorder()
		// TODO: in Go1.7 httptest.NewRequest was introduced this should be used
		// once 1.6 is not longer supported.
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
		assert.Equal(t, test.contentEncoding, header.Get("Content-Encoding"), fmt.Sprintf("for test iteration %d", num))
		assert.Equal(t, "Accept-Encoding", header.Get("Vary"), fmt.Sprintf("for test iteration %d", num))
		if test.emptyBody {
			assert.Empty(t, body, fmt.Sprintf("for test iteration %d", num))
		} else {
			assert.NotEmpty(t, body, fmt.Sprintf("for test iteration %d", num))
			assert.NotEqual(t, test.body, body, fmt.Sprintf("for test iteration %d", num))
		}
	}
}

func TestGzipHandlerContentLength(t *testing.T) {
	t.Parallel()

	testBodyBytes := []byte(testBody)
	tests := []struct {
		bodyLen   int
		bodies    [][]byte
		emptyBody bool
	}{
		{len(testBody), [][]byte{testBodyBytes}, false},
		// each of these writes is less than the DefaultMinSize
		{len(testBody), [][]byte{testBodyBytes[:200], testBodyBytes[200:]}, false},
		// without a defined Content-Length it should still gzip
		{0, [][]byte{testBodyBytes[:200], testBodyBytes[200:]}, false},
		// simulate a HEAD request with an empty write (to populate headers)
		{len(testBody), [][]byte{nil}, true},
	}

	// httptest.NewRecorder doesn't give you access to the Content-Length
	// header so instead, we create a server on a random port and make
	// a request to that instead
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("failed creating listen socket: %v", err)
	}
	defer ln.Close()
	srv := &http.Server{
		Handler: nil,
	}
	go srv.Serve(ln)

	for num, test := range tests {
		mw, _ := DefaultAdapter()
		srv.Handler = mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if test.bodyLen > 0 {
				w.Header().Set("Content-Length", strconv.Itoa(test.bodyLen))
			}
			for _, b := range test.bodies {
				w.Write(b)
			}
		}))
		req := &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/", Scheme: "http", Host: ln.Addr().String()},
			Header: make(http.Header),
			Close:  true,
		}
		req.Header.Set("Accept-Encoding", "gzip")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Unexpected error making http request in test iteration %d: %v", num, err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Unexpected error reading response body in test iteration %d: %v", num, err)
		}

		l, err := strconv.Atoi(res.Header.Get("Content-Length"))
		if err != nil {
			t.Fatalf("Unexpected error parsing Content-Length in test iteration %d: %v", num, err)
		}
		if test.emptyBody {
			assert.Empty(t, body, fmt.Sprintf("for test iteration %d", num))
			assert.Equal(t, 0, l, fmt.Sprintf("for test iteration %d", num))
		} else {
			assert.Len(t, body, l, fmt.Sprintf("for test iteration %d", num))
		}
		assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"), fmt.Sprintf("for test iteration %d", num))
		assert.NotEqual(t, test.bodyLen, l, fmt.Sprintf("for test iteration %d", num))
	}
}

func TestGzipHandlerMinSizeMustBePositive(t *testing.T) {
	t.Parallel()

	_, err := DefaultAdapter(MinSize(-1))
	assert.Error(t, err)
}

func TestGzipHandlerMinSize(t *testing.T) {
	t.Parallel()

	responseLength := 0
	b := []byte{'x'}

	wrapper, _ := DefaultAdapter(MinSize(128))
	handler := wrapper(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Write responses one byte at a time to ensure that the flush
			// mechanism, if used, is working properly.
			for i := 0; i < responseLength; i++ {
				n, err := w.Write(b)
				assert.Equal(t, 1, n)
				assert.Nil(t, err)
			}
		},
	))

	r, _ := http.NewRequest("GET", "/whatever", &bytes.Buffer{})
	r.Header.Add("Accept-Encoding", "gzip")

	// Short response is not compressed
	responseLength = 127
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Result().Header.Get(contentEncoding) == "gzip" {
		t.Error("Expected uncompressed response, got compressed")
	}

	// Long response is compressed
	responseLength = 128
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Result().Header.Get(contentEncoding) != "gzip" {
		t.Error("Expected compressed response, got uncompressed")
	}
}

type panicOnSecondWriteHeaderWriter struct {
	http.ResponseWriter
	headerWritten bool
}

func (w *panicOnSecondWriteHeaderWriter) WriteHeader(s int) {
	if w.headerWritten {
		panic("header already written")
	}
	w.headerWritten = true
	w.ResponseWriter.WriteHeader(s)
}

func TestGzipHandlerDoubleWriteHeader(t *testing.T) {
	t.Parallel()

	mw, _ := DefaultAdapter()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "15000")
		// Specifically write the header here
		w.WriteHeader(304)
		// Ensure that after a Write the header isn't triggered again on close
		w.Write(nil)
	}))
	wrapper := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = &panicOnSecondWriteHeaderWriter{
			ResponseWriter: w,
		}
		handler.ServeHTTP(w, r)
	})

	rec := httptest.NewRecorder()
	// TODO: in Go1.7 httptest.NewRequest was introduced this should be used
	// once 1.6 is not longer supported.
	req := &http.Request{
		Method:     "GET",
		URL:        &url.URL{Path: "/"},
		Proto:      "HTTP/1.1",
		ProtoMinor: 1,
		RemoteAddr: "192.0.2.1:1234",
		Header:     make(http.Header),
	}
	req.Header.Set("Accept-Encoding", "gzip")
	wrapper.ServeHTTP(rec, req)
	body, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %v", err)
	}
	assert.Empty(t, body)
	header := rec.Header()
	assert.Equal(t, "gzip", header.Get("Content-Encoding"))
	assert.Equal(t, "Accept-Encoding", header.Get("Vary"))
	assert.Equal(t, 304, rec.Code)
}

func TestGzipHandlerDoubleVary(t *testing.T) {
	t.Parallel()

	mw, _ := DefaultAdapter()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testBody))
	}))
	wrapper := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Add("Vary", "X-Something")
		handler.ServeHTTP(w, r)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	wrapper.ServeHTTP(rec, req)
	body, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %v", err)
	}
	assert.NotEmpty(t, body)
	header := rec.Header()
	assert.Equal(t, "gzip", header.Get("Content-Encoding"))
	// assert no duplicate Vary: Accept-Encoding
	assert.Equal(t, []string{"Accept-Encoding", "X-Something"}, header.Values("Vary"))
}

func TestStatusCodes(t *testing.T) {
	t.Parallel()

	mw, _ := DefaultAdapter()
	handler := mw(http.NotFoundHandler())
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	result := w.Result()
	if result.StatusCode != 404 {
		t.Errorf("StatusCode should have been 404 but was %d", result.StatusCode)
	}
}

func TestFlushBeforeWrite(t *testing.T) {
	t.Parallel()

	b := []byte(testBody)
	mw, _ := DefaultAdapter()
	handler := mw(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		rw.(http.Flusher).Flush()
		rw.Write(b)
	}))
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	assert.NotEqual(t, b, w.Body.Bytes())
}

func TestImplementCloseNotifier(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(acceptEncoding, "gzip")
	mw, _ := DefaultAdapter()
	mw(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, ok := rw.(http.CloseNotifier)
		assert.True(t, ok, "response writer must implement http.CloseNotifier")
	})).ServeHTTP(&mockRWCloseNotify{}, request)
}

func TestImplementFlusherAndCloseNotifier(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(acceptEncoding, "gzip")
	mw, _ := DefaultAdapter()
	mw(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, okCloseNotifier := rw.(http.CloseNotifier)
		assert.True(t, okCloseNotifier, "response writer must implement http.CloseNotifier")
		_, okFlusher := rw.(http.Flusher)
		assert.True(t, okFlusher, "response writer must implement http.Flusher")
	})).ServeHTTP(&mockRWCloseNotify{}, request)
}

func TestNotImplementCloseNotifier(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(acceptEncoding, "gzip")
	mw, _ := DefaultAdapter()
	mw(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, ok := rw.(http.CloseNotifier)
		assert.False(t, ok, "response writer must not implement http.CloseNotifier")
	})).ServeHTTP(httptest.NewRecorder(), request)
}

type mockRWCloseNotify struct{}

func (m *mockRWCloseNotify) CloseNotify() <-chan bool {
	panic("implement me")
}

func (m *mockRWCloseNotify) Header() http.Header {
	return http.Header{}
}

func (m *mockRWCloseNotify) Write([]byte) (int, error) {
	panic("implement me")
}

func (m *mockRWCloseNotify) WriteHeader(int) {
	panic("implement me")
}

func TestIgnoreSubsequentWriteHeader(t *testing.T) {
	t.Parallel()

	mw, _ := DefaultAdapter()
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.WriteHeader(404)
	}))
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	result := w.Result()
	if result.StatusCode != 500 {
		t.Errorf("StatusCode should have been 500 but was %d", result.StatusCode)
	}
}

func TestDontWriteWhenNotWrittenTo(t *testing.T) {
	t.Parallel()

	// When using gzip as middleware without ANY writes in the handler,
	// ensure the gzip middleware doesn't touch the actual ResponseWriter
	// either.

	mw, _ := DefaultAdapter()
	handler0 := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler0.ServeHTTP(w, r)
		w.WriteHeader(404) // this only works if gzip didn't do a WriteHeader(200)
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler1.ServeHTTP(w, r)

	result := w.Result()
	if result.StatusCode != 404 {
		t.Errorf("StatusCode should have been 404 but was %d", result.StatusCode)
	}
}

func TestContentTypes(t *testing.T) {
	t.Parallel()

	var contentTypeTests = []struct {
		name                 string
		contentType          string
		acceptedContentTypes []string
		expectedGzip         bool
	}{
		{
			name:                 "Always gzip when content types are empty",
			contentType:          "",
			acceptedContentTypes: []string{},
			expectedGzip:         true,
		},
		{
			name:                 "MIME match",
			contentType:          "application/json",
			acceptedContentTypes: []string{"application/json"},
			expectedGzip:         true,
		},
		{
			name:                 "MIME no match",
			contentType:          "text/xml",
			acceptedContentTypes: []string{"application/json"},
			expectedGzip:         false,
		},
		{
			name:                 "MIME match with no other directive ignores non-MIME directives",
			contentType:          "application/json; charset=utf-8",
			acceptedContentTypes: []string{"application/json"},
			expectedGzip:         true,
		},
		{
			name:                 "MIME match with other directives requires all directives be equal, different charset",
			contentType:          "application/json; charset=ascii",
			acceptedContentTypes: []string{"application/json; charset=utf-8"},
			expectedGzip:         false,
		},
		{
			name:                 "MIME match with other directives requires all directives be equal, same charset",
			contentType:          "application/json; charset=utf-8",
			acceptedContentTypes: []string{"application/json; charset=utf-8"},
			expectedGzip:         true,
		},
		{
			name:                 "MIME match with other directives requires all directives be equal, missing charset",
			contentType:          "application/json",
			acceptedContentTypes: []string{"application/json; charset=ascii"},
			expectedGzip:         false,
		},
		{
			name:                 "MIME match case insensitive",
			contentType:          "Application/Json",
			acceptedContentTypes: []string{"application/json"},
			expectedGzip:         true,
		},
		{
			name:                 "MIME match ignore whitespace",
			contentType:          "application/json;charset=utf-8",
			acceptedContentTypes: []string{"application/json;            charset=utf-8"},
			expectedGzip:         true,
		},
	}

	for _, tt := range contentTypeTests {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", tt.contentType)
			io.WriteString(w, testBody)
		})

		wrapper, err := DefaultAdapter(ContentTypes(tt.acceptedContentTypes, false))
		if !assert.Nil(t, err, "NewGzipHandlerWithOpts returned error", tt.name) {
			continue
		}

		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		resp := httptest.NewRecorder()
		wrapper(handler).ServeHTTP(resp, req)
		res := resp.Result()

		assert.Equal(t, 200, res.StatusCode)
		if tt.expectedGzip {
			assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"), tt.name)
		} else {
			assert.NotEqual(t, "gzip", res.Header.Get("Content-Encoding"), tt.name)
		}
	}
}

// --------------------------------------------------------------------

func BenchmarkAdapter(b *testing.B) {
	for _, size := range []int{10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			for _, ae := range []string{"gzip", "br", "zstd"} {
				b.Run(ae, func(b *testing.B) {
					b.Run("serial", func(b *testing.B) {
						benchmark(b, false, size, ae)
					})
					b.Run("parallel", func(b *testing.B) {
						benchmark(b, true, size, ae)
					})
				})
			}
		})
	}
}

// --------------------------------------------------------------------

func gzipStrLevel(s string, lvl int) []byte {
	var b bytes.Buffer
	w, _ := gzip.NewWriterLevel(&b, lvl)
	io.WriteString(w, s)
	w.Close()
	return b.Bytes()
}

func brotliStrLevel(s string, lvl int) []byte {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, lvl)
	io.WriteString(w, s)
	w.Close()
	return b.Bytes()
}

func benchmark(b *testing.B, parallel bool, size int, ae string) {
	bin, err := ioutil.ReadFile("testdata/benchmark.json")
	if err != nil {
		b.Fatal(err)
	}

	zenc, _ := zstd.New()

	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Set("Accept-Encoding", ae)
	handler := newTestHandler(
		string(bin[:size]),
		Compressor(zstd.Encoding, 2, zenc),
	)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if size < 20 {
		if res.Code != 200 || res.Header().Get("Content-Encoding") != "" || res.Body.Len() != size {
			b.Fatal(res)
		}
	} else {
		if res.Code != 200 || res.Header().Get("Content-Encoding") != ae || res.Body.Len() < size/10 {
			b.Fatal(res)
		}
	}

	b.ResetTimer()
	if parallel {
		b.RunParallel(func(pb *testing.PB) {
			res := &discardResponseWriter{}
			for pb.Next() {
				res.reset()
				handler.ServeHTTP(res, req)
			}
		})
	} else {
		res := &discardResponseWriter{}
		for i := 0; i < b.N; i++ {
			res.reset()
			handler.ServeHTTP(res, req)
		}
	}
}

type discardResponseWriter struct {
	h http.Header
}

func (w *discardResponseWriter) Header() http.Header {
	if w.h == nil {
		w.h = http.Header{}
	}
	return w.h
}

func (*discardResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (*discardResponseWriter) WriteHeader(int) {
}

func (w *discardResponseWriter) reset() {
	if w.h == nil {
		return
	}
	for k := range w.h {
		delete(w.h, k)
	}
}

func newTestHandler(body string, opts ...Option) http.Handler {
	mw, err := DefaultAdapter(opts...)
	if err != nil {
		panic(err)
	}
	buf := []byte(body)
	return mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/gzipped":
			w.Header().Set("Content-Encoding", "gzip")
			w.Write(buf)
		default:
			w.Write(buf)
		}
	}))
}
