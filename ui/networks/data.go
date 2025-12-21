package networks

import (
	"fmt"
	"sort"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"go.podman.io/common/libnetwork/types"
)

// SortView sorts data view called from sort dialog.
func (nets *Networks) SortView(option string, ascending bool) {
	log.Debug().Msgf("view: networks sort by %s", option)

	nets.networkList.mu.Lock()
	defer nets.networkList.mu.Unlock()

	nets.networkList.sortBy = option
	nets.networkList.ascending = ascending

	sort.Sort(netsListSorted{nets.networkList.report, option, ascending})
}

// UpdateData retrieves networks list data.
func (nets *Networks) UpdateData() {
	netList, err := networks.List()
	if err != nil {
		log.Error().Msgf("view: networks update %v", err)
		nets.errorDialog.SetText(fmt.Sprintf("%v", err))
		nets.errorDialog.Display()
	}

	nets.networkList.mu.Lock()
	defer nets.networkList.mu.Unlock()

	sort.Sort(netsListSorted{netList, nets.networkList.sortBy, nets.networkList.ascending})

	nets.networkList.report = netList
}

func (nets *Networks) getData() []types.Network {
	nets.networkList.mu.Lock()
	defer nets.networkList.mu.Unlock()

	data := nets.networkList.report

	return data
}

// ClearData clears table data.
func (nets *Networks) ClearData() {
	nets.networkList.mu.Lock()
	defer nets.networkList.mu.Unlock()

	nets.networkList.report = nil

	nets.table.Clear()

	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor

	for i := range nets.headers {
		nets.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(nets.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(bgColor).
													SetTextColor(fgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	nets.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(nets.title)))
}

type lprSort []types.Network

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type netsListSorted struct {
	lprSort

	option    string
	ascending bool
}

func (a netsListSorted) Less(i, j int) bool {
	if a.option == "driver" {
		if a.ascending {
			return a.lprSort[i].Driver < a.lprSort[j].Driver
		}

		return a.lprSort[i].Driver > a.lprSort[j].Driver
	}

	if a.ascending {
		return a.lprSort[i].Name < a.lprSort[j].Name
	}

	return a.lprSort[i].Name > a.lprSort[j].Name
}
