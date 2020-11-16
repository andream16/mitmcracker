package repository

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