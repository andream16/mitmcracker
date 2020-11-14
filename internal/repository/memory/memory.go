package memory

import (
	"sync"

	"github.com/andream16/mitmcracker/internal/repository"
)

// InMemo represents an in-memory map.
type InMemo struct {
	sync.Mutex

	KeyPairs map[string]*repository.KeyPair
}

// Insert
func (im *InMemo) Insert(key, cipherText string, mode repository.Mode) (*repository.KeyPair, bool, error) {
	im.Lock()
	defer im.Unlock()

	if err := mode.Validate(); err != nil {
		return nil, false, err
	}

	if len(im.KeyPairs) == 0 {
		im.KeyPairs = map[string]*repository.KeyPair{}
	}

	if pair, ok := im.KeyPairs[cipherText]; ok {
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

	im.KeyPairs[cipherText] = newKeyPair

	return nil, false, nil
}
