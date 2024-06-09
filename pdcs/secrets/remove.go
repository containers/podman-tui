package secrets

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/secrets"
	"github.com/rs/zerolog/log"
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
