package formatter

import (
	"fmt"
	"strconv"
)

// Formatter formats a given key based on the key length.
type Formatter func(key int, keyLength uint) string

// DefaultFormatter formats a given key to 24, 28 or 32 bits.
func DefaultFormatter(key int, keyLength uint) string {
	var s string
	switch keyLength {
	case 24:
		s = "%06x"
	case 28:
		s = "%07x"
	case 32:
		s = "%08x"
	}
	return fmt.Sprintf(s, key)
}

// FastFormatter >> fast.
func FastFormatter(key int, keyLength uint) string {
	var kL int
	switch keyLength {
	case 24:
		kL = 6
	case 28:
		kL = 7
	case 32:
		kL = 8
	}
	s := strconv.AppendUint(
		[]byte{'0', '0', '0', '0', '0', '0', '0'},
		uint64(key),
		16,
	)
	return string(s[len(s)-kL:])
}
