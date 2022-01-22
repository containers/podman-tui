package volumes

import (
	"sort"

	"github.com/containers/podman/v3/pkg/bindings/volumes"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// List returns list of volumes
func List() ([]*entities.VolumeListReport, error) {
	log.Debug().Msg("pdcs: podman volume ls")

	conn, err := connection.GetConnection()
	if err != nil {
		return nil, err
	}
	response, err := volumes.List(conn, new(volumes.ListOptions))
	if err != nil {
		return nil, err
	}
	sort.Sort(volumeListSortedName{response})
	log.Debug().Msgf("pdcs: %v", response)
	return response, nil
}

type lprSort []*entities.VolumeListReport

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type volumeListSortedName struct{ lprSort }

func (a volumeListSortedName) Less(i, j int) bool { return a.lprSort[i].Name < a.lprSort[j].Name }
