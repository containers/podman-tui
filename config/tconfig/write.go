package tconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/containers/podman-tui/config/utils"
	"github.com/rs/zerolog/log"
)

// Write writes config.
func (c *Config) write() error {
	var err error

	c.mu.Lock()
	defer c.mu.Unlock()

	path, err := utils.ConfigPath()
	if err != nil {
		return err
	}

	log.Debug().Msgf("config: write configuration file %q", path)

	err = os.MkdirAll(filepath.Dir(path), 0o750) //nolint:mnd
	if err != nil {
		return err
	}

	configFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o600) //nolint:mnd,gosec
	if err != nil {
		return err
	}

	defer func() {
		err := configFile.Close()
		if err != nil {
			log.Error().Msgf("failed to close config file after write: %s", err.Error())
		}
	}()

	jsonData, err := json.Marshal(c.Connection)
	if err != nil {
		return fmt.Errorf("config: configuration json marshal %w", err)
	}

	_, err = configFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("config: %w write configuration %q", err, path)
	}

	return nil
}
