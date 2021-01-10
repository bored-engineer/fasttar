package fasttar

import (
	"bytes"
	"strconv"
)

/*
These are slight variations of the stdlib archive/tar versions

Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

func parseString(b []byte) []byte {
	if idx := bytes.IndexByte(b, 0); idx == 0 {
		return nil
	} else if idx >= 0 {
		return b[:idx]
	}
	return b
}

func parseOctal(b []byte) int64 {
	// Because unused fields are filled with NULs, we need
	// to skip leading NULs. Fields may also be padded with
	// spaces or NULs.
	// So we remove leading and trailing NULs and spaces to
	// be sure.
	b = bytes.Trim(b, " \x00")
	if len(b) == 0 {
		return 0
	}
	x, err := strconv.ParseUint(unsafeString(parseString(b)), 8, 64)
	if err != nil {
		return 0
	}
	return int64(x)
}

func parseNumeric(b []byte) int64 {
	// Check for base-256 (binary) format first.
	// If the first bit is set, then all following bits constitute a two's
	// complement encoded number in big-endian byte order.
	if len(b) > 0 && b[0]&0x80 != 0 {
		// Handling negative numbers relies on the following identity:
		//	-a-1 == ^a
		//
		// If the number is negative, we use an inversion mask to invert the
		// data bytes and treat the value as an unsigned number.
		var inv byte // 0x00 if positive or zero, 0xff if negative
		if b[0]&0x40 != 0 {
			inv = 0xff
		}

		var x uint64
		for i, c := range b {
			c ^= inv // Inverts c only if inv is 0xff, otherwise does nothing
			if i == 0 {
				c &= 0x7f // Ignore signal bit in first byte
			}
			if (x >> 56) > 0 {
				return 0 // Overflow
			}
			x = x<<8 | uint64(c)
		}
		if (x >> 63) > 0 {
			return 0 // Overflow
		}
		if inv == 0xff {
			return ^int64(x)
		}
		return int64(x)
	}

	// Normal case is base-8 (octal) format.
	return parseOctal(b)
}
