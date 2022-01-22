package system

import (
	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/rs/zerolog/log"
)

func (sys *System) runCommand(cmd string) {
	switch cmd {
	case "disk usage":
		sys.df()
	case "info":
		sys.info()
	case "prune":
		sys.cprune()
	}

}

func (sys *System) info() {
	data, err := sysinfo.Info()
	if err != nil {
		log.Error().Msgf("view: system %s", err.Error())
		sys.errorDialog.SetText(err.Error())
		sys.errorDialog.Display()
		return
	}
	sys.messageDialog.SetTitle("podman system info")
	sys.messageDialog.SetText(data)
	sys.messageDialog.Display()
}

func (sys *System) df() {
	sys.progressDialog.SetTitle("podman disk usage in progress")
	sys.progressDialog.Display()
	diskUsage := func() {
		response, err := sysinfo.DiskUsage()
		sys.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: system %s", err.Error())
			sys.errorDialog.SetText(err.Error())
			sys.errorDialog.Display()
			return
		}
		sys.dfDialog.UpdateDiskSummary(response)
		sys.dfDialog.Display()
	}
	go diskUsage()
}

func (sys *System) cprune() {
	sys.confirmDialog.SetTitle("podman system prune")
	sys.confirmData = "prune"
	sys.confirmDialog.SetText("Are you sure you want to remove all unused pod, container, image and volume data ?")
	sys.confirmDialog.Display()
}

func (sys *System) prune() {
	sys.progressDialog.SetTitle("system purne in progress")
	sys.progressDialog.Display()
	prune := func() {
		report, err := sysinfo.Prune()
		sys.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: system %s", err.Error())
			sys.errorDialog.SetText(err.Error())
			sys.errorDialog.Display()
			return
		}
		sys.messageDialog.SetText("PODMAN SYSTEM PRUNE")
		sys.messageDialog.SetText(report)
		sys.messageDialog.Display()
	}
	go prune()
}
