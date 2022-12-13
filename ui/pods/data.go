package pods

import (
	"fmt"
	"strings"

	ppods "github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retrieves pods list data
func (pods *Pods) UpdateData() {
	podList, err := ppods.List()
	if err != nil {
		log.Error().Msgf("view: pods update %v", err)
		pods.errorDialog.SetText(fmt.Sprintf("%v", err))
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
	pods.podsList.mu.Lock()
	pods.podsList.report = nil
	pods.podsList.mu.Unlock()
	pods.table.Clear()
	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor

	for i := 0; i < len(pods.headers); i++ {
		pods.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(pods.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	pods.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(pods.title)))

}
