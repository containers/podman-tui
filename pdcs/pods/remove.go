package pods

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/pods"
	"github.com/rs/zerolog/log"
)

// Remove removes the pod.
func Remove(id string) ([]string, error) {
	log.Debug().Msgf("pdcs: podman pod remove %s", id)

	var report []string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := pods.Remove(conn, id, new(pods.RemoveOptions))
	if err != nil {
		return report, err
	}

	if response.Err != nil {
		respData := fmt.Sprintf("error removing %s: %s", response.Id, response.Err.Error())
		report = append(report, respData)
	}

	return report, nil
}
