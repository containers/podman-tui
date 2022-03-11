package containers

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Status returns status of a container
func Status(id string) (string, error) {
	log.Debug().Msg("pdcs: podman container ls")
	conn, err := connection.GetConnection()
	if err != nil {
		return "", err
	}

	filter := make(map[string][]string)
	filter["id"] = []string{id}
	response, err := containers.List(conn, new(containers.ListOptions).WithFilters(filter))
	if err != nil {
		return "", err
	}
	if len(response) == 0 {
		return "", err
	}
	log.Debug().Msgf("pdcs: %v", response)
	return response[0].State, nil
}
