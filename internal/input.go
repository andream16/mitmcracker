package mitmcracker

import "github.com/pkg/errors"

// input represents the input that can be passed to the command.
type input struct {
	// encText is the encrypted word.
	encText string
	// plainText is the plain word.
	plainText string
	// keyLength is the length of the key.
	// Can only be:
	//  - 24-bit
	//  - 28-bit
	//  - 32-bit
	keyLength uint
}

var (
	errInvalidField = errors.New("invalid_field")

	keyLengths = map[uint]struct{}{
		24: {},
		28: {},
		32: {},
	}
)

func (i *input) validate() error {
	_, ok := keyLengths[i.keyLength]
	if !ok {
		return errors.Wrap(errInvalidField, "key length")
	}
	if i.encText == "" {
		return errors.Wrap(errInvalidField, "known encrypted text")
	}
	if i.plainText == "" {
		return errors.Wrap(errInvalidField, "known encrypted text")
	}
	return nil
}
