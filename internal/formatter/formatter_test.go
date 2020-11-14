package formatter_test

import (
	"testing"

	"github.com/andream16/mitmcracker/internal/formatter"

	"github.com/stretchr/testify/assert"
)

func TestDefaultFormatter(t *testing.T) {
	t.Run("it formats correctly values based on key length", func(t *testing.T) {
		for _, tc := range []struct {
			key       int
			keyLength uint
			expected  string
		}{
			{key: 1, keyLength: 24, expected: "000001"},
			{key: 16777215, keyLength: 24, expected: "ffffff"},
			{key: 1, keyLength: 28, expected: "0000001"},
			{key: 268435455, keyLength: 28, expected: "fffffff"},
			{key: 1, keyLength: 32, expected: "00000001"},
			{key: 4294967295, keyLength: 32, expected: "ffffffff"},
		} {
			assert.Equal(t, tc.expected, formatter.DefaultFormatter(tc.key, tc.keyLength))
		}
	})
}
