Golang server middleware for HTTP compression
=============================================

[![Build status](https://github.com/CAFxX/httpcompression/workflows/Build/badge.svg)](https://github.com/CAFxX/httpcompression/actions)
[![codecov](https://codecov.io/gh/CAFxX/httpcompression/branch/master/graph/badge.svg)](https://codecov.io/gh/CAFxX/httpcompression)
[![Go Report](https://goreportcard.com/badge/github.com/CAFxX/httpcompression)](https://goreportcard.com/report/github.com/CAFxX/httpcompression)
[![CodeFactor](https://www.codefactor.io/repository/github/cafxx/httpcompression/badge)](https://www.codefactor.io/repository/github/cafxx/httpcompression)
[![Go Reference](https://pkg.go.dev/badge/github.com/CAFxX/httpcompression.svg)](https://pkg.go.dev/github.com/CAFxX/httpcompression) 

This is a small Go package which wraps HTTP handlers to transparently compress
response bodies using zstd, brotli or gzip - for clients which support them. Although 
it's usually simpler to leave that to a reverse proxy (like nginx or Varnish),
this package is useful when that is undesirable. In addition, this package allows
users to extend it by plugging in third-party or custom compression encoders.

**Note: This package was recently forked from NYTimes/gziphandler.
Maintaining drop-in compatibility is not a goal of this fork, as the scope of this fork
is significantly wider than the original package.**

:warning: As we have not reached 1.0 yet, API is still subject to changes.

## Features

- gzip and brotli compression by default, zstd and alternate (faster) gzip implementations are optional
- Apply compression only if response body size is greater than a threshold
- Apply compression only to a allowlist/denylist of MIME content types
- Define encoding priority (e.g. give brotli a higher priority than gzip)
- Control whether the client or the server defines the encoder priority
- Plug in third-party/custom compression schemes or implementations
- Custom dictionary compression for zstd

## Install
```bash
go get -u github.com/CAFxX/httpcompression
```

## Usage

Call `httpcompression.Handler` to get an adapter that can be used to wrap
any handler (an object which implements the `http.Handler` interface),
to transparently provide response body compression. 
Note that, despite the name, `httpcompression` automatically compresses using 
Brotli or Gzip, depending on the capabilities of the client (`Accept-Encoding`)
and the configuration of this handler (by default, both Gzip and Brotli are 
enabled and Brotli is used by default if the client supports both).

As a simple example:

```go
package main

import (
	"io"
	"net/http"
	"github.com/CAFxX/httpcompression"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Hello, World")
	})
	compress := httpcompression.DefaultAdapter() // Use the default configuration
	http.Handle("/", compress(handler))
	http.ListenAndServe("0.0.0.0:8080", nil)
}
```

## Benchmark

See the [benchmark results](results.md) to get an idea of the relative performance and
compression efficiency of gzip, brotli and zstd in the current implementation.

## TODO

- Add dictionary support to gzip and brotli (zstd already supports it)
- Allow to choose dictionary based on content-type

## License

[Apache 2.0][license].




[docs]:     https://godoc.org/github.com/CAFxX/httpcompression
[license]:  https://github.com/CAFxX/httpcompression/blob/master/LICENSE
