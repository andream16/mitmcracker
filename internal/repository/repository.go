package repository

import "errors"

// ErrNotFound is returned when no common cipher-text is found.
var ErrNotFound = errors.New("key_not_found")

// Repositorer represents the repository interface.
type Repositorer interface {
	InsertEnc(key, cipherText string) error
	InsertDec(key, cipherText string) error
	FindKeys() (*Keys, error)
}

// Keys is the final result of FindKeys.
type Keys struct {
	EncKey string
	DecKey string
}
