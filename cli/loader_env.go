package cli

import (
	"fmt"
	"os"

	"github.com/traefik/paerser/env"
)

// EnvLoader loads a configuration from all the environment variables prefixed with Prefix (default: "TRAEFIK_").
type EnvLoader struct {
	Prefix string
}

// Load loads the command's configuration from the environment variables.
func (e *EnvLoader) Load(_ []string, cmd *Command) (bool, error) {
	prefix := env.DefaultNamePrefix
	if e.Prefix != "" {
		prefix = e.Prefix
	}

	vars := env.FindPrefixedEnvVars(os.Environ(), prefix, cmd.Configuration)
	if len(vars) == 0 {
		return false, nil
	}

	if err := env.Decode(vars, prefix, cmd.Configuration); err != nil {
		return false, fmt.Errorf("failed to decode configuration from environment variables: %w ", err)
	}

	return true, nil
}
