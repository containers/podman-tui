package containers

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Prune removes all non running containers.
func Prune() ([]string, error) {
	log.Debug().Msgf("pdcs: podman container prune")

	var report []string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := containers.Prune(conn, new(containers.PruneOptions))
	if err != nil {
		return report, err
	}

	for _, r := range response {
		if r.Err != nil {
			respData := fmt.Sprintf("error removing %s: %s", r.Id, r.Err.Error())
			report = append(report, respData)
		}
	}

	return report, nil
}
