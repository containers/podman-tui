package config

import (
	"github.com/containers/podman-tui/pdcs/registry"
)

type Config interface {
	RemoteConnections() []registry.Connection
	SetDefaultConnection(name string) error
	GetDefaultConnection() registry.Connection
	Add(name string, uri string, identity string) error
	Remove(name string) error
}

func NewConfig() (Config, error) { //nolint:ireturn
	var cfg Config

	// load podman remote connections config
	pconfig, err := NewPodmanRemoteConfig()
	if err != nil {
		return nil, err
	}

	cfg = pconfig

	defaultConn := cfg.GetDefaultConnection()
	if defaultConn.URI != "" && defaultConn.Name != "" {
		registry.SetConnection(defaultConn)
	}

	return cfg, nil
}
