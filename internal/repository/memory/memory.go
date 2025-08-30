package memory

import (
	"context"
	"sync"

	"github.com/andream16/mitmcracker/internal/repository"
)

// InMemo represents an in-memory map with thread-safe operations.
type InMemo struct {
	KeyPairs map[string]*repository.KeyPair
	mu       sync.RWMutex
}

// Insert adds a key-ciphertext pair to the repository.
func (im *InMemo) Insert(ctx context.Context, key, cipherText, mode string) (*repository.KeyPair, bool, error) {
	select {
	case <-ctx.Done():
		return nil, false, ctx.Err()
	default:
	}

	im.mu.Lock()
	defer im.mu.Unlock()
	if p, ok := im.KeyPairs[cipherText]; ok {
		if mode == repository.EncodeMode {
			p.EncodeKey = key
		} else {
			p.DecodeKey = key
		}
		return p, true, nil
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

func (im *InMemo) InsertBulk(ctx context.Context, reqs ...repository.InsertBulkRequest) (*repository.KeyPair, bool, error) {
	select {
	case <-ctx.Done():
		return nil, false, ctx.Err()
	default:
	}

	im.mu.Lock()
	defer im.mu.Unlock()
	for _, r := range reqs {
		if p, ok := im.KeyPairs[r.CipherText]; ok {
			if r.Mode == repository.EncodeMode {
				p.EncodeKey = r.Key
			} else {
				p.DecodeKey = r.Key
			}
			return p, true, nil
		}

		newKeyPair := &repository.KeyPair{}

		if r.Mode == repository.EncodeMode {
			newKeyPair.EncodeKey = r.Key
		} else {
			newKeyPair.DecodeKey = r.Key
		}

		im.KeyPairs[r.CipherText] = newKeyPair
	}

	return nil, false, nil
}
