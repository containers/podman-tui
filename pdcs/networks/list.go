package networks

import (
	"sort"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/network"
	"github.com/rs/zerolog/log"
)

//List returns list of podman networks
func List() ([][]string, error) {
	log.Debug().Msg("pdcs: podman network ls")

	var report [][]string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := network.List(conn, new(network.ListOptions))
	if err != nil {
		return report, err
	}
	sort.Sort(netListSortedName{response})

	for _, item := range response {
		report = append(report, []string{
			item.ID,
			item.Name,
			item.Driver,
		})
	}

	log.Debug().Msgf("pdcs: %v", report)

	return report, nil
}

type lprSort []types.Network

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type netListSortedName struct{ lprSort }

func (a netListSortedName) Less(i, j int) bool { return a.lprSort[i].Name < a.lprSort[j].Name }
