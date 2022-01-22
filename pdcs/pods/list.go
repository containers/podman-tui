package pods

import (
	"sort"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/rs/zerolog/log"
)

// List returns list of pods
func List() ([]*entities.ListPodsReport, error) {
	log.Debug().Msg("pdcs: podman pod ls")
	conn, err := connection.GetConnection()
	if err != nil {
		return nil, err
	}
	response, err := pods.List(conn, new(pods.ListOptions))
	if err != nil {
		return nil, err
	}
	sort.Sort(podPsSortedName{response})
	log.Debug().Msgf("pdcs: %v", response)
	return response, nil

}

type lprSort []*entities.ListPodsReport

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type podPsSortedName struct{ lprSort }

func (a podPsSortedName) Less(i, j int) bool { return a.lprSort[i].Name < a.lprSort[j].Name }
