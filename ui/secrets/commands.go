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

func (s *Secrets) inspect() {
	secID, secName := s.getSelectedItem()
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
	s.messageDialog.Display()
}

func (s *Secrets) rm() {
	secID, secName := s.getSelectedItem()
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
	secID, _ := s.getSelectedItem()
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

		s.UpdateData()
	}

	go remove(secID)
}
