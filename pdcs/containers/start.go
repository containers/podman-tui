package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/containers"
)

// Start starts a container.
func Start(id string) error {
	log.Debug().Msgf("pdcs: podman container start %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return containers.Start(conn, id, new(containers.StartOptions))
}
