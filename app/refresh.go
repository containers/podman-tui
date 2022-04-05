package app

import (
	"time"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rs/zerolog/log"
)

func (app *App) refresh() {
	log.Debug().Msgf("app: starting refresh loop (interval=%v)", utils.RefreshInterval)
	tick := time.NewTicker(utils.RefreshInterval)
	for {
		select {
		case <-tick.C:
			connStatus, connMsg := app.health.ConnStatus()
			switch connStatus {
			case registry.ConnectionStatusConnected:
				app.refreshConnOK()
			case registry.ConnectionStatusConnectionError:
				app.refreshNotConnOK()
				if registry.ConnectionIsSet() {
					name := registry.ConnectionName()
					app.system.ConnectionProgressDisplay(true)
					app.system.SetConnectionProgressDestName(name)
					app.system.SetConnectionProgressMessage(connMsg)
				}
			case registry.ConnectionStatusDisconnected:
				app.refreshNotConnOK()
				if registry.ConnectionIsSet() {
					name := registry.ConnectionName()
					app.system.ConnectionProgressDisplay(true)
					app.system.SetConnectionProgressDestName(name)
				}
			}
			app.initInfoBar()
			app.infoBar.UpdateConnStatus(connStatus)
			app.Application.Draw()
		}
	}
}

func (app *App) refreshConnOK() {
	if app.needInitUI {
		// init ui after reconnection
		app.initUI()
		app.system.ConnectionProgressDisplay(false)
		app.needInitUI = false
		app.pages.SwitchToPage(app.currentPage)
		app.setPageFocus(app.currentPage)
	}
	app.flushEvents()
}

func (app *App) refreshNotConnOK() {
	// only switch to system view one time
	if !app.needInitUI {
		app.clearViewsData()
		app.switchToScreen(app.system.GetTitle())
	}
	app.system.UpdateConnectionsData()
	app.needInitUI = true
}

func (app *App) flushEvents() {
	// update events
	eventTypes := app.health.GetEvents()
	for _, evt := range eventTypes {
		app.updatePageDataFromEvent(evt)
	}
	if app.health.HasNewEvent() {
		app.system.SetEventMessage(app.health.GetEventMessages())
	}
}

// fastRefresh method will refresh the screen as soon as it receives
// the refresh signal. Its required for some feature e.g. container exec
func (app *App) fastRefresh() {
	log.Debug().Msg("app: starting fast refresh loop")
	for {
		select {
		case refresh := <-app.fastRefreshChan:
			{
				if refresh {
					app.Application.Draw()
				}
			}
		}
	}
}
