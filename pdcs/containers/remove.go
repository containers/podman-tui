package containers

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Remove removes the container.
func Remove(id string) ([]string, error) {
	log.Debug().Msgf("pdcs: podman container remove %s", id)
	var report []string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := containers.Remove(conn, id, new(containers.RemoveOptions))
	if err != nil {
		return report, err
	}
	for _, rmRept := range response {
		if rmRept != nil {
			respData := fmt.Sprintf("error removing %s: %s", rmRept.Id, rmRept.Err.Error())
			report = append(report, respData)
		}
	}
	return report, err
}
