package system

import (
	"fmt"
	"sort"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rs/zerolog/log"
)

type connectionItemStatus struct {
	status registry.ConnStatus
}

func (connStatus connectionItemStatus) StatusString() string {
	var status string

	switch connStatus.status {
	case registry.ConnectionStatusConnected:
		status = fmt.Sprintf("%s %s", style.HeavyGreenCheckMark, "connected")
	case registry.ConnectionStatusConnectionError:
		status = fmt.Sprintf("%s %s", style.HeavyRedCrossMark, "connection error")
	}

	return status
}

// SortView sorts data view called from sort dialog.
func (sys *System) SortView(option string, ascending bool) {
	log.Debug().Msgf("view: system sort by %s", option)

	sys.connectionList.mu.Lock()
	defer sys.connectionList.mu.Unlock()

	sys.connectionList.sortBy = option
	sys.connectionList.ascending = ascending
	sort.Sort(ConnectionListSorted{sys.connectionList.report, option, ascending})
}

// UpdateData retrieves connections list data.
func (sys *System) UpdateData() {
	destinations := sys.connectionListFunc()

	sys.connectionList.mu.Lock()

	sort.Sort(ConnectionListSorted{destinations, sys.connectionList.sortBy, sys.connectionList.ascending})
	sys.connectionList.report = destinations

	sys.connectionList.mu.Unlock()

	sys.udpateConnectionDataStatus()
}

func (sys *System) udpateConnectionDataStatus() {
	sys.connectionList.mu.Lock()
	defer sys.connectionList.mu.Unlock()

	name := registry.ConnectionName()
	status := registry.ConnectionStatus()

	for i := range sys.connectionList.report {
		if sys.connectionList.report[i].Name == name {
			sys.connectionList.report[i].Status = status

			return
		}
	}
}

func (sys *System) getConnectionsData() []registry.Connection {
	sys.connectionList.mu.Lock()
	defer sys.connectionList.mu.Unlock()

	destReport := sys.connectionList.report

	return destReport
}

type lprSort []registry.Connection

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ConnectionListSorted struct {
	lprSort

	option    string
	ascending bool
}

func (a ConnectionListSorted) Less(i, j int) bool { //nolint:cyclop
	switch a.option {
	case "default":
		if a.ascending {
			return a.lprSort[i].Default || a.lprSort[j].Default
		}

		return !a.lprSort[i].Default && a.lprSort[j].Default
	case "status":
		if a.ascending {
			return a.lprSort[i].Status.String() < a.lprSort[j].Status.String()
		}

		return a.lprSort[i].Status.String() > a.lprSort[j].Status.String()
	case "uri":
		if a.ascending {
			return a.lprSort[i].URI < a.lprSort[j].URI
		}

		return a.lprSort[i].URI > a.lprSort[j].URI
	case "identity":
		if a.ascending {
			return a.lprSort[i].Identity < a.lprSort[j].Identity
		}

		return a.lprSort[i].Identity > a.lprSort[j].Identity
	}

	if a.ascending {
		return a.lprSort[i].Name < a.lprSort[j].Name
	}

	return a.lprSort[i].Name > a.lprSort[j].Name
}
