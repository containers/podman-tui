package containers

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Remove removes the container.
func Remove(id string) error {
	log.Debug().Msgf("pdcs: podman container remove %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	return containers.Remove(conn, id, new(containers.RemoveOptions))
}
