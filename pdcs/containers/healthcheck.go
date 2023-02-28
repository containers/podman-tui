package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// HealthCheck runs the health check of a container.
func HealthCheck(id string) (string, error) {
	log.Debug().Msgf("pdcs: podman container healthcheck %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return "", err
	}

	options := new(containers.HealthCheckOptions)

	response, err := containers.RunHealthCheck(conn, id, options)
	if err != nil {
		return "", err
	}

	report, err := utils.GetJSONOutput(response)
	if err != nil {
		return "", err
	}

	return report, nil
}
