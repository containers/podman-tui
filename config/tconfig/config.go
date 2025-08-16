package tconfig

import (
	"errors"
	"os"
	"sort"
	"sync"

	"github.com/containers/podman-tui/config/utils"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
)

const (
	// _configPath is the path to the podman-tui/podman-tui.json
	// inside a given config directory.
	_configPath = "podman-tui/podman-tui.json"
	// UserAppConfig holds the user podman-tui config path.
	UserAppConfig = ".config/" + _configPath
)

var ErrDuplicatedConnectionName = errors.New("duplicated connection name")

// Config contains configuration options for container tools.
type Config struct {
	mu         sync.Mutex
	Connection RemoteConnections
}

type RemoteConnections struct {
	Connections map[string]RemoteConnection `json:"connections"`
}

type RemoteConnection struct {
	// URI, required. Example: ssh://root@example.com:22/run/podman/podman.sock
	URI string `json:"uri"`

	// Identity file with ssh key, optional
	Identity string `json:"identity,omitempty"`

	// Default if its default connection, optional
	Default bool `json:"default,omitempty"`
}

// NewConfig returns new config.
func NewConfig() (*Config, error) {
	path, err := utils.ConfigPath()
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("config: loading from %q", path)

	newConfig := &Config{}
	newConfig.Connection.Connections = make(map[string]RemoteConnection)

	_, err = os.Stat(path)
	if err == nil {
		err := newConfig.readConfigFromFile(path)
		if err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	newConfig.addLocalHostIfEmptyConfig()

	return newConfig, nil
}

// RemoteConnections returns list of available connections.
func (c *Config) RemoteConnections() []registry.Connection {
	rconn := make([]registry.Connection, 0)

	c.mu.Lock()
	defer c.mu.Unlock()

	for name, conn := range c.Connection.Connections {
		rconn = append(rconn, registry.Connection{
			Name:     name,
			URI:      conn.URI,
			Identity: conn.Identity,
			Default:  conn.Default,
		})
	}

	sort.Sort(utils.ConnectionListSortedName{rconn}) //nolint:govet

	return rconn
}

func (c *Config) addLocalHostIfEmptyConfig() {
	if len(c.Connection.Connections) > 0 {
		return
	}

	c.Connection.Connections = make(map[string]RemoteConnection)
	c.Connection.Connections["localhost"] = RemoteConnection{
		URI:     utils.LocalNodeUnixSocket(),
		Default: true,
	}
}
