package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

// Write writes config.
func (c *Config) Write() error {
	var err error

	c.mu.Lock()
	defer c.mu.Unlock()

	path, err := configPath()
	if err != nil {
		return err
	}

	log.Debug().Msgf("config: write configuration file %q", path)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { //nolint:gomnd
		return err
	}

	configFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o640) //nolint:gomnd
	if err != nil {
		return err
	}

	defer configFile.Close()

	enc := toml.NewEncoder(configFile)

	return enc.Encode(c)
}
