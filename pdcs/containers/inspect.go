package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/containers"
)

// Inspect returns inspect results of the specific container.
func Inspect(id string) (string, error) {
	log.Debug().Msgf("pdcs: podman container inspect %s", id)

	var report string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := containers.Inspect(conn, id, new(containers.InspectOptions))
	if err != nil {
		return report, err
	}

	report, err = utils.GetJSONOutput(response)
	if err != nil {
		return report, err
	}

	log.Debug().Msgf("pdcs: %s", report)

	return report, nil
}
