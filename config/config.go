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
	// _configPath is the path to the podman-tui/podman-tui.conf
	// inside a given config directory.
	_configPath = "podman-tui/podman-tui.conf"
	// UserAppConfig holds the user podman-tui config path.
	UserAppConfig = ".config/" + _configPath
)

var (
	ErrRemotePodmanUDSReport   = errors.New("remote podman failed to report its UDS socket")
	ErrInvalidURISchemaName    = errors.New("invalid schema name")
	ErrInvalidTCPSchemaOption  = errors.New("invalid option for tcp")
	ErrInvalidUnixSchemaOption = errors.New("invalid option for unix")
	ErrFileNotUnixSocket       = errors.New("not a unix domain socket")
	ErrEmptySSHIdentity        = errors.New("empty identity field for SSH connection")
	ErrEmptyURIDestination     = errors.New("empty URI destination")
	ErrEmptyServiceName        = errors.New("empty service name")
	ErrDuplicatedServiceName   = errors.New("duplicated service name")
)

// Config contains configuration options for container tools.
type Config struct {
	mu sync.Mutex
	// Services specify the service destination connections
	Services map[string]Service `toml:"services,omitempty"`
}

// Service represents remote service destination.
type Service struct {
	// URI, required. Example: ssh://root@example.com:22/run/podman/podman.sock
	URI string `toml:"uri"`

	// Identity file with ssh key, optional
	Identity string `toml:"identity,omitempty"`

	// Default if its default service, optional
	Default bool `toml:"default,omitempty"`
}

// NewConfig returns new config.
func NewConfig() (*Config, error) {
	log.Debug().Msgf("config: new")

	path, err := configPath()
	if err != nil {
		return nil, err
	}

	newConfig := &Config{}
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
	if len(c.Services) > 0 {
		return
	}

	c.Services = make(map[string]Service)
	c.Services["localhost"] = Service{
		URI:     localNodeUnixSocket(),
		Default: true,
	}
}

// ServicesConnections returns list of available connections.
func (c *Config) ServicesConnections() []registry.Connection {
	conn := make([]registry.Connection, 0)

	c.mu.Lock()
	defer c.mu.Unlock()

	for name, service := range c.Services {
		conn = append(conn, registry.Connection{
			Name:     name,
			URI:      service.URI,
			Identity: service.Identity,
			Default:  service.Default,
		})
	}

	sort.Sort(connectionListSortedName{conn})

	return conn
}

type connSort []registry.Connection

func (a connSort) Len() int      { return len(a) }
func (a connSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type connectionListSortedName struct{ connSort }

func (a connectionListSortedName) Less(i, j int) bool {
	return a.connSort[i].Name < a.connSort[j].Name
}
