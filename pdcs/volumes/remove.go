package volumes

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/rs/zerolog/log"
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
