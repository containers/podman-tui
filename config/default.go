package config

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
)

// SetDefaultService sets default service name.
func (c *Config) SetDefaultService(name string) error {
	log.Debug().Msgf("config: set %s as default service", name)

	if err := c.setDef(name); err != nil {
		return err
	}

	if err := c.Write(); err != nil {
		return err
	}

	return c.reload()
}

func (c *Config) setDef(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.Services {
		dest := c.Services[key]
		dest.Default = false

		if key == name {
			dest.Default = true
		}

		c.Services[key] = dest
	}

	return nil
}

func (c *Config) getDefault() registry.Connection {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, service := range c.Services {
		if service.Default {
			return registry.Connection{
				Name:     name,
				Identity: service.Identity,
				URI:      service.URI,
			}
		}
	}

	return registry.Connection{}
}
