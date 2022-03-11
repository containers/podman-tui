package containers

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Diff returns the diff of the specified container ID
func Diff(id string) ([]string, error) {
	log.Debug().Msgf("pdcs: podman container diff %s", id)
	var report []string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := containers.Diff(conn, id, new(containers.DiffOptions))
	if err != nil {
		return report, err
	}

	for _, row := range response {
		report = append(report, row.String())
	}

	log.Debug().Msgf("pdcs: %s", report)
	return report, nil
}
