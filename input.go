package mitmcracker

import (
	"flag"

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
	Storage   *Storage
}

// Storage contains storage information.
type Storage struct {
	// Type is the type of storage to be used.
	// Can only be:
	// - memory
	// - disk
	// If no option is provided, memory will be used by default.
	Type string
	// Address is the address of redis. It's optional.
	Address string
	// Password is the password of the redis database. It's optional.
	Password string
	// DB is the redis DB. It's optional.
	DB int
}

// New returns a new Input.
func New() (*Input, error) {

	in := &Input{}
	storage := &Storage{}

	flag.UintVar(&in.KeyLength, "key", 0, "key lenght")
	flag.StringVar(&in.PlainText, "plain", "", "known plain text")
	flag.StringVar(&in.EncText, "encoded", "", "known encoded text")
	flag.StringVar(&storage.Type, "storage.type", "", "storage to be used")
	flag.StringVar(&storage.Address, "storage.address", "", "storage address")
	flag.StringVar(&storage.Password, "storage.password", "", "storage password")
	flag.IntVar(&storage.DB, "storage.db", 0, "storage db")

	flag.Parse()

	in.Storage = storage

	return in, in.validate()

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
	if i.Storage != nil && i.Storage.Type != "" {
		_, ok := storages[i.Storage.Type]
		if !ok {
			return errors.Wrap(errInvalidField, "storage")
		}
	}
	return nil
}
