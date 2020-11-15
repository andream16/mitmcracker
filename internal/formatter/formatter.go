package formatter

import "fmt"

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
