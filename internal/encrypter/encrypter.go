package encrypter

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Encrypter encrypts a given plaintext using a key.
type Encrypter func(key, plainText string) (string, error)

// DefaultEncrypt encrypts a given plaintext using a key using the supplied encrypt binary.
func DefaultEncrypt(key, plainText string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("encryption key cannot be empty")
	}
	if plainText == "" {
		return "", fmt.Errorf("plaintext cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		command = exec.CommandContext(ctx, "./resources/encrypt", "-s", key, plainText)
		out     bytes.Buffer
		errOut  bytes.Buffer
	)

	command.Stdout = &out
	command.Stderr = &errOut

	if err := command.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("encryption timeout after 10 seconds")
		}
		if errOut.Len() > 0 {
			return "", fmt.Errorf("encryption failed: %s", strings.TrimSpace(errOut.String()))
		}
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		return "", fmt.Errorf("encryption produced empty result")
	}

	return result, nil
}
