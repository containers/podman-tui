package tconfig

import (
	"github.com/containers/podman-tui/config/utils"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
)

// SetDefaultConnection sets default connection.
func (c *Config) SetDefaultConnection(name string) error {
	log.Debug().Msgf("config: set %s as default connection", name)

	if err := c.setDef(name); err != nil {
		return err
	}

	if err := c.write(); err != nil {
		return err
	}

	return c.reload()
}

func (c *Config) setDef(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for connName := range c.Connection.Connections {
		if connName == name {
			dest := c.Connection.Connections[connName]
			dest.Default = true
			c.Connection.Connections[connName] = dest

			return nil
		}
	}

	return utils.ErrConnectionNotFound
}

func (c *Config) GetDefaultConnection() (registry.Connection, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for connName, conn := range c.Connection.Connections {
		if conn.Default {
			return registry.Connection{
				Name:     connName,
				Identity: conn.Identity,
				URI:      conn.URI,
			}, nil
		}
	}

	return registry.Connection{}, utils.ErrDefaultConnectionNotFound
}
