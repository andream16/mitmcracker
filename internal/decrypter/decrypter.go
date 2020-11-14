package decrypter

import "os/exec"

// Decrypter decrypts a given cipherText using a key.
type Decrypter func(key, cipherText string) (string, error)

// DefaultDecrypt decrypts a given cipherText using a key using the supplied decrypter binary.
func DefaultDecrypt(key, plainText string) (string, error) {
	out, err := exec.Command("./resources/decrypt", "-s", key, plainText).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
