package volumes

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retreives pods list data
func (vols *Volumes) UpdateData() {
	volList, err := volumes.List()
	if err != nil {
		log.Error().Msgf("view: volumes update %v", err)
		vols.errorDialog.SetText(fmt.Sprintf("%v", err))
		vols.errorDialog.Display()
		return
	}
	vols.volumeList.mu.Lock()
	vols.volumeList.report = volList
	vols.volumeList.mu.Unlock()
}

func (vols *Volumes) getData() []*entities.VolumeListReport {
	vols.volumeList.mu.Lock()
	data := vols.volumeList.report
	vols.volumeList.mu.Unlock()
	return data
}

// ClearData clears table data
func (vols *Volumes) ClearData() {
	vols.table.Clear()
	expand := 1
	fgColor := utils.Styles.PageTable.HeaderRow.FgColor
	bgColor := utils.Styles.PageTable.HeaderRow.BgColor

	for i := 0; i < len(vols.headers); i++ {
		vols.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(vols.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	vols.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(vols.title)))
}
