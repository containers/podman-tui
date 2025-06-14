package config

import (
	"github.com/containers/podman-tui/config/pconfig"
	"github.com/containers/podman-tui/config/tconfig"
	"github.com/containers/podman-tui/pdcs/registry"
)

type Config interface {
	RemoteConnections() []registry.Connection
	SetDefaultConnection(name string) error
	GetDefaultConnection() (registry.Connection, error)
	Add(name string, uri string, identity string) error
	Remove(name string) error
}

func NewConfig() (Config, error) { //nolint:ireturn
	var cfg Config

	pconfig, err := pconfig.NewConfig()
	if err != nil {
		return nil, err
	}

	premoteConns := pconfig.RemoteConnections()
	if len(premoteConns) > 0 {
		cfg = pconfig
	} else {
		tconfig, err := tconfig.NewConfig()
		if err != nil {
			return nil, err
		}

		cfg = tconfig
	}

	defaultConn, err := cfg.GetDefaultConnection()
	if err != nil {
		return nil, err
	}

	registry.SetConnection(defaultConn)

	return cfg, nil
}
