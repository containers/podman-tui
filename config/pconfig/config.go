package pconfig

import (
	"slices"

	cconfig "github.com/containers/common/pkg/config"
	"github.com/containers/podman-tui/config/utils"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rs/zerolog/log"
)

type Config struct {
	podmanOptions *entities.PodmanConfig
}

func NewConfig() (*Config, error) {
	log.Debug().Msg("config: loading podman remote connections")

	newConfig := &Config{}

	defaultConfig, err := cconfig.New(&cconfig.Options{
		SetDefault: true,
		Modules:    nil,
	})
	if err != nil {
		return nil, err
	}

	podmanOptions := entities.PodmanConfig{ContainersConf: &cconfig.Config{}, ContainersConfDefaultsRO: defaultConfig}
	newConfig.podmanOptions = &podmanOptions

	return newConfig, nil
}

func (c *Config) RemoteConnections() []registry.Connection {
	rconn := make([]registry.Connection, 0)

	conns, err := c.podmanOptions.ContainersConfDefaultsRO.GetAllConnections()
	if err != nil {
		log.Err(err).Msgf("config: podman remote connection")

		return nil
	}

	log.Debug().Msgf("connections: %v", conns)

	for _, conn := range conns {
		rconn = append(rconn, registry.Connection{
			Name:     conn.Name,
			URI:      conn.URI,
			Identity: conn.Identity,
			Default:  conn.Default,
		})
	}

	return rconn
}

func (c *Config) Remove(name string) error {
	return cconfig.EditConnectionConfig(func(cfg *cconfig.ConnectionsFile) error {
		delete(cfg.Connection.Connections, name)

		if cfg.Connection.Default == name {
			cfg.Connection.Default = ""
		}

		// If there are existing farm, remove the deleted connection that might be part of a farm
		for k, v := range cfg.Farm.List {
			index := slices.Index(v, name)
			if index > -1 {
				cfg.Farm.List[k] = append(v[:index], v[index+1:]...)
			}
		}

		return nil
	})
}

func (c *Config) Add(name string, uri string, identity string) error {
	connURI, err := utils.ValidateNewConnection(name, uri, identity)
	if err != nil {
		return err
	}

	dst := cconfig.Destination{
		URI:      connURI,
		Identity: identity,
	}

	return cconfig.EditConnectionConfig(func(cfg *cconfig.ConnectionsFile) error {
		if cfg.Connection.Connections == nil {
			cfg.Connection.Connections = map[string]cconfig.Destination{
				name: dst,
			}
			cfg.Connection.Default = name
		} else {
			cfg.Connection.Connections[name] = dst
		}

		return nil
	})
}

func (c *Config) SetDefaultConnection(name string) error {
	return cconfig.EditConnectionConfig(func(cfg *cconfig.ConnectionsFile) error {
		if _, found := cfg.Connection.Connections[name]; !found {
			return utils.ErrConnectionNotFound
		}

		cfg.Connection.Default = name

		return nil
	})
}

func (c *Config) GetDefaultConnection() registry.Connection {
	for _, conn := range c.RemoteConnections() {
		if conn.Default {
			return registry.Connection{
				Name:     conn.Name,
				Identity: conn.Identity,
				URI:      conn.URI,
			}
		}
	}

	return registry.Connection{}
}
