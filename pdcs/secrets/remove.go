package secrets

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/secrets"
)

// Remove removes the secret.
func Remove(id string) error {
	log.Debug().Msgf("pdcs: podman secret remove %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return secrets.Remove(conn, id)
}
