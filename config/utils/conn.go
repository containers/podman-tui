package utils

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidURISchemaName      = errors.New("invalid schema name")
	ErrInvalidTCPSchemaOption    = errors.New("invalid option for tcp")
	ErrInvalidUnixSchemaOption   = errors.New("invalid option for unix")
	ErrFileNotUnixSocket         = errors.New("not a unix domain socket")
	ErrEmptySSHIdentity          = errors.New("empty identity field for SSH connection")
	ErrEmptyURIDestination       = errors.New("empty URI destination")
	ErrEmptyConnectionName       = errors.New("empty connection name")
	ErrConnectionNotFound        = errors.New("connection not found")
	ErrDefaultConnectionNotFound = errors.New("default connection not found")
)

func ValidateNewConnection(name string, dest string, identity string) (string, error) { //nolint:cyclop
	var connIdentity string

	if name == "" {
		return "", ErrEmptyConnectionName
	}

	if dest == "" {
		return "", ErrEmptyURIDestination
	}

	if match, err := regexp.Match("^[A-Za-z][A-Za-z0-9+.-]*://", []byte(dest)); err != nil { //nolint:mirror
		return "", fmt.Errorf("%w invalid destition", err)
	} else if !match {
		dest = "ssh://" + dest
	}

	uri, err := url.Parse(dest)
	if err != nil {
		return "", err
	}

	switch uri.Scheme {
	case "ssh":
		if uri.User.Username() == "" {
			if uri.User, err = getUserInfo(uri); err != nil {
				return "", err
			}
		}

		connIdentity, err = utils.ResolveHomeDir(identity)
		if err != nil {
			return "", err
		}

		if identity == "" {
			return "", ErrEmptySSHIdentity
		}

		if uri.Port() == "" {
			uri.Host = net.JoinHostPort(uri.Hostname(), "22")
		}

		if uri.Path == "" || uri.Path == "/" {
			if uri.Path, err = getUDS(uri, connIdentity); err != nil {
				return "", err
			}
		}
	case "unix":
		if identity != "" {
			return "", fmt.Errorf("%w identity", ErrInvalidUnixSchemaOption)
		}

		info, err := os.Stat(uri.Path)

		switch {
		case errors.Is(err, os.ErrNotExist):
			log.Warn().Msgf("config: %q does not exists", uri.Path)
		case errors.Is(err, os.ErrPermission):
			log.Warn().Msgf("config: You do not have permission to read %q", uri.Path)
		case err != nil:
			return "", err
		case info.Mode()&os.ModeSocket == 0:
			return "", fmt.Errorf("%w %q", ErrFileNotUnixSocket, uri.Path)
		}
	case "tcp":
		if identity != "" {
			return "", fmt.Errorf("%w identity", ErrInvalidTCPSchemaOption)
		}
	default:
		return "", fmt.Errorf("%w %q", ErrInvalidURISchemaName, uri.Scheme)
	}

	return uri.String(), nil
}
