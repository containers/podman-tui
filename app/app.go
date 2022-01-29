package app

import (
	"time"

	health "github.com/containers/podman-tui/system"
	"github.com/containers/podman-tui/ui/connection"
	"github.com/containers/podman-tui/ui/containers"
	"github.com/containers/podman-tui/ui/images"
	"github.com/containers/podman-tui/ui/infobar"
	"github.com/containers/podman-tui/ui/networks"
	"github.com/containers/podman-tui/ui/pods"
	"github.com/containers/podman-tui/ui/system"
	"github.com/containers/podman-tui/ui/volumes"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	refreshInterval = 1000 * time.Millisecond
)

// App represents main application struct
type App struct {
	*tview.Application
	infoBar         *infobar.InfoBar
	pages           *tview.Pages
	pods            *pods.Pods
	containers      *containers.Containers
	volumes         *volumes.Volumes
	images          *images.Images
	networks        *networks.Networks
	system          *system.System
	connection      *connection.Connection
	menu            *tview.TextView
	health          *health.Engine
	currentPage     string
	needInitUI      bool
	fastRefreshChan chan bool
}

// NewApp returns new app
func NewApp() *App {
	log.Info().Msg("app: new application")
	app := App{
		Application:     tview.NewApplication(),
		pages:           tview.NewPages(),
		needInitUI:      false,
		fastRefreshChan: make(chan bool, 10),
	}
	app.health = health.NewEngine(refreshInterval)

	app.infoBar = infobar.NewInfoBar()

	app.pods = pods.NewPods()
	app.containers = containers.NewContainers()
	app.volumes = volumes.NewVolumes()
	app.images = images.NewImages()
	app.networks = networks.NewNetworks()
	app.system = system.NewSystem()
	app.connection = connection.NewConnection()

	// set refresh channel for container page
	// its requried for container exec dialog
	app.containers.SetFastRefreshChannel(app.fastRefreshChan)

	// menu items
	var menuItems = [][]string{
		{"F1", app.pods.GetTitle()},
		{"F2", app.containers.GetTitle()},
		{"F3", app.volumes.GetTitle()},
		{"F4", app.images.GetTitle()},
		{"F5", app.networks.GetTitle()},
		{"F6", app.system.GetTitle()},
		{"Enter", "commands"},
	}
	app.menu = newMenu(menuItems)
	app.pages.AddPage(app.pods.GetTitle(), app.pods, true, false)
	app.pages.AddPage(app.containers.GetTitle(), app.containers, true, false)
	app.pages.AddPage(app.images.GetTitle(), app.images, true, false)
	app.pages.AddPage(app.volumes.GetTitle(), app.volumes, true, false)
	app.pages.AddPage(app.networks.GetTitle(), app.networks, true, false)
	app.pages.AddPage(app.system.GetTitle(), app.system, true, false)
	app.pages.AddPage(app.connection.GetTitle(), app.connection, true, false)

	return &app
}

// Run starts the application loop.
func (app *App) Run() error {
	log.Info().Msg("app: run")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(app.infoBar, infobar.InfoBarViewHeight, 0, false).
		AddItem(app.pages, 0, 1, false).
		AddItem(app.menu, 1, 1, false)

	// start health check and event parser
	app.health.Start()

	// initial update
	app.initUI()

	// listen for user input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		connOK, _ := app.health.ConnOK()
		if !connOK {
			return event
		}
		switch event.Key() {

		case tcell.KeyF1:
			//pods page
			log.Debug().Msgf("app: switching to %s view", app.pods.GetTitle())
			app.pages.SwitchToPage(app.pods.GetTitle())
			app.SetFocus(app.pods)
			app.pods.UpdateData()
			app.currentPage = app.pods.GetTitle()
			return nil

		case tcell.KeyF2:
			// containers page
			log.Debug().Msgf("app: switching to %s view", app.containers.GetTitle())
			app.pages.SwitchToPage(app.containers.GetTitle())
			app.SetFocus(app.containers)
			app.containers.UpdateData()
			app.currentPage = app.containers.GetTitle()
			return nil

		case tcell.KeyF3:
			// volumes page
			log.Debug().Msgf("app: switching to %s view", app.volumes.GetTitle())
			app.pages.SwitchToPage(app.volumes.GetTitle())
			app.SetFocus(app.volumes)
			app.volumes.UpdateData()
			app.currentPage = app.volumes.GetTitle()
			return nil

		case tcell.KeyF4:
			// images page
			log.Debug().Msgf("app: switching to %s view", app.images.GetTitle())
			app.pages.SwitchToPage(app.images.GetTitle())
			app.SetFocus(app.images)
			app.images.UpdateData()
			app.currentPage = app.images.GetTitle()
			return nil

		case tcell.KeyF5:
			// networks page
			log.Debug().Msgf("app: switching to %s view", app.networks.GetTitle())
			app.pages.SwitchToPage(app.networks.GetTitle())
			app.SetFocus(app.networks)
			app.networks.UpdateData()
			app.currentPage = app.networks.GetTitle()
			return nil
		case tcell.KeyF6:
			// system page
			log.Debug().Msgf("app: switching to %s view", app.system.GetTitle())
			app.pages.SwitchToPage(app.system.GetTitle())
			app.SetFocus(app.system)
			app.currentPage = app.system.GetTitle()
			return nil

		}

		return event
	})
	app.currentPage = app.pods.GetTitle()
	app.pages.SwitchToPage(app.connection.GetTitle())

	// start refresh loop
	go app.refresh()

	// start fast refresh loop
	go app.fastRefresh()

	if err := app.SetRoot(flex, true).SetFocus(app.pods).EnableMouse(false).Run(); err != nil {
		return err
	}
	return nil
}
