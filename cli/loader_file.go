package cli

import (
	"errors"
	"os"
	"strings"

	"github.com/traefik/paerser/file"
	"github.com/traefik/paerser/flag"
	"github.com/traefik/paerser/parser"
)

// FileLoader loads a configuration from a file.
type FileLoader struct {
	ConfigFileFlag string
	filename       string
	BasePaths      []string
	Extensions     []string
}

// GetFilename returns the configuration file if any.
func (f *FileLoader) GetFilename() string {
	return f.filename
}

// Load loads the command's configuration from a file either specified with the ConfigFileFlag flag, or from default locations.
func (f *FileLoader) Load(args []string, cmd *Command) (bool, error) {
	ref, err := flag.Parse(args, cmd.Configuration)
	if err != nil {
		_ = cmd.PrintHelp(os.Stdout)
		return false, err
	}

	var configFileFlag string
	if f.ConfigFileFlag == "" {
		return false, errors.New("missing config file flag")
	}

	configFileFlag = parser.DefaultRootName + "." + f.ConfigFileFlag
	if _, ok := ref[strings.ToLower(configFileFlag)]; ok {
		configFileFlag = parser.DefaultRootName + "." + strings.ToLower(f.ConfigFileFlag)
	}

	configFile, err := f.loadConfigFiles(ref[configFileFlag], cmd.Configuration)
	if err != nil {
		return false, err
	}

	f.filename = configFile

	if configFile == "" {
		return false, nil
	}

	return true, nil
}

// loadConfigFiles tries to decode the given configuration file and all default locations for the configuration file.
// It stops as soon as decoding one of them is successful.
func (f *FileLoader) loadConfigFiles(configFile string, element interface{}) (string, error) {
	extensions := []string{"toml", "yaml", "yml"}
	if len(f.Extensions) != 0 {
		extensions = f.Extensions
	}

	if len(f.BasePaths) == 0 {
		return "", errors.New("missing base paths")
	}

	finder := Finder{
		BasePaths:  f.BasePaths,
		Extensions: extensions,
	}

	filePath, err := finder.Find(configFile)
	if err != nil {
		return "", err
	}

	if len(filePath) == 0 {
		return "", nil
	}

	if err = file.Decode(filePath, element); err != nil {
		return "", err
	}

	return filePath, nil
}
