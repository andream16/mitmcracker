package repository

import "context"

const (
	EncodeMode = "encode"
	DecodeMode = "decode"
)

// Inserter represents the repository interface.
type Inserter interface {
	Insert(ctx context.Context, key, cipherText, mode string) (*KeyPair, bool, error)
	InsertBulk(ctx context.Context, reqs ...InsertBulkRequest) (*KeyPair, bool, error)
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
