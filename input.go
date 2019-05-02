package mitmcracker

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

var (
	errInvalidField = errors.New("invalid_field")

	keyLengths = map[uint]struct{}{
		24: {},
		28: {},
		32: {},
	}
)

// Input represents the input that can be passed to the command.
type Input struct {
	// EncText is the encrypted word.
	EncText string
	// PlainText is the plain word.
	PlainText string
	// KeyLength is the length of the key.
	// Can only be:
	//  - 24-bit
	//  - 28-bit
	//  - 32-bit
	KeyLength uint
}

// New returns a new Input.
func New() (*Input, error) {

	in := &Input{}

	flag.UintVar(&in.KeyLength, "k", 0, "key lenght")
	flag.StringVar(&in.PlainText, "p", "", "known plain text")
	flag.StringVar(&in.EncText, "e", "", "known encoded text")

	flag.Parse()

	return in, in.validate()

}

// String returns the string representation of an input.
func (i *Input) String() string {
	return fmt.Sprintf("key: %d, plain-text: %s, encoded-text: %s", i.KeyLength, i.PlainText, i.EncText)
}

func (i *Input) validate() error {
	_, ok := keyLengths[i.KeyLength]
	if !ok {
		return errors.Wrap(errInvalidField, "key length")
	}
	if i.EncText == "" {
		return errors.Wrap(errInvalidField, "known encrypted text")
	}
	if i.PlainText == "" {
		return errors.Wrap(errInvalidField, "known encrypted text")
	}
	return nil
}
