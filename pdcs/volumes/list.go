package volumes

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/volumes"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

// List returns list of volumes.
func List() ([]*entities.VolumeListReport, error) {
	log.Debug().Msg("pdcs: podman volume ls")

	conn, err := registry.GetConnection()
	if err != nil {
		return nil, err
	}

	response, err := volumes.List(conn, new(volumes.ListOptions))
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("pdcs: %v", response)

	return response, nil
}
