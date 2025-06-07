package secrets

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/secrets"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retrieves secrets list data.
func (s *Secrets) UpdateData() {
	secResponse, err := secrets.List()
	if err != nil {
		log.Error().Msgf("view: secrets update %v", err)

		s.errorDialog.SetText(fmt.Sprintf("%v", err))
		s.errorDialog.Display()
	}

	s.secretList.mu.Lock()
	defer s.secretList.mu.Unlock()

	s.secretList.report = secResponse
}

func (s *Secrets) getData() []*types.SecretInfoReport {
	s.secretList.mu.Lock()
	defer s.secretList.mu.Unlock()

	data := s.secretList.report

	return data
}

// ClearData clears table data.
func (s *Secrets) ClearData() {
	s.secretList.mu.Lock()
	defer s.secretList.mu.Unlock()

	s.secretList.report = nil

	s.table.Clear()

	expand := 1

	for i := range s.headers {
		s.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(s.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(style.PageHeaderBgColor).
													SetTextColor(style.PageHeaderFgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	s.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(s.title)))
}
