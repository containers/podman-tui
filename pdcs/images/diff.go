package images

import (
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Diff returns the diff of specified image ID
func Diff(id string) ([]string, error) {
	log.Debug().Msgf("pdcs: podman image diff %s", id)
	var report []string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := images.Diff(conn, id, new(images.DiffOptions))
	if err != nil {
		return report, err
	}

	for _, row := range response {
		report = append(report, row.String())
	}
	log.Debug().Msgf("pdcs: %s", report)
	return report, nil
}
