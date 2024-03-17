package containers

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Remove removes the container.
func Remove(id string) ([]string, error) {
	log.Debug().Msgf("pdcs: podman container remove %s", id)

	var report []string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := containers.Remove(conn, id, new(containers.RemoveOptions))
	if err != nil {
		return report, err
	}

	for _, rmRept := range response {
		if rmRept != nil {
			if rmRept.Err != nil {
				respData := fmt.Sprintf("error removing %s: %v", rmRept.Id, rmRept.Err)
				report = append(report, respData)
			}
		}
	}

	return report, err
}
