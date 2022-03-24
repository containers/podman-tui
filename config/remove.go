package config

import "github.com/rs/zerolog/log"

// Remove removes a service from config
func (c *Config) Remove(name string) error {
	log.Debug().Msgf("config: remove service %q", name)
	c.remove(name)
	if err := c.Write(); err != nil {
		return err
	}
	return c.reload()
}

func (c *Config) remove(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for serviceName := range c.Services {
		if serviceName == name {
			delete(c.Services, name)
		}
	}

}
