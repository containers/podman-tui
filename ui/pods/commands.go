package pods

import (
	"fmt"
	"strings"

	ppods "github.com/containers/podman-tui/pdcs/pods"
	"github.com/rs/zerolog/log"
)

func (p *Pods) runCommand(cmd string) {
	switch cmd {
	case "create":
		p.createDialog.Display()
	case "inspect":
		p.inspect()
	case "kill":
		p.kill()
	case "pause":
		p.pause()
	case "prune":
		p.confirmDialog.SetTitle("podman pod prune")
		p.confirmData = "prune"
		p.confirmDialog.SetText("Are you sure you want to remove all stopped pods ?")
		p.confirmDialog.Display()
	case "restart":
		p.restart()
	case "rm":
		p.rm()
	case "start":
		p.start()
	case "stop":
		p.stop()
	case "top":
		p.top()
	case "unpause":
		p.unpause()
	}
}

func (p *Pods) create() {
	podSpec := p.createDialog.GetPodSpec()
	err := ppods.Create(podSpec)
	if err != nil {
		log.Error().Msgf("view: pods create %s", err.Error())
		p.errorDialog.SetText(err.Error())
		p.errorDialog.Display()
		return
	}
}

func (p *Pods) inspect() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to inspect")
		p.errorDialog.Display()
		return
	}
	data, err := ppods.Inspect(p.selectedID)
	if err != nil {
		log.Error().Msgf("view: pods %s", err.Error())
		p.errorDialog.SetText(err.Error())
		p.errorDialog.Display()
		return
	}
	p.messageDialog.SetTitle("podman pod inspect")
	p.messageDialog.SetText(data)
	p.messageDialog.Display()
}

func (p *Pods) kill() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to kill")
		p.errorDialog.Display()
		return
	}
	p.progressDialog.SetTitle("pod kill in progress")
	p.progressDialog.Display()
	kill := func(id string) {
		err := ppods.Kill(id)
		p.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: pods %s", err.Error())
			p.errorDialog.SetText(err.Error())
			p.errorDialog.Display()
			return
		}
	}
	go kill(p.selectedID)
}

func (p *Pods) pause() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to pause")
		p.errorDialog.Display()
		return
	}
	p.progressDialog.SetTitle("pod pause in progress")
	p.progressDialog.Display()
	pause := func(id string) {
		err := ppods.Pause(id)
		p.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: pods %s", err.Error())
			p.errorDialog.SetText(err.Error())
			p.errorDialog.Display()
			return
		}
	}
	go pause(p.selectedID)
}

func (p *Pods) prune() {
	p.progressDialog.SetTitle("pod purne in progress")
	p.progressDialog.Display()
	unpause := func() {
		errData, err := ppods.Prune()
		p.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: pods %s", err.Error())
			p.errorDialog.SetText(err.Error())
			p.errorDialog.Display()
			return
		}
		if len(errData) > 0 {
			p.errorDialog.SetText(strings.Join(errData, "\n"))
			p.errorDialog.Display()
		}

	}
	go unpause()
}

func (p *Pods) restart() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to restart")
		p.errorDialog.Display()
		return
	}
	p.progressDialog.SetTitle("pod restart in progress")
	p.progressDialog.Display()
	restart := func(id string) {
		err := ppods.Restart(id)
		p.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: pods %s", err.Error())
			p.errorDialog.SetText(err.Error())
			p.errorDialog.Display()
			return
		}
	}
	go restart(p.selectedID)
}

func (p *Pods) rm() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to remove")
		p.errorDialog.Display()
		return
	}
	p.confirmDialog.SetTitle("podman pod rm")
	p.confirmData = "rm"
	description := fmt.Sprintf("Are you sure you want to remove following pod ? \n\nPOD ID : %s", p.selectedID)
	p.confirmDialog.SetText(description)
	p.confirmDialog.Display()
}

func (p *Pods) remove() {
	p.progressDialog.SetTitle("pod remove in progress")
	p.progressDialog.Display()
	remove := func(id string) {
		errData, err := ppods.Remove(id)
		p.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: pods %s", err.Error())
			p.errorDialog.SetText(err.Error())
			p.errorDialog.Display()
			return
		}
		if len(errData) > 0 {
			p.errorDialog.SetText(strings.Join(errData, "\n"))
			p.errorDialog.Display()
		}
	}
	go remove(p.selectedID)
}

func (p *Pods) start() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to start")
		p.errorDialog.Display()
		return
	}
	p.progressDialog.SetTitle("pod start in progress")
	p.progressDialog.Display()
	start := func(id string) {
		err := ppods.Start(id)
		p.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: pods %s", err.Error())
			p.errorDialog.SetText(err.Error())
			p.errorDialog.Display()
			return
		}
	}
	go start(p.selectedID)
}

func (p *Pods) stop() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to stop")
		p.errorDialog.Display()
		return
	}
	p.progressDialog.SetTitle("pod stop in progress")
	p.progressDialog.Display()
	stop := func(id string) {
		err := ppods.Stop(id)
		p.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: pods %s", err.Error())
			p.errorDialog.SetText(err.Error())
			p.errorDialog.Display()
			return
		}
	}
	go stop(p.selectedID)
}

func (p *Pods) top() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to display top")
		p.errorDialog.Display()
		return
	}
	data, err := ppods.Top(p.selectedID)
	if err != nil {
		log.Error().Msgf("view: pods %s", err.Error())
		p.errorDialog.SetText(err.Error())
		p.errorDialog.Display()
		return
	}
	p.topDialog.UpdateResults(data)
	p.topDialog.Display()
}

func (p *Pods) unpause() {
	if p.selectedID == "" {
		p.errorDialog.SetText("there is no pod to unpause")
		p.errorDialog.Display()
		return
	}
	p.progressDialog.SetTitle("pod unpause in progress")
	p.progressDialog.Display()
	unpause := func(id string) {
		err := ppods.Unpause(id)
		p.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: pods %s", err.Error())
			p.errorDialog.SetText(err.Error())
			p.errorDialog.Display()
			return
		}
	}
	go unpause(p.selectedID)
}
