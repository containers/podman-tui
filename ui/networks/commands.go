package networks

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/rs/zerolog/log"
)

func (nets *Networks) runCommand(cmd string) {
	switch cmd {
	case "create":
		nets.createDialog.Display()
	case "inspect":
		nets.inspect()
	case "prune":
		nets.cprune()
	case "rm":
		nets.rm()
	}
}

func (nets *Networks) create() {
	createOpts := nets.createDialog.NetworkCreateOptions()
	filename, err := networks.Create(createOpts)
	if err != nil {
		log.Error().Msgf("view: newtork create %s", err.Error())
		nets.errorDialog.SetText(err.Error())
		nets.errorDialog.Display()
		return
	}
	nets.UpdateData()
	nets.messageDialog.SetTitle("podman network create")
	nets.messageDialog.SetText(filename)
	nets.messageDialog.Display()

}

func (nets *Networks) inspect() {
	if nets.selectedID == "" {
		nets.errorDialog.SetText("there is no network to inspect")
		nets.errorDialog.Display()
		return
	}
	data, err := networks.Inspect(nets.selectedID)
	if err != nil {
		log.Error().Msgf("view: networks %s", err.Error())
		nets.errorDialog.SetText(err.Error())
		nets.errorDialog.Display()
		return
	}
	nets.messageDialog.SetTitle("podman network inspect")
	nets.messageDialog.SetText(data)
	nets.messageDialog.Display()
}

func (nets *Networks) cprune() {
	nets.confirmDialog.SetTitle("podman network prune")
	nets.confirmData = "prune"
	nets.confirmDialog.SetText("Are you sure you want to remove all un used network ?")
	nets.confirmDialog.Display()
}

func (nets *Networks) prune() {
	nets.progressDialog.SetTitle("network purne in progress")
	nets.progressDialog.Display()
	prune := func() {
		err := networks.Prune()
		nets.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: networks %s", err.Error())
			nets.errorDialog.SetText(err.Error())
			nets.errorDialog.Display()
			return
		}
	}
	go prune()
}

func (nets *Networks) rm() {
	if nets.selectedID == "" {
		nets.errorDialog.SetText("there is no network to remove")
		nets.errorDialog.Display()
		return
	}
	nets.confirmDialog.SetTitle("podman network remove")
	nets.confirmData = "rm"
	description := fmt.Sprintf("Are you sure you want to remove following network? \n\nNETWORK NAME : %s", nets.selectedID)
	nets.confirmDialog.SetText(description)
	nets.confirmDialog.Display()
}

func (nets *Networks) remove() {
	nets.progressDialog.SetTitle("newtork remove in progress")
	nets.progressDialog.Display()
	remove := func(id string) {
		err := networks.Remove(id)
		nets.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: networks %s", err.Error())
			nets.errorDialog.SetText(err.Error())
			nets.errorDialog.Display()
			return
		}
		nets.UpdateData()
	}
	go remove(nets.selectedID)
}
