package fasttar

import "unsafe"

// unsafeString performs an _unsafe_ no-copy string allocation from buf
// https://github.com/golang/go/issues/25484 has more info on this.
// The implementation is roughly taken from strings.Builder's
//
// This function is used internally to build encoding/xml elements
// without copying the underlying values on the assumption the
// original bytes slice given to NewScanner was immutable.
func unsafeString(buf []byte) string {
	if buf == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&buf))
}
