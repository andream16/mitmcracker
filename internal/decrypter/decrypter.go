package decrypter

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Decrypter decrypts a given cipherText using a key.
type Decrypter func(key, cipherText string) (string, error)

// DefaultDecrypt decrypts a given cipherText using a key using the supplied decrypter binary.
func DefaultDecrypt(key, cipherText string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("decryption key cannot be empty")
	}
	if cipherText == "" {
		return "", fmt.Errorf("ciphertext cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		command = exec.CommandContext(ctx, "./resources/decrypt", "-s", key, cipherText)
		out     bytes.Buffer
		errOut  bytes.Buffer
	)

	command.Stdout = &out
	command.Stderr = &errOut

	if err := command.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("decryption timeout after 10 seconds")
		}
		if errOut.Len() > 0 {
			return "", fmt.Errorf("decryption failed: %s", strings.TrimSpace(errOut.String()))
		}
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		return "", fmt.Errorf("decryption produced empty result")
	}

	return result, nil
}
