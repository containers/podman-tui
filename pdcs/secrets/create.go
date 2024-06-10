package secrets

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/secrets"
	"github.com/rs/zerolog/log"
)

var (
	errSecretInvalidLabelFormat        = errors.New("invalid label format")
	errSecretInvalidDriverOptionFormat = errors.New("invalid driver option format")
)

// SecretCreateOptions secret create options.
type SecretCreateOptions struct {
	Name          string
	Replace       bool
	File          string
	Text          string
	Labels        []string
	Driver        string
	DriverOptions []string
}

// Create creates a new secret.
func Create(opts *SecretCreateOptions) error { //nolint:cyclop
	log.Debug().Msgf("pdcs: podman secret create %v", opts)

	var reader io.Reader

	createOpts := new(secrets.CreateOptions)
	createOpts = createOpts.WithReplace(opts.Replace)
	createOpts = createOpts.WithName(opts.Name)
	createOpts = createOpts.WithDriver(opts.Driver)
	labels := make(map[string]string)

	for _, label := range opts.Labels {
		if label == "" {
			continue
		}

		key, value, _ := strings.Cut(label, "=")
		if key == "" {
			return errSecretInvalidLabelFormat
		}

		labels[key] = value
	}

	createOpts.WithLabels(labels)

	driverOptions := make(map[string]string)

	for _, driverOpt := range opts.DriverOptions {
		if driverOpt == "" {
			continue
		}

		key, value, _ := strings.Cut(driverOpt, "=")
		if key == "" {
			return errSecretInvalidDriverOptionFormat
		}

		driverOptions[key] = value
	}

	createOpts.WithDriverOpts(driverOptions)

	if opts.File != "" {
		file, err := os.Open(opts.File)
		if err != nil {
			return err
		}

		defer file.Close()
		reader = file
	}

	if opts.Text != "" {
		reader = strings.NewReader(opts.Text)
	}

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	if _, err := secrets.Create(conn, reader, createOpts); err != nil {
		return err
	}

	return nil
}
