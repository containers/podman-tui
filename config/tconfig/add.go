package tconfig

import (
	"github.com/containers/podman-tui/config/utils"
	"github.com/rs/zerolog/log"
)

// Add adds a new remote connection.
func (c *Config) Add(name string, uri string, identity string) error {
	log.Debug().Msgf("config: adding new remote connection %s %s %s", name, uri, identity)

	connURI, err := utils.ValidateNewConnection(name, uri, identity)
	if err != nil {
		return err
	}

	conn := RemoteConnection{
		URI:      connURI,
		Identity: identity,
		Default:  false,
	}

	if err := c.add(name, conn); err != nil {
		return err
	}

	if err := c.write(); err != nil {
		return err
	}

	return c.reload()
}

func (c *Config) add(name string, conn RemoteConnection) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for connName := range c.Connection.Connections {
		if connName == name {
			return ErrDuplicatedConnectionName
		}
	}

	c.Connection.Connections[name] = conn

	return nil
}
