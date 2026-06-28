package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/containers"
)

// Stop starts a container.
func Stop(id string) error {
	log.Debug().Msgf("pdcs: podman container stop %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return containers.Stop(conn, id, new(containers.StopOptions))
}
