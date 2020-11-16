package encrypter

import (
	"bytes"
	"os/exec"
)

// Encrypter encrypts a given plaintext using a key.
type Encrypter func(key, plainText string) (string, error)

// DefaultEncrypt encrypts a given plaintext using a key using the supplied encrypt binary.
func DefaultEncrypt(key, plainText string) (string, error) {
	var (
		command = exec.Command("./resources/encrypt", "-s", key, plainText)
		out bytes.Buffer
	)

	command.Stdout = &out

	if err := command.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}
