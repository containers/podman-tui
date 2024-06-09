package app

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
)

func (app *App) switchToScreen(name string) {
	log.Debug().Msgf("app: switching to %s screen", name)
	app.pages.SwitchToPage(name)
	app.setPageFocus(name)
	app.updatePageData(name)

	app.currentPage = name
}

func (app *App) frontScreenHasActiveDialog() bool {
	switch app.currentPage {
	case app.system.GetTitle():
		return app.system.SubDialogHasFocus()
	case app.pods.GetTitle():
		return app.pods.SubDialogHasFocus()
	case app.containers.GetTitle():
		return app.containers.SubDialogHasFocus()
	case app.networks.GetTitle():
		return app.networks.SubDialogHasFocus()
	case app.images.GetTitle():
		return app.images.SubDialogHasFocus()
	case app.volumes.GetTitle():
		return app.volumes.SubDialogHasFocus()
	case app.secrets.GetTitle():
		return app.secrets.SubDialogHasFocus()
	}

	return false
}

func (app *App) switchToPreviousScreen() {
	var previousScreen string

	switch app.currentPage {
	case app.help.GetTitle():
		previousScreen = app.secrets.GetTitle()
	case app.system.GetTitle():
		previousScreen = app.secrets.GetTitle()
	case app.pods.GetTitle():
		previousScreen = app.system.GetTitle()
	case app.containers.GetTitle():
		previousScreen = app.pods.GetTitle()
	case app.volumes.GetTitle():
		previousScreen = app.containers.GetTitle()
	case app.images.GetTitle():
		previousScreen = app.volumes.GetTitle()
	case app.networks.GetTitle():
		previousScreen = app.images.GetTitle()
	case app.secrets.GetTitle():
		previousScreen = app.networks.GetTitle()
	}

	app.switchToScreen(previousScreen)
}

func (app *App) switchToNextScreen() {
	var nextScreen string

	switch app.currentPage {
	case app.help.GetTitle():
		nextScreen = app.system.GetTitle()
	case app.system.GetTitle():
		nextScreen = app.pods.GetTitle()
	case app.pods.GetTitle():
		nextScreen = app.containers.GetTitle()
	case app.containers.GetTitle():
		nextScreen = app.volumes.GetTitle()
	case app.volumes.GetTitle():
		nextScreen = app.images.GetTitle()
	case app.images.GetTitle():
		nextScreen = app.networks.GetTitle()
	case app.networks.GetTitle():
		nextScreen = app.secrets.GetTitle()
	case app.secrets.GetTitle():
		nextScreen = app.system.GetTitle()
	}

	app.switchToScreen(nextScreen)
}

func (app *App) setPageFocus(page string) {
	switch page {
	case app.help.GetTitle():
		app.Application.SetFocus(app.help)
	case app.system.GetTitle():
		app.Application.SetFocus(app.system)
	case app.pods.GetTitle():
		app.Application.SetFocus(app.pods)
	case app.containers.GetTitle():
		app.Application.SetFocus(app.containers)
	case app.networks.GetTitle():
		app.Application.SetFocus(app.networks)
	case app.images.GetTitle():
		app.Application.SetFocus(app.images)
	case app.volumes.GetTitle():
		app.Application.SetFocus(app.volumes)
	case app.secrets.GetTitle():
		app.Application.SetFocus(app.secrets)
	}
}

func (app *App) updatePageData(page string) {
	connStatus, _ := app.health.ConnStatus()
	if connStatus != registry.ConnectionStatusConnected {
		return
	}

	switch page {
	case app.system.GetTitle():
		app.system.UpdateConnectionsData()
	case app.pods.GetTitle():
		app.pods.UpdateData()
	case app.containers.GetTitle():
		app.containers.UpdateData()
	case app.networks.GetTitle():
		app.networks.UpdateData()
	case app.images.GetTitle():
		app.images.UpdateData()
	case app.volumes.GetTitle():
		app.volumes.UpdateData()
	case app.secrets.GetTitle():
		app.secrets.UpdateData()
	}
}

func (app *App) updatePageDataFromEvent(eventType string) {
	switch eventType {
	case "pod":
		app.pods.UpdateData()
	case "container":
		app.containers.UpdateData()
	case "network":
		app.networks.UpdateData()
	case "image":
		app.images.UpdateData()
	case "volume":
		app.volumes.UpdateData()
	case "secrets":
		app.secrets.UpdateData()
	}
}

func (app *App) clearViewsData() {
	app.pods.ClearData()
	app.pods.HideAllDialogs()

	app.containers.ClearData()
	app.containers.HideAllDialogs()

	app.volumes.ClearData()
	app.volumes.HideAllDialogs()

	app.networks.ClearData()
	app.networks.HideAllDialogs()

	app.images.ClearData()
	app.images.HideAllDialogs()

	app.secrets.ClearData()
	app.secrets.HideAllDialogs()
}

func (app *App) clearInfoUIData() {
	app.infoBar.UpdateBasicInfo("", "", "")
	app.infoBar.UpdateSystemUsageInfo(0.00, 0.00) //nolint:gomnd
	app.infoBar.UpdatePodmanInfo("", "", "", "")
}
