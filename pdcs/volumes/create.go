package volumes

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rs/zerolog/log"
)

// CreateOptions implements new volume create options.
type CreateOptions struct {
	Name          string
	Labels        map[string]string
	Driver        string
	DriverOptions map[string]string
}

// Create creates a new volume.
func Create(opts CreateOptions) (string, error) {
	log.Debug().Msgf("pdcs: podman volume create %v", opts)

	var report string

	volCreateOptions := entities.VolumeCreateOptions{
		Name:    opts.Name,
		Label:   opts.Labels,
		Driver:  opts.Driver,
		Options: opts.DriverOptions,
	}

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := volumes.Create(conn, volCreateOptions, &volumes.CreateOptions{})
	if err != nil {
		return report, err
	}

	report, err = utils.GetJSONOutput(response)
	if err != nil {
		return report, err
	}

	return report, nil
}
