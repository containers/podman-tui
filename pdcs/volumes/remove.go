package volumes

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/pkg/bindings/volumes"
	"github.com/rs/zerolog/log"
)

// Remove removes the volulme
func Remove(name string) error {
	log.Debug().Msgf("pdcs: podman volume remove %s", name)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	return volumes.Remove(conn, name, new(volumes.RemoveOptions))
}
