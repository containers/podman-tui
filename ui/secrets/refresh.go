package secrets

import (
	"fmt"
	"strings"
	"time"

	"github.com/containers/podman-tui/ui/style"
	"github.com/docker/go-units"
	"github.com/rivo/tview"
)

func (s *Secrets) refresh(_ int) {
	s.table.Clear()

	expand := 1
	alignment := tview.AlignLeft

	for i := range s.headers {
		s.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(s.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(style.PageHeaderBgColor).
													SetTextColor(style.PageHeaderFgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	rowIndex := 1
	secResponse := s.getData()

	s.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(s.title), len(secResponse)))

	for i := range secResponse {
		secID := secResponse[i].ID
		secName := secResponse[i].Spec.Name
		secDriver := secResponse[i].Spec.Driver.Name
		secCreated := units.HumanDuration(time.Since(secResponse[i].CreatedAt)) + " ago"
		secUpdated := units.HumanDuration(time.Since(secResponse[i].UpdatedAt)) + " ago"

		// ID column
		s.table.SetCell(rowIndex, viewSecretsIDColIndex,
			tview.NewTableCell(secID).
				SetExpansion(expand).
				SetAlign(alignment))

		// Name column
		s.table.SetCell(rowIndex, viewSecretsNameColIndex,
			tview.NewTableCell(secName).
				SetExpansion(expand).
				SetAlign(alignment))

		// Driver column
		s.table.SetCell(rowIndex, viewSecretsDriverColIndex,
			tview.NewTableCell(secDriver).
				SetExpansion(expand).
				SetAlign(alignment))

		// Created column
		s.table.SetCell(rowIndex, viewSecretsCreatedColIndex,
			tview.NewTableCell(secCreated).
				SetExpansion(expand).
				SetAlign(alignment))

		// Updated column
		s.table.SetCell(rowIndex, viewSecretsUpdatedColIndex,
			tview.NewTableCell(secUpdated).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
}
