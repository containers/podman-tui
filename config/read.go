package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

func (c *Config) readConfigFromFile(path string) error {
	log.Debug().Msgf("config: reading configuration file %q", path)

	c.mu.Lock()
	defer c.mu.Unlock()

	rawConfig, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config: %w read configuration %q", err, path)
	}

	config := os.ExpandEnv(string(rawConfig))

	meta, err := toml.Decode(config, c)
	if err != nil {
		return fmt.Errorf("config: %w decode configuration %q", err, path)
	}

	keys := meta.Undecoded()
	if len(keys) > 0 {
		log.Debug().Msgf("config: failed to decode the keys %q from %q.", keys, path)
	}

	return nil
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
