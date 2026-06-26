package volumes

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/volumes"
)

// Prune removes all unused volumes.
func Prune() ([]string, error) {
	log.Debug().Msg("pdcs: podman volume prune")

	var (
		report   []string
		respData string
	)

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := volumes.Prune(conn, new(volumes.PruneOptions))
	if err != nil {
		return report, err
	}

	for _, r := range response {
		if r.Err != nil {
			respData = fmt.Sprintf("error removing %s: %s", r.Id, r.Err.Error())
			report = append(report, respData)
		}
	}

	return report, nil
}
