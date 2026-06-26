package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

// Stats returns live stream of containers stats result.
func Stats(id string, opts *containers.StatsOptions) (chan entities.ContainerStatsReport, error) {
	log.Debug().Msgf("pdcs: podman container stats %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return nil, err
	}

	statReportChan, err := containers.Stats(conn, []string{id}, opts)
	if err != nil {
		return nil, err
	}

	return statReportChan, nil
}
