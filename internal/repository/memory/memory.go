package memory

import (
	"errors"
	"sync"

	"github.com/andream16/mitmcracker/internal/repository"
)

// InMemo represents an in-memory map.
type InMemo struct {
	KeyPairs sync.Map
}

// Insert
func (im *InMemo) Insert(key, cipherText string, mode repository.Mode) (*repository.KeyPair, bool, error) {
	if p, ok := im.KeyPairs.Load(cipherText); ok {
		pair, ok := p.(*repository.KeyPair)
		if ! ok {
			return nil, false, errors.New("could not convert found pair to *repository.KeyPain")
		}
		if mode == repository.EncodeMode {
			pair.EncodeKey = key
		} else {
			pair.DecodeKey = key
		}
		return pair, true, nil
	}

	newKeyPair := &repository.KeyPair{}

	if mode == repository.EncodeMode {
		newKeyPair.EncodeKey = key
	} else {
		newKeyPair.DecodeKey = key
	}

	im.KeyPairs.Store(cipherText, newKeyPair)

	return nil, false, nil
}
