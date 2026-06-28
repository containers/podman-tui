package volumes

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/volumes"
)

// Remove removes the volulme.
func Remove(name string) error {
	log.Debug().Msgf("pdcs: podman volume remove %s", name)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return volumes.Remove(conn, name, new(volumes.RemoveOptions))
}
