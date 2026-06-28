package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/containers"
)

// Kill sends SIGKILL signal to container processes.
func Kill(id string) error {
	log.Debug().Msgf("pdcs: podman container kill %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return containers.Kill(conn, id, new(containers.KillOptions))
}
