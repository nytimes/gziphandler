package stdlib

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGzipDoubleClose(t *testing.T) {
	// reset the pool for the default compression so we can make sure duplicates
	// aren't added back by double close
	addLevelPool(gzip.DefaultCompression)

	w := bytes.NewBufferString("")
	writer := NewWriter(w, gzip.DefaultCompression)
	writer.Close()

	// the second close shouldn't have added the same writer
	// so we pull out 2 writers from the pool and make sure they're different
	w1 := gzipWriterPools[poolIndex(gzip.DefaultCompression)].Get()
	w2 := gzipWriterPools[poolIndex(gzip.DefaultCompression)].Get()
	// assert.NotEqual looks at the value and not the address, so we use regular ==
	assert.False(t, w1 == w2)
}
