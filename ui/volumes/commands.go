package volumes

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/rs/zerolog/log"
)

func (vols *Volumes) runCommand(cmd string) {
	switch cmd {
	case "create":
		vols.createDialog.Display()
	case "inspect":
		vols.inspect()
	case "prune":
		vols.cprune()
	case "rm":
		vols.rm()
	}
}

func (vols *Volumes) create() {
	createOpts := vols.createDialog.VolumeCreateOptions()
	report, err := volumes.Create(createOpts)
	if err != nil {
		log.Error().Msgf("view: newtork create %s", err.Error())
		vols.errorDialog.SetText(err.Error())
		vols.errorDialog.Display()
		return
	}
	vols.messageDialog.SetTitle("podman volume create")
	vols.messageDialog.SetText(report)
	vols.messageDialog.Display()
}

func (vols *Volumes) inspect() {
	if vols.selectedID == "" {
		vols.errorDialog.SetText("there is no volume to inspect")
		vols.errorDialog.Display()
		return
	}
	data, err := volumes.Inspect(vols.selectedID)
	if err != nil {
		log.Error().Msgf("view: volumes %s", err.Error())
		vols.errorDialog.SetText(err.Error())
		vols.errorDialog.Display()
		return
	}
	vols.messageDialog.SetTitle("podman volume inspect")
	vols.messageDialog.SetText(data)
	vols.messageDialog.Display()
}

func (vols *Volumes) cprune() {
	vols.confirmDialog.SetTitle("podman pod prune")
	vols.confirmData = "prune"
	vols.confirmDialog.SetText("Are you sure you want to remove all unused volumes ?")
	vols.confirmDialog.Display()
}

func (vols *Volumes) prune() {
	vols.progressDialog.SetTitle("pod purne in progress")
	vols.progressDialog.Display()
	unpause := func() {
		errData, err := volumes.Prune()
		vols.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: volumes %s", err.Error())
			vols.errorDialog.SetText(err.Error())
			vols.errorDialog.Display()
			return
		}
		if len(errData) > 0 {
			vols.errorDialog.SetText(strings.Join(errData, "\n"))
			vols.errorDialog.Display()
		}

	}
	go unpause()
}

func (vols *Volumes) rm() {
	if vols.selectedID == "" {
		vols.errorDialog.SetText("there is no volume to remove")
		vols.errorDialog.Display()
		return
	}
	vols.confirmDialog.SetTitle("podman pod rm")
	vols.confirmData = "rm"
	description := fmt.Sprintf("Are you sure you want to remove following volume ? \n\nVOLUME NAME : %s", vols.selectedID)
	vols.confirmDialog.SetText(description)
	vols.confirmDialog.Display()
}

func (vols *Volumes) remove() {
	vols.progressDialog.SetTitle("volume remove in progress")
	vols.progressDialog.Display()
	remove := func(name string) {
		err := volumes.Remove(name)
		vols.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: volumes %s", err.Error())
			vols.errorDialog.SetText(err.Error())
			vols.errorDialog.Display()
			return
		}

	}
	go remove(vols.selectedID)
}
