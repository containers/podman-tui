package tconfig

import "github.com/rs/zerolog/log"

// Remove removes a connection from config.
func (c *Config) Remove(name string) error {
	log.Debug().Msgf("config: remove remote connection %q", name)

	c.remove(name)

	if err := c.write(); err != nil {
		return err
	}

	return c.reload()
}

func (c *Config) remove(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for connName := range c.Connection.Connections {
		if connName != name {
			delete(c.Connection.Connections, name)
		}
	}
}
