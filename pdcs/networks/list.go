package networks

import (
	"sort"
	"strings"

	"github.com/containers/podman/v3/pkg/bindings/network"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman-tui/pdcs/connection"
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
		var plugins []string
		for _, p := range item.Plugins {
			plugins = append(plugins, p.Network.Type)
		}
		report = append(report, []string{
			item.Name,
			item.CNIVersion,
			strings.Join(plugins, ","),
		})
	}

	log.Debug().Msgf("pdcs: %v", report)

	return report, nil
}

type lprSort []*entities.NetworkListReport

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type netListSortedName struct{ lprSort }

func (a netListSortedName) Less(i, j int) bool { return a.lprSort[i].Name < a.lprSort[j].Name }
