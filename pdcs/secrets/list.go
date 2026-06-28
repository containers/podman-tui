package secrets

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/secrets"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

// List returns list of podman secrets.
func List() ([]*types.SecretInfoReport, error) {
	log.Debug().Msg("pdcs: podman secret ls")

	conn, err := registry.GetConnection()
	if err != nil {
		return nil, err
	}

	response, err := secrets.List(conn, new(secrets.ListOptions))
	if err != nil {
		return nil, err
	}

	return response, nil
}
