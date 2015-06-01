package gziphandler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	handler := GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, testBody)
	}))

	// requests without accept-encoding are passed along as-is

	req1, _ := http.NewRequest("GET", "/whatever", nil)
	res1 := httptest.NewRecorder()
	handler.ServeHTTP(res1, req1)

	assert.Equal(t, 200, res1.Code)
	assert.Equal(t, "", res1.Header().Get("Content-Encoding"))
	assert.Equal(t, testBody, res1.Body.String())

	// but requests with accept-encoding:gzip are compressed if possible

	req2, _ := http.NewRequest("GET", "/whatever", nil)
	req2.Header.Set("Accept-Encoding", "gzip")
	res2 := httptest.NewRecorder()
	handler.ServeHTTP(res2, req2)

	assert.Equal(t, 200, res2.Code)
	assert.Equal(t, "gzip", res2.Header().Get("Content-Encoding"))
	assert.Equal(t, gzipStr(testBody), res2.Body.Bytes())
}

func gzipStr(s string) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	io.WriteString(w, s)
	w.Close()
	return b.Bytes()
}

func TestMultiError(t *testing.T) {
	tests := []string{
		"",
		"compress;q=0.....1, gzip;q=1.r",
		"gzip;q=1.0, identity; q=///, *;q=[",

		"AAA;q=r",
		"BBB ; q = f",
	}

	want := ErrorList{
		KeyError{"", ErrEmptyContentCoding},
		KeyError{"compress;q=0.....1", &strconv.NumError{
			Func: "ParseFloat",
			Num:  "0.....1",
			Err:  strconv.ErrSyntax,
		}},
		KeyError{" gzip;q=1.r", &strconv.NumError{
			Func: "ParseFloat",
			Num:  "1.r",
			Err:  strconv.ErrSyntax,
		}},
		KeyError{" identity; q=///", &strconv.NumError{
			Func: "ParseFloat",
			Num:  "///",
			Err:  strconv.ErrSyntax,
		}},
		KeyError{" *;q=[", &strconv.NumError{
			Func: "ParseFloat",
			Num:  "[",
			Err:  strconv.ErrSyntax,
		}},
		KeyError{"AAA;q=r", &strconv.NumError{
			Func: "ParseFloat",
			Num:  "r",
			Err:  strconv.ErrSyntax,
		}},
	}

	masterList := new(ErrorList)
	for _, eg := range tests {
		_, errList := parseEncodings(eg)

		masterList.Append(errList)
	}

	if !equalErrorSlice(*masterList, want) {
		t.Errorf("Slices do not match %+v\n\n%+v", *masterList, want)
	}

	assert.Equal(t, "6 errors", masterList.Error())
}

func equalErrorSlice(s1, s2 ErrorList) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		e1 := s1[i]
		e2 := s2[i]

		if e1.Key != e2.Key ||
			// If the string outputs match, the errors are the same.
			// There's no need for deep-equality checking, particularly
			// because the error could be something like this:
			// gziphandler.KeyError{Key:"compress;q=0.....1", Err:(*strconv.NumError)(0xc208038d80)}
			// and it's pointless to chase the (potential) pointer.
			e1.Err.Error() != e2.Err.Error() {
			return false
		}
	}

	return true
}
