package images

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/images"
)

// Pull pulls image from registry.
func Pull(name string) error {
	log.Debug().Msgf("pdcs: podman image pull %s", name)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	_, err = images.Pull(conn, name, new(images.PullOptions).WithQuiet(true))
	if err != nil {
		return err
	}

	return nil
}
