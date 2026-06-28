package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/containers"
)

// Rename renames existing container's name.
func Rename(id string, name string) error {
	log.Debug().Msgf("pdcs: podman container rename %s -> %s", id, name)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return containers.Rename(conn, id, new(containers.RenameOptions).WithName(name))
}
