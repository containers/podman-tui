package pods

import (
	"fmt"

	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Prune removes all stop pods
func Prune() ([]string, error) {
	log.Debug().Msgf("pdcs: podman pod purne")
	var report []string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := pods.Prune(conn, new(pods.PruneOptions))
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
