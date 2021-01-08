package fasttar

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	var buf bytes.Buffer
	w := tar.NewWriter(&buf)
	w.WriteHeader(&tar.Header{
		Name: "foo",
	})
	w.WriteHeader(&tar.Header{
		Name: "bar",
		Size: 3,
	})
	w.Write([]byte("har"))
	w.WriteHeader(&tar.Header{
		Name: "far",
	})
	w.Close()
	r := NewReader(buf.Bytes())
	// foo
	header, file, err := r.Next()
	assert.NoError(t, err)
	assert.Len(t, header, 512)
	assert.Len(t, file, 0)
	assert.NoError(t, err)
	assert.Equal(t, "foo", NameString(header))
	// bar
	header, file, err = r.Next()
	assert.NoError(t, err)
	assert.Len(t, header, 512)
	assert.Len(t, file, 3)
	assert.NoError(t, err)
	assert.Equal(t, "bar", NameString(header))
	assert.Equal(t, []byte("har"), file)
	// far
	header, file, err = r.Next()
	assert.NoError(t, err)
	assert.Len(t, header, 512)
	assert.Len(t, file, 0)
	assert.NoError(t, err)
	assert.Equal(t, "far", NameString(header))
	// EOF
	header, file, err = r.Next()
	assert.Equal(t, io.EOF, err)
	assert.Len(t, header, 512)
	assert.Nil(t, file)
}
