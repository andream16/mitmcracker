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

	storages = map[string]struct{}{
		"memory": {},
		"disk":   {},
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
	// Storage is the type of storage to be used.
	// Can only be:
	// - memory
	// - disk
	// If no option is provided, memory will be used by default.
	Storage string
}

// New returns a new Input.
func New() (*Input, error) {

	in := &Input{}

	flag.UintVar(&in.KeyLength, "key", 0, "key lenght")
	flag.StringVar(&in.PlainText, "plain", "", "known plain text")
	flag.StringVar(&in.EncText, "encoded", "", "known encoded text")
	flag.StringVar(&in.Storage, "storage", "", "storage to be used")

	flag.Parse()

	return in, in.validate()

}

// String returns the string representation of an input.
func (i *Input) String() string {
	return fmt.Sprintf("key: %d, plain-text: %s, encoded-text: %s, storage: %s", i.KeyLength, i.PlainText, i.EncText, i.Storage)
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
		return errors.Wrap(errInvalidField, "known plain text")
	}
	if i.Storage != "" {
		_, ok := storages[i.Storage]
		if !ok {
			return errors.Wrap(errInvalidField, "storage")
		}
	}
	return nil
}
