package containers

import (
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Start starts a container
func Start(id string) error {
	log.Debug().Msgf("pdcs: podman container start %s", id)

	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}

	return containers.Start(conn, id, new(containers.StartOptions))
}
