package decrypter

import (
	"bytes"
	"os/exec"
)

// Decrypter decrypts a given cipherText using a key.
type Decrypter func(key, cipherText string) (string, error)

// DefaultDecrypt decrypts a given cipherText using a key using the supplied decrypter binary.
func DefaultDecrypt(key, plainText string) (string, error) {
	var (
		command = exec.Command("./resources/decrypt", "-s", key, plainText)
		out bytes.Buffer
	)

	command.Stdout = &out

	if err := command.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}
