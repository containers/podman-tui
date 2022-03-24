package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Add adds new service connection
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
			return fmt.Errorf("duplicated service name")
		}
	}

	c.Services[name] = newService
	return nil
}

// most of codes are from:
// https://github.com/containers/podman/blob/main/cmd/podman/system/connection/add.go
func validateNewService(name string, dest string, identity string) (Service, error) {
	var (
		service         Service
		serviceIdentity string
	)
	if name == "" {
		return service, fmt.Errorf("empty service name %q", name)
	}
	if dest == "" {
		return service, fmt.Errorf("empty URI %q", dest)
	}
	if match, err := regexp.Match("^[A-Za-z][A-Za-z0-9+.-]*://", []byte(dest)); err != nil {
		return service, fmt.Errorf("%v invalid destition", err)
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
		serviceIdentity, err = resolveHomeDir(identity)
		if err != nil {
			return service, err
		}
		if identity == "" {
			return service, fmt.Errorf("%q empty identity field for SSH connection", identity)
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
			return service, fmt.Errorf("identity option not supported for unix scheme")
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
			return service, fmt.Errorf("%q exists and is not a unix domain socket", uri.Path)
		}

	case "tcp":
		if identity != "" {
			return service, fmt.Errorf("identity option not supported for tcp scheme")
		}
	default:
		return service, fmt.Errorf("%q invalid schema name", uri.Scheme)
	}

	service.Identity = serviceIdentity
	service.URI = uri.String()
	service.Default = false

	return service, nil
}
