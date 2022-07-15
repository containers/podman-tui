package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
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
