package containers

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Prune removes all non running containers
func Prune() ([]string, error) {
	log.Debug().Msgf("pdcs: podman container purne")
	var report []string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := containers.Prune(conn, new(containers.PruneOptions))
	if err != nil {
		return report, err
	}
	for _, r := range response {
		respData := ""
		if r.Err != nil {
			respData = fmt.Sprintf("error removing %s: %s", r.Id, r.Err.Error())
			report = append(report, respData)
		}

	}
	return report, nil
}
