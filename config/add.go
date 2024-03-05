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

// Add adds new service connection.
func (c *Config) Add(name string, uri string, identity string) error {
	log.Debug().Msgf("config: adding new service %s %s %s", name, uri, identity)

	newService, err := validateNewService(name, uri, identity)
	if err != nil {
		return err
	}

	if err := c.add(name, newService); err != nil {
		return err
	}

	if err := c.Write(); err != nil {
		return err
	}

	return c.reload()
}

func (c *Config) add(name string, newService Service) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for serviceName := range c.Services {
		if serviceName == name {
			return ErrDuplicatedServiceName
		}
	}

	c.Services[name] = newService

	return nil
}

// most of codes are from:
// https://github.com/containers/podman/blob/main/cmd/podman/system/connection/add.go.
func validateNewService(name string, dest string, identity string) (Service, error) { //nolint:cyclop
	var (
		service         Service
		serviceIdentity string
	)

	if name == "" {
		return service, ErrEmptyServiceName
	}

	if dest == "" {
		return service, ErrEmptyURIDestination
	}

	if match, err := regexp.Match("^[A-Za-z][A-Za-z0-9+.-]*://", []byte(dest)); err != nil { //nolint:mirror
		return service, fmt.Errorf("%w invalid destition", err)
	} else if !match {
		dest = "ssh://" + dest
	}

	uri, err := url.Parse(dest)
	if err != nil {
		return service, err
	}

	switch uri.Scheme {
	case "ssh":
		if uri.User.Username() == "" {
			if uri.User, err = getUserInfo(uri); err != nil {
				return service, err
			}
		}

		serviceIdentity, err = utils.ResolveHomeDir(identity)
		if err != nil {
			return service, err
		}

		if identity == "" {
			return service, ErrEmptySSHIdentity
		}

		if uri.Port() == "" {
			uri.Host = net.JoinHostPort(uri.Hostname(), "22")
		}

		if uri.Path == "" || uri.Path == "/" {
			if uri.Path, err = getUDS(uri, serviceIdentity); err != nil {
				return service, err
			}
		}
	case "unix":
		if identity != "" {
			return service, fmt.Errorf("%w identity", ErrInvalidUnixSchemaOption)
		}

		info, err := os.Stat(uri.Path)

		switch {
		case errors.Is(err, os.ErrNotExist):
			log.Warn().Msgf("config: %q does not exists", uri.Path)
		case errors.Is(err, os.ErrPermission):
			log.Warn().Msgf("config: You do not have permission to read %q", uri.Path)
		case err != nil:
			return service, err
		case info.Mode()&os.ModeSocket == 0:
			return service, fmt.Errorf("%w %q", ErrFileNotUnixSocket, uri.Path)
		}
	case "tcp":
		if identity != "" {
			return service, fmt.Errorf("%w identity", ErrInvalidTCPSchemaOption)
		}
	default:
		return service, fmt.Errorf("%w %q", ErrInvalidURISchemaName, uri.Scheme)
	}

	service.Identity = serviceIdentity
	service.URI = uri.String()
	service.Default = false

	return service, nil
}
