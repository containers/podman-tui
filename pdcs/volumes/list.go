package volumes

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rs/zerolog/log"
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
