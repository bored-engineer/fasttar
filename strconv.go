package fasttar

import (
	"bytes"
	"strconv"
)

/* These are slight variations of the stdlib archive/tar versions */

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
