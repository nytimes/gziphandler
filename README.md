Gzip (and Brotli) Handler
============

[![Documentation](https://godoc.org/github.com/CAFxX/gziphandler?status.svg)](https://godoc.org/github.com/CAFxX/gziphandler)
[![Coverage](https://gocover.io/_badge/github.com/CAFxX/gziphandler)](https://gocover.io/github.com/CAFxX/gziphandler)

This is a tiny Go package which wraps HTTP handlers to transparently compress
response bodies, using brotli or gzip, for clients which support it. Although 
it's usually simpler to leave that to a reverse proxy (like nginx or Varnish),
this package is useful when that's undesirable.

**Note: This package was recently forked from NYTimes/gziphandler, so this is where
the name comes from. Since maintaining drop-in compatibility is not a goal of this
fork, and since the scope of the fork is wider than the original package, this
package will likely be renamed in the near future.**

## Install
```bash
go get -u github.com/CAFxX/gziphandler
```

## Usage

Call `GzipHandler` with any handler (an object which implements the
`http.Handler` interface), and it'll return a new handler which gzips
the response. Note that, despite the name, `GzipHandler` automatically
compresses using Brotli or Gzip, depending on the capabilities of the
client (`Accept-Encoding`) and the configuration of this handler (by
default, both Gzip and Brotli are enabled and, unless the client prefers
Gzip over Brotli, Brotli is used by default).

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


## Documentation

The docs can be found at [godoc.org][docs], as usual.


## License

[Apache 2.0][license].




[docs]:     https://godoc.org/github.com/CAFxX/gziphandler
[license]:  https://github.com/CAFxX/gziphandler/blob/master/LICENSE
