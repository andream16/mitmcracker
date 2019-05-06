package memory

// InMemo represents an in-memory map.
type InMemo struct {
	Dec map[string]string
	Enc map[string]string
}

// InsertDec adds a new entry in dec map.
func (im *InMemo) InsertDec(key, cipherText string) {

	if len(im.Dec) == 0 {
		im.Dec = map[string]string{}
	}

	im.Dec[cipherText] = key

}

// InsertEnc adds a new entry in enc map.
func (im *InMemo) InsertEnc(key, cipherText string) {

	if len(im.Enc) == 0 {
		im.Enc = map[string]string{}
	}

	im.Enc[cipherText] = key

}

// FindKey returns the key used to encode & decode the same message.
func (im *InMemo) FindKey() string {

	for ek := range im.Enc {
		_, ok := im.Dec[ek]
		if ok {
			return ek
		}
	}

	return ""

}
