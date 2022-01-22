package pods

import (
	"fmt"
	"strings"

	"github.com/containers/podman/v3/pkg/domain/entities"
	ppods "github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retreives pods list data
func (pods *Pods) UpdateData() {
	podList, err := ppods.List()
	if err != nil {
		log.Error().Msgf("view: pods %s", err.Error())
		pods.errorDialog.SetText(err.Error())
		pods.errorDialog.Display()
		return
	}
	pods.podsList.mu.Lock()
	pods.podsList.report = podList
	pods.podsList.mu.Unlock()
}

func (pods *Pods) getData() []*entities.ListPodsReport {
	pods.podsList.mu.Lock()
	data := pods.podsList.report
	pods.podsList.mu.Unlock()
	return data
}

// ClearData clears table data
func (pods *Pods) ClearData() {
	pods.table.Clear()
	expand := 1
	fgColor := utils.Styles.PageTable.HeaderRow.FgColor
	bgColor := utils.Styles.PageTable.HeaderRow.BgColor

	for i := 0; i < len(pods.headers); i++ {
		pods.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(pods.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	pods.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(pods.title)))

}
