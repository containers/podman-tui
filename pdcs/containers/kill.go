package containers

import (
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Kill sends SIGKILL signal to container processes
func Kill(id string) error {
	log.Debug().Msgf("pdcs: podman container kill %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	return containers.Kill(conn, id, new(containers.KillOptions))

}
