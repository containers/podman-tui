package system

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
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
	case "prune": //nolint:goconst
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
	selectedItem := sys.getSelectedItem()
	// empty table
	if selectedItem.name == "" {
		return
	}

	dest := registry.Connection{
		Name:     selectedItem.name,
		URI:      selectedItem.uri,
		Identity: selectedItem.identity,
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
		sys.dfDialog.SetServiceName(connName)
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
	sys.eventDialog.SetServiceName(connName)
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

	sys.messageDialog.SetTitle("SYSTEM INFORMATION")
	sys.messageDialog.SetText(dialogs.MessageSystemInfo, connName, data)
	sys.messageDialog.Display()
}

func (sys *System) cprune() {
	if !sys.destIsSet() {
		return
	}

	connName := registry.ConnectionName()

	sys.confirmDialog.SetTitle("podman system prune")
	sys.confirmData = "prune"
	confirmMsg := fmt.Sprintf(
		"Are you sure you want to remove all unused pod, container, image and volume data on %s?",
		connName)
	sys.confirmDialog.SetText(confirmMsg)
	sys.confirmDialog.Display()
}

func (sys *System) prune() {
	sys.progressDialog.SetTitle("system prune in progress")
	sys.progressDialog.Display()

	prune := func() {
		report, err := sysinfo.Prune()

		sys.progressDialog.Hide()

		if err != nil {
			sys.displayError("SYSTEM PRUNE ERROR", err)

			return
		}

		sys.messageDialog.SetTitle("PODMAN SYSTEM PRUNE")
		sys.messageDialog.SetText(dialogs.MessageSystemInfo, registry.ConnectionName(), report)
		sys.messageDialog.Display()
	}

	go prune()
}

func (sys *System) cremove() {
	selectedItem := sys.getSelectedItem()
	if selectedItem.status != "" {
		sys.displayError(
			"SYSTEM CONNECTION REMOVE",
			fmt.Errorf("%w %q", ErrConnectionInprogres, selectedItem.name))

		return
	}

	if selectedItem.name == "" {
		return
	}

	title := "podman system connection remove"
	sys.confirmDialog.SetTitle(title)
	sys.confirmData = "remove_conn"
	bgColor := style.GetColorHex(style.DialogBorderColor)
	fgColor := style.GetColorHex(style.DialogFgColor)
	serviceItem := fmt.Sprintf("[%s:%s:b]SERVICE NAME:[:-:-] %s", fgColor, bgColor, selectedItem.name)

	confirmMsg := fmt.Sprintf("%s\n\nAre you sure you want to remove the selected service connection ?", //nolint:perfsprint,lll
		serviceItem)
	sys.confirmDialog.SetText(confirmMsg)
	sys.confirmDialog.Display()
}

func (sys *System) remove() {
	selectedItem := sys.getSelectedItem()

	sys.progressDialog.SetTitle("removing connection")
	sys.progressDialog.Display()

	go func() {
		err := sys.connectionRemoveFunc(selectedItem.name)
		sys.progressDialog.Hide()

		if err != nil {
			sys.displayError("SYSTEM CONNECTION REMOVE ERROR", err)

			return
		}

		sys.UpdateConnectionsData()
	}()
}

func (sys *System) setDefault() {
	selectedItem := sys.getSelectedItem()
	setDefFunc := func() {
		sys.progressDialog.Hide()

		if err := sys.connectionSetDefaultFunc(selectedItem.name); err != nil {
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
