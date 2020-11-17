package repository

const (
	EncodeMode = "encode"
	DecodeMode = "decode"
)

// Inserter represents the repository interface.
type Inserter interface {
	Insert(key, cipherText, mode string) (*KeyPair, bool, error)
	InsertBulk(reqs ...InsertBulkRequest) (*KeyPair, bool, error)
}

// InsertBulkRequest represents an InsertBulk request.
type InsertBulkRequest struct {
	Key        string
	CipherText string
	Mode       string
}

// KeyPair is the final result of FindKeys.
type KeyPair struct {
	EncodeKey string
	DecodeKey string
}
