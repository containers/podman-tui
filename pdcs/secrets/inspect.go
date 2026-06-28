package secrets

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/secrets"
)

// Inspect inspects the specified secret.
func Inspect(id string) (string, error) {
	log.Debug().Msgf("pdcs: podman secret inspect %s", id)

	var report string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := secrets.Inspect(conn, id, new(secrets.InspectOptions).WithShowSecret(true))
	if err != nil {
		return report, err
	}

	report, err = utils.GetJSONOutput(response)
	if err != nil {
		return report, err
	}

	return report, nil
}
