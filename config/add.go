package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Add adds a new remote connection.
func (c *Config) Add(name string, uri string, identity string) error {
	log.Debug().Msgf("config: adding new remote connection %s %s %s", name, uri, identity)

	conn, err := validateNewConnection(name, uri, identity)
	if err != nil {
		return err
	}

	if err := c.add(name, conn); err != nil {
		return err
	}

	if err := c.Write(); err != nil {
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

func validateNewConnection(name string, dest string, identity string) (RemoteConnection, error) { //nolint:cyclop
	var (
		conn         RemoteConnection
		connIdentity string
	)

	if name == "" {
		return conn, ErrEmptyConnectionName
	}

	if dest == "" {
		return conn, ErrEmptyURIDestination
	}

	if match, err := regexp.Match("^[A-Za-z][A-Za-z0-9+.-]*://", []byte(dest)); err != nil { //nolint:mirror
		return conn, fmt.Errorf("%w invalid destition", err)
	} else if !match {
		dest = "ssh://" + dest
	}

	uri, err := url.Parse(dest)
	if err != nil {
		return conn, err
	}

	switch uri.Scheme {
	case "ssh":
		if uri.User.Username() == "" {
			if uri.User, err = getUserInfo(uri); err != nil {
				return conn, err
			}
		}

		connIdentity, err = utils.ResolveHomeDir(identity)
		if err != nil {
			return conn, err
		}

		if identity == "" {
			return conn, ErrEmptySSHIdentity
		}

		if uri.Port() == "" {
			uri.Host = net.JoinHostPort(uri.Hostname(), "22")
		}

		if uri.Path == "" || uri.Path == "/" {
			if uri.Path, err = getUDS(uri, connIdentity); err != nil {
				return conn, err
			}
		}
	case "unix":
		if identity != "" {
			return conn, fmt.Errorf("%w identity", ErrInvalidUnixSchemaOption)
		}

		info, err := os.Stat(uri.Path)

		switch {
		case errors.Is(err, os.ErrNotExist):
			log.Warn().Msgf("config: %q does not exists", uri.Path)
		case errors.Is(err, os.ErrPermission):
			log.Warn().Msgf("config: You do not have permission to read %q", uri.Path)
		case err != nil:
			return conn, err
		case info.Mode()&os.ModeSocket == 0:
			return conn, fmt.Errorf("%w %q", ErrFileNotUnixSocket, uri.Path)
		}
	case "tcp":
		if identity != "" {
			return conn, fmt.Errorf("%w identity", ErrInvalidTCPSchemaOption)
		}
	default:
		return conn, fmt.Errorf("%w %q", ErrInvalidURISchemaName, uri.Scheme)
	}

	conn.Identity = connIdentity
	conn.URI = uri.String()
	conn.Default = false

	return conn, nil
}
