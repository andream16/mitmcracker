package cli

import "fmt"

// InvalidInputFlagError is used when an invalid flag is passed as cli argument.
type InvalidInputFlagError struct {
	field  string
	reason string
}

func (i InvalidInputFlagError) Error() string {
	return fmt.Sprintf("invalid input flag %s: %s", i.field, i.reason)
}
