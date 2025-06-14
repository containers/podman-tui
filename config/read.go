package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

func (c *Config) readConfigFromFile(path string) error {
	log.Debug().Msgf("config: reading configuration file %q", path)

	c.mu.Lock()
	defer c.mu.Unlock()

	cfgFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config: %w open configuration %q", err, path)
	}

	cfgData, err := io.ReadAll(cfgFile)
	if err != nil {
		return fmt.Errorf("config: %w read configuration %q", err, path)
	}

	return json.Unmarshal(cfgData, &c.Connection)
}

func (c *Config) reload() error {
	log.Debug().Msgf("config: reload configuration")

	path, err := configPath()
	if err != nil {
		return err
	}

	if err := c.readConfigFromFile(path); err != nil {
		return err
	}

	return nil
}
