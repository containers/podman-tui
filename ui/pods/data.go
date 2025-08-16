package pods

import (
	"fmt"
	"strings"

	ppods "github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retrieves pods list data.
func (pods *Pods) UpdateData() {
	podList, err := ppods.List()
	if err != nil {
		log.Error().Msgf("view: pods update %v", err)
		pods.errorDialog.SetText(fmt.Sprintf("%v", err))
		pods.errorDialog.Display()

		return
	}

	pods.podsList.mu.Lock()
	defer pods.podsList.mu.Unlock()

	pods.podsList.report = podList
}

func (pods *Pods) getData() []*entities.ListPodsReport {
	pods.podsList.mu.Lock()
	defer pods.podsList.mu.Unlock()

	data := pods.podsList.report

	return data
}

// ClearData clears table data.
func (pods *Pods) ClearData() {
	pods.podsList.mu.Lock()
	defer pods.podsList.mu.Unlock()

	pods.podsList.report = nil

	pods.table.Clear()

	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor

	for i := range pods.headers {
		pods.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(pods.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(bgColor).
													SetTextColor(fgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	pods.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(pods.title)))
}
