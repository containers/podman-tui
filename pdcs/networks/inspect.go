package networks

import (
	"github.com/containers/podman/v3/pkg/bindings/network"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/rs/zerolog/log"
)

// Inspect inspects the specified network
func Inspect(id string) (string, error) {
	log.Debug().Msgf("pdcs: podman network inspect %s", id)

	var report string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := network.Inspect(conn, id, new(network.InspectOptions))
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
