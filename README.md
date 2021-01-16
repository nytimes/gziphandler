Golang server middleware for HTTP compression
=============================================

[![Documentation](https://godoc.org/github.com/CAFxX/gziphandler?status.svg)](https://godoc.org/github.com/CAFxX/gziphandler)
[![Coverage](https://gocover.io/_badge/github.com/CAFxX/gziphandler)](https://gocover.io/github.com/CAFxX/gziphandler)

This is a small Go package which wraps HTTP handlers to transparently compress
response bodies, using zstd, brotli or gzip, for clients which support it. Although 
it's usually simpler to leave that to a reverse proxy (like nginx or Varnish),
this package is useful when that's undesirable. In addition, this package allows
users to extend it by providing third-party or custom compression encodings.

**Note: This package was recently forked from NYTimes/gziphandler, so this is where
the name comes from. Since maintaining drop-in compatibility is not a goal of this
fork, and since the scope of the fork is wider than the original package, this
package will likely be renamed in the near future.**

## Features

- Apply compression only if response body size is greater than a threshold
- Apply compression only to a allowlist/denylist of MIME content types
- Define encoding priority (e.g. give brotli a higher priority than gzip)
- Control whether the client or the server defines the encoder priority
- Plug in third-party/custom compression schemes or implementations
- Custom dictionary compression for zstd

## Install
```bash
go get -u github.com/CAFxX/gziphandler
```

## Usage

Call `GzipHandler` with any handler (an object which implements the
`http.Handler` interface), and it will return a new handler which gzips
the response. Note that, despite the name, `GzipHandler` automatically
compresses using Brotli or Gzip, depending on the capabilities of the
client (`Accept-Encoding`) and the configuration of this handler (by
default, both Gzip and Brotli are enabled and Brotli is used by default
if the client supports both).

As a simple example:

```go
package main

import (
	"io"
	"net/http"
	"github.com/CAFxX/gziphandler"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Hello, World")
	})

	compressHandler := gziphandler.GzipHandler(handler)

	http.Handle("/", compressHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)
}
```

## TODO

- Add dictionary support to gzip and brotli
- Allow to choose dictionary based on content-type

## License

[Apache 2.0][license].




[docs]:     https://godoc.org/github.com/CAFxX/gziphandler
[license]:  https://github.com/CAFxX/gziphandler/blob/master/LICENSE
