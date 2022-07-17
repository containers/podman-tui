package containers

import (
	"sort"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/rs/zerolog/log"
)

// List returns list of containers information.
func List() ([]entities.ListContainer, error) {
	log.Debug().Msg("pdcs: podman container ls")

	conn, err := registry.GetConnection()
	if err != nil {
		return nil, err
	}

	response, err := containers.List(conn, new(containers.ListOptions).WithAll(true))
	if err != nil {
		return nil, err
	}

	sort.Sort(containerListSortedName{response})

	log.Debug().Msgf("pdcs: %v", response)

	return response, nil
}

type lprSort []entities.ListContainer

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type containerListSortedName struct{ lprSort }

func (a containerListSortedName) Less(i, j int) bool {
	return a.lprSort[i].Names[0] < a.lprSort[j].Names[0]
}
