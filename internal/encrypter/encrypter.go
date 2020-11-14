package encrypter

import "os/exec"

// Encrypter encrypts a given plaintext using a key.
type Encrypter func(key, plainText string) (string, error)

// DefaultEncrypt encrypts a given plaintext using a key using the supplied encrypt binary.
func DefaultEncrypt(key, plainText string) (string, error) {
	out, err := exec.Command("./resources/encrypt", "-s", key, plainText).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
