package app

import (
	"time"

	"github.com/rs/zerolog/log"
)

func (app *App) refresh() {
	log.Debug().Msgf("app: starting refresh loop (interval=%v)", refreshInterval)
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
					app.setPageFocus(app.currentPage)

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
				app.setPageFocus(app.connection.GetTitle())
			}
			app.initInfoBar()
			app.infoBar.UpdateConnStatus(connOK)
			app.Application.Draw()
		}
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
