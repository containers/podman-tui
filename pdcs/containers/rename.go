package containers

import (
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Rename renames existing container's name
func Rename(id string, name string) error {
	log.Debug().Msgf("pdcs: podman container rename %s -> %s", id, name)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}

	return containers.Rename(conn, id, new(containers.RenameOptions).WithName(name))
}
