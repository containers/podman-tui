package secrets

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/secrets"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rs/zerolog/log"
)

func (s *Secrets) runCommand(cmd string) {
	switch cmd {
	case "create":
		s.createDialog.Display()
	case "inspect":
		s.inspect()
	case "rm":
		s.rm()
	}
}

func (s *Secrets) displayError(title string, err error) {
	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	s.errorDialog.SetTitle(title)
	s.errorDialog.SetText(fmt.Sprintf("%v", err))
	s.errorDialog.Display()
}

func (s *Secrets) create() {
	createOpts := s.createDialog.GetCreateOptions()

	if createOpts.File != "" && createOpts.Text != "" {
		s.displayError("SECRET CREATE ERROR", errSecretFileAndText)

		return
	}

	if createOpts.File == "" && createOpts.Text == "" {
		s.displayError("SECRET CREATE ERROR", errEmptySecretFileOrText)

		return
	}

	if err := secrets.Create(createOpts); err != nil {
		s.displayError("SECRET CREATE ERROR", err)

		return
	}

	s.UpdateData()
}

func (s *Secrets) inspect() {
	_, secID, secName := s.getSelectedItem()
	if secID == "" {
		s.displayError("", errNoSecretInspect)

		return
	}

	data, err := secrets.Inspect(secID)
	if err != nil {
		title := fmt.Sprintf("SECRET (%s) INSPECT ERROR", secID)
		s.displayError(title, err)

		return
	}

	headerLabel := fmt.Sprintf("%s (%s)", secID, secName)

	s.messageDialog.SetTitle("podman secret inspect")
	s.messageDialog.SetText(dialogs.MessageSecretInfo, headerLabel, data)
	s.messageDialog.DisplayFullSize()
}

func (s *Secrets) rm() {
	_, secID, secName := s.getSelectedItem()
	if secID == "" {
		s.displayError("", errNoSecretRemove)

		return
	}

	s.confirmDialog.SetTitle("podman secret remove")

	bgColor := style.GetColorHex(style.DialogBorderColor)
	fgColor := style.GetColorHex(style.DialogFgColor)
	networkItem := fmt.Sprintf("[%s:%s:b]SECRET ID:[:-:-] %s (%s)", fgColor, bgColor, secID, secName)

	description := fmt.Sprintf("%s\n\nAre you sure you want to remove the selected secret?", //nolint:perfsprint
		networkItem)
	s.confirmDialog.SetText(description)
	s.confirmDialog.Display()
}

func (s *Secrets) remove() {
	rowIndex, secID, _ := s.getSelectedItem()
	if secID == "" {
		s.displayError("", errNoSecretRemove)

		return
	}

	s.progressDialog.SetTitle("secret remove in progress")
	s.progressDialog.Display()

	remove := func(id string) {
		err := secrets.Remove(id)

		s.progressDialog.Hide()

		if err != nil {
			title := fmt.Sprintf("SECRET (%s) REMOVE ERROR", secID)
			s.displayError(title, err)

			return
		}

		rowIndex--
		if rowIndex > 0 {
			s.table.Select(rowIndex, 0)
		}

		s.UpdateData()
	}

	go remove(secID)
}
