package config

import (
	"errors"
	"os"
	"sort"
	"sync"

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

var (
	ErrRemotePodmanUDSReport    = errors.New("remote podman failed to report its UDS socket")
	ErrInvalidURISchemaName     = errors.New("invalid schema name")
	ErrInvalidTCPSchemaOption   = errors.New("invalid option for tcp")
	ErrInvalidUnixSchemaOption  = errors.New("invalid option for unix")
	ErrFileNotUnixSocket        = errors.New("not a unix domain socket")
	ErrEmptySSHIdentity         = errors.New("empty identity field for SSH connection")
	ErrEmptyURIDestination      = errors.New("empty URI destination")
	ErrEmptyConnectionName      = errors.New("empty connection name")
	ErrDuplicatedConnectionName = errors.New("duplicated connection name")
)

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
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("loading config from %q", path)

	newConfig := &Config{}
	newConfig.Connection.Connections = make(map[string]RemoteConnection)

	if _, err := os.Stat(path); err == nil {
		if err := newConfig.readConfigFromFile(path); err != nil {
			return nil, err
		}
	} else {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	newConfig.addLocalHostIfEmptyConfig()

	defaultConn := newConfig.getDefault()
	if defaultConn.URI != "" {
		registry.SetConnection(defaultConn)
	}

	return newConfig, nil
}

func (c *Config) addLocalHostIfEmptyConfig() {
	if len(c.Connection.Connections) > 0 {
		return
	}

	c.Connection.Connections = make(map[string]RemoteConnection)
	c.Connection.Connections["localhost"] = RemoteConnection{
		URI:     localNodeUnixSocket(),
		Default: true,
	}
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

	sort.Sort(connectionListSortedName{rconn})

	return rconn
}

type connSort []registry.Connection

func (a connSort) Len() int      { return len(a) }
func (a connSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type connectionListSortedName struct{ connSort }

func (a connectionListSortedName) Less(i, j int) bool {
	return a.connSort[i].Name < a.connSort[j].Name
}
