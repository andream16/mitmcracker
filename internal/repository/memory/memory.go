package memory

import (
	"sync"

	"github.com/andream16/mitmcracker/internal/repository"
)

var _ repository.Repositorer = (*InMemo)(nil)

// InMemo represents an in-memory map.
type InMemo struct {
	sync.Mutex
	Dec map[string]string
	Enc map[string]string
}

// New returns a new *InMemo with fixed size maps.
func New(size int) *InMemo {
	return &InMemo{
		Dec: make(map[string]string, size),
		Enc: make(map[string]string, size),
	}
}

// InsertDec adds a new entry in dec map.
func (im *InMemo) InsertDec(key, cipherText string) error {

	im.Lock()
	defer im.Unlock()

	if len(im.Dec) == 0 {
		im.Dec = map[string]string{}
	}

	im.Dec[cipherText] = key
	return nil

}

// InsertEnc adds a new entry in enc map.
func (im *InMemo) InsertEnc(key, cipherText string) error {

	im.Lock()
	defer im.Unlock()

	if len(im.Enc) == 0 {
		im.Enc = map[string]string{}
	}

	im.Enc[cipherText] = key
	return nil

}

// FindKeys returns the keys used to encode and decode the message.
func (im *InMemo) FindKeys() (*repository.Keys, error) {

	for ek, ev := range im.Enc {
		dv, ok := im.Dec[ek]
		if ok {
			return &repository.Keys{
				Encode: ev,
				Decode: dv,
			}, nil
		}
	}

	return nil, repository.ErrNotFound

}
