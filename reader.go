package fasttar

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
)

var (
	errHeader = errors.New("invalid tar header")
)

var zeroBlock [512]byte

// Reader reads a tar file
type Reader struct {
	buf []byte
	pos int
}

// Cap returns the capacity of the readers's underlying byte slice.
func (r *Reader) Cap() int {
	return cap(r.buf)
}

// Len returns the number of bytes stored in the reader.
func (r *Reader) Len() int {
	return len(r.buf)
}

// Truncate resets the buffer to be empty, but it retains the underlying storage for use by future writes.
func (r *Reader) Truncate() {
	r.buf = r.buf[:0]
	r.pos = 0
}

// Write implement io.Writer by appending data to the internal slice
func (r *Reader) Write(p []byte) (int, error) {
	r.buf = append(r.buf, p...)
	return len(p), nil
}

// ReadFrom implements io.ReaderFrom, _replacing_ the internal buffer with the read bytes
func (r *Reader) ReadFrom(reader io.Reader) (int64, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return int64(len(buf)), err
	}
	r.buf = buf
	r.pos = 0
	return int64(len(buf)), nil
}

// readHeader checks if we will read past the end
func (r *Reader) readHeader() ([]byte, error) {
	start := r.pos
	r.pos += BlockSize
	if r.pos > len(r.buf) {
		return nil, io.ErrUnexpectedEOF
	}
	return r.buf[start:r.pos], nil
}

// Next reads the next file from the archive
func (r *Reader) Next() (header []byte, file []byte, err error) {
	// Read the header
	header, err = r.readHeader()
	if err != nil {
		return
	}
	// If it was zero blocks, expect EOF
	if bytes.Equal(header, zeroBlock[:]) {
		header, err = r.readHeader()
		if err != nil {
			return
		}
		if !bytes.Equal(header, zeroBlock[:]) {
			err = errHeader
			return
		}
		err = io.EOF
		return
	}
	// Read the file size from the header
	size := int(Size(header))
	offset := r.pos
	r.pos += size
	if r.pos > len(r.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	file = r.buf[offset:r.pos]
	if rem := (size % BlockSize); rem > 0 {
		// Move position past the end of the block
		r.pos += BlockSize - rem
	}
	return
}

// Reset sets the internal buffer to buf
func (r *Reader) Reset(buf []byte) {
	r.buf = buf
	r.pos = 0
}

// NewReader creates a new *Reader
func NewReader(buf []byte) *Reader {
	return &Reader{buf: buf, pos: 0}
}
