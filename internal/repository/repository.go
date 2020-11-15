package repository

import "fmt"

const (
	EncodeMode Mode = "encode"
	DecodeMode Mode = "decode"
)

// Inserter represents the repository interface.
type Inserter interface {
	Insert(key, cipherText string, mode Mode) (*KeyPair, bool, error)
}

// KeyPair is the final result of FindKeys.
type KeyPair struct {
	EncodeKey string
	DecodeKey string
}

// Mode represents the operation mode.
type Mode string

// Validate validates the mode.
func (m Mode) Validate() error {
	if _, ok := map[Mode]struct{}{
		EncodeMode: {},
		DecodeMode: {},
	}[m]; !ok {
		return fmt.Errorf("invalid mode %s. Should be either %s or %s", m, EncodeMode, DecodeMode)
	}
	return nil
}
