package system

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/rs/zerolog/log"
)

func (sys *System) runCommand(cmd string) {
	switch cmd {
	case "add connection":
		sys.connAddDialog.Display()
	case "connect":
		sys.connect()
	case "disconnect":
		sys.disconnect()
	case "disk usage":
		sys.df()
	case "events":
		sys.events()
	case "info":
		sys.info()
	case "prune":
		sys.cprune()
	case "remove connection":
		sys.cremove()
	case "set default":
		sys.setDefault()
	}
}

func (sys *System) displayError(title string, err error) {
	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	sys.errorDialog.SetTitle(title)
	sys.errorDialog.SetText(fmt.Sprintf("%v", err))
	sys.errorDialog.Display()
}

func (sys *System) addConnection() {
	sys.connAddDialog.Hide()
	name, uri, identity := sys.connAddDialog.GetItems()
	sys.progressDialog.SetTitle("adding new connection")
	sys.progressDialog.Display()
	go func() {
		err := sys.connectionAddFunc(name, uri, identity)
		sys.progressDialog.Hide()
		sys.UpdateConnectionsData()
		if err != nil {
			sys.displayError("ADD NEW CONNECTION ERROR", err)
			return
		}
	}()
}

func (sys *System) connect() {
	connName, _, connURI, connIdentity := sys.getSelectedItem()
	// empty table
	if connName == "" {
		return
	}
	dest := registry.Connection{
		Name:     connName,
		URI:      connURI,
		Identity: connIdentity,
	}
	sys.eventDialog.SetText("")
	sys.connectionConnectFunc(dest)
	sys.UpdateConnectionsData()
}

func (sys *System) disconnect() {
	sys.connectionDisconnectFunc()
	sys.eventDialog.SetText("")
	sys.UpdateConnectionsData()
}

func (sys *System) df() {
	if !sys.destIsSet() {
		return
	}
	sys.progressDialog.SetTitle("podman disk usage in progress")
	sys.progressDialog.Display()
	diskUsage := func() {
		response, err := sysinfo.DiskUsage()
		sys.progressDialog.Hide()
		if err != nil {
			sys.displayError("SYSTEM DISK USAGE ERROR", err)
			return
		}
		connName := registry.ConnectionName()
		sys.dfDialog.SetTitle(connName)
		sys.dfDialog.UpdateDiskSummary(response)
		sys.dfDialog.Display()
	}
	go diskUsage()
}

func (sys *System) events() {
	if !sys.destIsSet() {
		return
	}
	connName := registry.ConnectionName()
	sys.eventDialog.SetTitle(connName)
	sys.eventDialog.Display()
}

func (sys *System) info() {
	if !sys.destIsSet() {
		return
	}
	data, err := sysinfo.Info()
	if err != nil {
		sys.displayError("SYSTEM INFO ERROR", err)
		return
	}
	connName := registry.ConnectionName()
	title := fmt.Sprintf("%s system info", connName)
	sys.messageDialog.SetTitle(title)
	sys.messageDialog.SetText(data)
	sys.messageDialog.Display()
}

func (sys *System) cprune() {
	if !sys.destIsSet() {
		return
	}
	connName := registry.ConnectionName()
	sys.confirmDialog.SetTitle("podman system prune")
	sys.confirmData = "prune"
	confirmMsg := fmt.Sprintf("Are you sure you want to remove all unused pod, container, image and volume data on %s?", connName)
	sys.confirmDialog.SetText(confirmMsg)
	sys.confirmDialog.Display()
}

func (sys *System) prune() {
	sys.progressDialog.SetTitle("system purne in progress")
	sys.progressDialog.Display()
	prune := func() {
		report, err := sysinfo.Prune()
		sys.progressDialog.Hide()
		if err != nil {
			sys.displayError("SYSTEM PRUNE ERROR", err)
			return
		}
		sys.messageDialog.SetText("PODMAN SYSTEM PRUNE")
		sys.messageDialog.SetText(report)
		sys.messageDialog.Display()
	}
	go prune()
}

func (sys *System) cremove() {
	connName, status, _, _ := sys.getSelectedItem()
	if status != "" {
		sys.displayError("SYSTEM CONNECTION REMOVE", fmt.Errorf("%q connection in progress, need to disconnect", connName))
		return
	}
	if connName == "" {
		return
	}
	title := "podman system connection remove"
	sys.confirmDialog.SetTitle(title)
	sys.confirmData = "remove_conn"
	confirmMsg := fmt.Sprintf("Are you sure you want to remove connection %q?", connName)
	sys.confirmDialog.SetText(confirmMsg)
	sys.confirmDialog.Display()
}

func (sys *System) remove() {
	connName, _, _, _ := sys.getSelectedItem()

	sys.progressDialog.SetTitle("removing connection")
	sys.progressDialog.Display()
	go func() {
		sys.connectionRemoveFunc(connName)
		sys.progressDialog.Hide()
		sys.UpdateConnectionsData()
	}()

}

func (sys *System) setDefault() {
	connName, _, _, _ := sys.getSelectedItem()
	setDefFunc := func() {
		sys.progressDialog.Hide()
		if err := sys.connectionSetDefaultFunc(connName); err != nil {
			sys.displayError("SYSTEM CONNECTION SET DEFAULT ERROR", err)
			return
		}
	}
	sys.progressDialog.Display()
	go setDefFunc()
}

func (sys *System) destIsSet() bool {
	if !registry.ConnectionIsSet() {
		sys.errorDialog.SetText("not connected to any podman service")
		sys.errorDialog.Display()
		return false
	}
	return true
}
