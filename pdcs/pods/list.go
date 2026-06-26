package pods

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/pods"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

// List returns list of pods.
func List() ([]*entities.ListPodsReport, error) {
	log.Debug().Msg("pdcs: podman pod ls")

	conn, err := registry.GetConnection()
	if err != nil {
		return nil, err
	}

	response, err := pods.List(conn, new(pods.ListOptions))
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("pdcs: %v", response)

	return response, nil
}
