package pods

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/pods"
	"github.com/rs/zerolog/log"
)

// Prune removes all stop pods.
func Prune() ([]string, error) {
	log.Debug().Msgf("pdcs: podman pod purne")

	var report []string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := pods.Prune(conn, new(pods.PruneOptions))
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
