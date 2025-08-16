package containers

import (
	"strings"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rs/zerolog/log"
)

// Port retrurn ports mapping of the container.
func Port(id string) ([]string, error) {
	log.Debug().Msgf("pdcs: podman container port %s", id)

	var report []string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	filters := make(map[string][]string)
	filters["id"] = []string{id}

	response, err := containers.List(conn, new(containers.ListOptions).WithFilters(filters))
	if err != nil {
		return report, err
	}

	if len(response) > 0 {
		report = strings.Split(conReporter{response[0]}.ports(), ",")
	}

	log.Debug().Msgf("pdcs: %v", report)

	return report, nil
}

type conReporter struct {
	entities.ListContainer
}

func (con conReporter) ports() string {
	if len(con.Ports) < 1 {
		return ""
	}

	return utils.PortsToString(con.Ports)
}
