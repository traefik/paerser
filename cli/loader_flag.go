package cli

import (
	"fmt"

	"github.com/traefik/paerser/flag"
)

// FlagLoader loads configuration from flags.
type FlagLoader struct{}

// Load loads the command's configuration from flag arguments.
func (*FlagLoader) Load(args []string, cmd *Command) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	if err := flag.Decode(args, cmd.Configuration); err != nil {
		return false, fmt.Errorf("failed to decode configuration from flags: %w", err)
	}

	return true, nil
}
