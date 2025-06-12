package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { //nolint:mnd
		return err
	}

	configFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o640) //nolint:mnd
	if err != nil {
		return err
	}

	defer configFile.Close()

	jsonData, err := json.Marshal(c.Connection)
	if err != nil {
		return fmt.Errorf("config: configuration json marshal %w", err)
	}

	if _, err := configFile.Write(jsonData); err != nil {
		return fmt.Errorf("config: %w write configuration %q", err, path)
	}

	return nil
}
