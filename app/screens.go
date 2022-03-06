package app

import "github.com/rs/zerolog/log"

func (app *App) switchToScreen(name string) {
	log.Debug().Msgf("app: switching to %s screen", name)
	app.pages.SwitchToPage(name)
	app.setPageFocus(name)
	app.pods.UpdateData()
	app.currentPage = name
}

func (app *App) frontScreenHasActiveDialog() bool {
	switch app.currentPage {
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
	case app.system.GetTitle():
		return app.system.SubDialogHasFocus()
	}
	return false
}

func (app *App) switchToPreviousScreen() {
	var previousScreen string
	switch app.currentPage {
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
	case app.system.GetTitle():
		previousScreen = app.networks.GetTitle()
	default:
		previousScreen = app.pods.GetTitle()
	}
	app.switchToScreen(previousScreen)
}

func (app *App) switchToNextScreen() {
	var nextScreen string
	switch app.currentPage {
	case app.pods.GetTitle():
		nextScreen = app.containers.GetTitle()
	case app.containers.GetTitle():
		nextScreen = app.volumes.GetTitle()
	case app.volumes.GetTitle():
		nextScreen = app.images.GetTitle()
	case app.images.GetTitle():
		nextScreen = app.networks.GetTitle()
	case app.networks.GetTitle():
		nextScreen = app.system.GetTitle()
	default:
		nextScreen = app.pods.GetTitle()
	}
	app.switchToScreen(nextScreen)
}

func (app *App) setPageFocus(page string) {
	switch page {
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
	case app.system.GetTitle():
		app.Application.SetFocus(app.system)
	case app.connection.GetTitle():
		app.Application.SetFocus(app.connection)
	case app.help.GetTitle():
		app.Application.SetFocus(app.help)
	}
}

func (app *App) updatePageData(eventType string) {
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
	}

}
