package networks

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

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

	nets.networkList.report = netList
}

func (nets *Networks) getData() [][]string {
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
