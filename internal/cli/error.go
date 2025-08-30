package cli

import "fmt"

// InvalidInputFlagError is used when an invalid flag is passed as cli argument.
type InvalidInputFlagError struct {
	field  string
	reason string
}

func (e InvalidInputFlagError) Error() string {
	return fmt.Sprintf("invalid input flag '%s': %s", e.field, e.reason)
}

// Field returns the field name that caused the error.
func (e InvalidInputFlagError) Field() string {
	return e.field
}

// Reason returns the reason for the error.
func (e InvalidInputFlagError) Reason() string {
	return e.reason
}
