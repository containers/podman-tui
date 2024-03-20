package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Diff returns the diff of the specified container ID.
func Diff(id string) ([]string, error) {
	log.Debug().Msgf("pdcs: podman container diff %s", id)

	report := make([]string, 0)

	conn, err := registry.GetConnection()
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
