package app

import (
	"time"

	"github.com/rs/zerolog/log"
)

func (app *App) refresh() {
	log.Debug().Msg("app: starting refresh loop")
	tick := time.NewTicker(refreshInterval)
	for {
		select {
		case <-tick.C:
			connOK, connMsg := app.health.ConnOK()
			if connOK {
				if app.needInitUI {
					// init ui after reconnection
					app.initUI()
					app.needInitUI = false
					app.pages.SwitchToPage(app.currentPage)
					app.setFocus(app.currentPage)

				}
				eventTypes := app.health.GetEvents()
				// update events
				for _, evt := range eventTypes {
					app.updatePageData(evt)
				}
				if app.health.HasNewEvent() {
					app.system.SetEventMessage(app.health.GetEventMessages())
				}
			} else {
				// set init ui to true
				app.needInitUI = true
				app.clearUIData()
				app.pages.SwitchToPage(app.connection.GetTitle())
				app.connection.SetErrorMessage(connMsg)
				app.setFocus(app.connection.GetTitle())
			}
			app.initInfoBar()
			app.infoBar.UpdateConnStatus(connOK)
			app.Application.Draw()
		}
	}
}

func (app *App) setFocus(page string) {
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
