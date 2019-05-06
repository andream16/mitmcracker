package repository

// Repositorer represents the repository interface.
type Repositorer interface {
	InsertEnc(key, cipherText string)
	InsertDec(key, cipherText string)
	FindKey() string
}
