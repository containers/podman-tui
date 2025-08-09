package volumes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rs/zerolog/log"
)

var errNoVolume = errors.New("there is no volume to perform command")

func (vols *Volumes) runCommand(cmd string) {
	switch cmd {
	case "create":
		vols.createDialog.Display()
	case "inspect":
		vols.inspect()
	case "prune":
		vols.prunePrep()
	case "rm":
		vols.removePrep()
	}
}

func (vols *Volumes) displayError(title string, err error) {
	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	vols.errorDialog.SetTitle(strings.ToUpper(title))
	vols.errorDialog.SetText(fmt.Sprintf("%v", err))
	vols.errorDialog.Display()
}

func (vols *Volumes) create() {
	createOpts := vols.createDialog.VolumeCreateOptions()

	report, err := volumes.Create(createOpts)
	if err != nil {
		vols.displayError("volume create error", err)

		return
	}

	vols.messageDialog.SetTitle("podman volume create")
	vols.messageDialog.SetText(dialogs.MessageVolumeInfo, createOpts.Name, report)
	vols.messageDialog.Display()
}

func (vols *Volumes) inspect() {
	volID := vols.getSelectedItem()
	if volID == "" {
		vols.displayError("", errNoVolume)

		return
	}

	data, err := volumes.Inspect(volID)
	if err != nil {
		title := fmt.Sprintf("volume (%s) inspect error", volID)
		vols.displayError(title, err)

		return
	}

	vols.messageDialog.SetTitle("podman volume inspect")
	vols.messageDialog.SetText(dialogs.MessageVolumeInfo, volID, data)
	vols.messageDialog.Display()
}

func (vols *Volumes) prunePrep() {
	vols.confirmDialog.SetTitle("podman volume prune")
	vols.confirmData = utils.PruneCommandLabel
	vols.confirmDialog.SetText("Are you sure you want to remove all unused volumes ?")
	vols.confirmDialog.Display()
}

func (vols *Volumes) prune() {
	vols.progressDialog.SetTitle("VOLUME prune in progress")
	vols.progressDialog.Display()

	prune := func() {
		errData, err := volumes.Prune()

		vols.progressDialog.Hide()

		errorTitle := "volume prune error"
		if err != nil {
			vols.displayError(errorTitle, err)
			vols.appFocusHandler()

			return
		}

		if len(errData) > 0 {
			pruneError := errors.New(strings.Join(errData, "\n")) //nolint:err113
			vols.displayError(errorTitle, pruneError)
			vols.appFocusHandler()
		}
	}

	go prune()
}

func (vols *Volumes) removePrep() {
	volID := vols.getSelectedItem()
	if volID == "" {
		vols.displayError("", errNoVolume)

		return
	}

	vols.confirmDialog.SetTitle("podman pod rm")

	vols.confirmData = "rm"
	bgColor := style.GetColorHex(style.DialogBorderColor)
	fgColor := style.GetColorHex(style.DialogFgColor)
	volumeItem := fmt.Sprintf("[%s:%s:b]VOLUME NAME:[:-:-] %s", fgColor, bgColor, volID)
	description := fmt.Sprintf("%s\n\nAre you sure you want to remove the selected volume?", //nolint:perfsprint
		volumeItem)

	vols.confirmDialog.SetText(description)
	vols.confirmDialog.Display()
}

func (vols *Volumes) remove() {
	volID := vols.getSelectedItem()
	if volID == "" {
		vols.displayError("", errNoVolume)

		return
	}

	vols.progressDialog.SetTitle("volume remove in progress")
	vols.progressDialog.Display()

	remove := func(name string) {
		err := volumes.Remove(name)

		vols.progressDialog.Hide()

		if err != nil {
			title := fmt.Sprintf("volume (%s) remove error", volID)
			vols.displayError(title, err)
			vols.appFocusHandler()

			return
		}
	}

	go remove(volID)
}
