package containers

import (
	"fmt"
	"strings"
	"time"

	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/docker/go-units"
	"github.com/containers/podman-tui/pdcs/containers"
	putils "github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retreives containers list data
func (cnt *Containers) UpdateData() {
	cntList, err := containers.List()
	if err != nil {
		log.Error().Msgf("view: containers %s", err.Error())
		cnt.errorDialog.SetText(err.Error())
		cnt.errorDialog.Display()
		return
	}
	cnt.containersList.mu.Lock()
	cnt.containersList.report = cntList
	cnt.containersList.mu.Unlock()
}

func (cnt *Containers) getData() []entities.ListContainer {
	cnt.containersList.mu.Lock()
	data := cnt.containersList.report
	cnt.containersList.mu.Unlock()
	return data
}

// ClearData clears table data
func (cnt *Containers) ClearData() {
	cnt.table.Clear()
	expand := 1
	fgColor := utils.Styles.PageTable.HeaderRow.FgColor
	bgColor := utils.Styles.PageTable.HeaderRow.BgColor
	for i := 0; i < len(cnt.headers); i++ {
		cnt.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(cnt.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	cnt.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(cnt.title)))
}

type conReporter struct {
	entities.ListContainer
}

func (con conReporter) names() string {
	return strings.Join(con.Names, ",")
}

func (con conReporter) state() string {
	var state string
	switch con.ListContainer.State {
	case "running":
		t := units.HumanDuration(time.Since(time.Unix(con.StartedAt, 0)))
		state = "Up " + t + " ago"
	case "configured":
		state = "Created"
	case "exited", "stopped":
		t := units.HumanDuration(time.Since(time.Unix(con.ExitedAt, 0)))
		state = fmt.Sprintf("Exited (%d) %s ago", con.ExitCode, t)
	default:
		state = con.ListContainer.State
	}
	return state
}

func (con conReporter) status() string {
	hc := con.ListContainer.Status
	if hc != "" {
		return con.state() + " (" + hc + ")"
	}
	return con.state()
}

func (con conReporter) ports() string {
	if len(con.ListContainer.Ports) < 1 {
		return ""
	}
	return putils.PortsToString(con.ListContainer.Ports)
}
