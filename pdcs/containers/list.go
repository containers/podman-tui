package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

// List returns list of containers information.
func List() ([]entities.ListContainer, error) {
	log.Debug().Msg("pdcs: podman container ls")

	conn, err := registry.GetConnection()
	if err != nil {
		return nil, err
	}

	response, err := containers.List(conn, new(containers.ListOptions).WithAll(true))
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("pdcs: %v", response)

	return response, nil
}
