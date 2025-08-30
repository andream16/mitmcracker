package cli

import (
	"regexp"
)

var (
	keyLengths = map[uint]struct{}{
		24: {},
		28: {},
		32: {},
	}

	storages = map[string]struct{}{
		"memory": {},
		"disk":   {},
	}

	// hexPattern matches exactly 16 hex digits (64-bit blocks as per OCaml requirement)
	hexPattern = regexp.MustCompile(`^[0-9A-Fa-f]{16}$`)
)

// Config represents the input that can be passed to the command.
type Config struct {
	// EncText is the encrypted word.
	EncText string
	// PlainText is the plain word.
	PlainText string
	// KeyLength is the length of the key.
	// Can only be:
	//  - 24-bit
	//  - 28-bit
	//  - 32-bit
	KeyLength uint
	Storage   *Storage
}

// Storage contains storage information.
type Storage struct {
	// Type is the type of storage to be used.
	// Can only be:
	// - memory
	// - disk
	// If no option is provided, memory will be used by default.
	Type string
	// Address is the address of redis. It's optional.
	Address string
	// Password is the password of the redis database. It's optional.
	Password string
	// DB is the redis DB. It's optional.
	DB int
}

func (c *Config) Validate() error {
	if _, ok := keyLengths[c.KeyLength]; !ok {
		return InvalidInputFlagError{
			field:  "key",
			reason: "it must be either '24', '28' or '32' bits",
		}
	}

	if c.EncText == "" {
		return InvalidInputFlagError{
			field:  "encoded",
			reason: "it must not be empty",
		}
	}

	if !hexPattern.MatchString(c.EncText) {
		return InvalidInputFlagError{
			field:  "encoded",
			reason: "it must be exactly 16 hexadecimal digits (0-9, A-F)",
		}
	}

	if c.PlainText == "" {
		return InvalidInputFlagError{
			field:  "plain",
			reason: "it must not be empty",
		}
	}

	if !hexPattern.MatchString(c.PlainText) {
		return InvalidInputFlagError{
			field:  "plain",
			reason: "it must be exactly 16 hexadecimal digits (0-9, A-F)",
		}
	}

	if c.Storage != nil && c.Storage.Type != "" {
		if _, ok := storages[c.Storage.Type]; !ok {
			return InvalidInputFlagError{
				field:  "storage",
				reason: "it must be either 'memory' or 'disk'",
			}
		}
	}

	return nil
}
