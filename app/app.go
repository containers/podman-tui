package app

import (
	"time"

	health "github.com/containers/podman-tui/system"
	"github.com/containers/podman-tui/ui/connection"
	"github.com/containers/podman-tui/ui/containers"
	"github.com/containers/podman-tui/ui/help"
	"github.com/containers/podman-tui/ui/images"
	"github.com/containers/podman-tui/ui/infobar"
	"github.com/containers/podman-tui/ui/networks"
	"github.com/containers/podman-tui/ui/pods"
	"github.com/containers/podman-tui/ui/system"
	"github.com/containers/podman-tui/ui/utils"
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
	help            *help.HelpScreen
	currentPage     string
	needInitUI      bool
	fastRefreshChan chan bool
}

// NewApp returns new app
func NewApp(name string, version string) *App {
	log.Debug().Msg("app: new application")
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
	app.help = help.NewHelpScreen(name, version)

	// set refresh channel for container page
	// its required for container exec dialog
	app.containers.SetFastRefreshChannel(app.fastRefreshChan)

	// menu items
	var menuItems = [][]string{
		{utils.HelpScreenKey.Label(), app.help.GetTitle()},
		{utils.PodsScreenKey.Label(), app.pods.GetTitle()},
		{utils.ContainersScreenKey.Label(), app.containers.GetTitle()},
		{utils.VolumesScreenKey.Label(), app.volumes.GetTitle()},
		{utils.ImagesScreenKey.Label(), app.images.GetTitle()},
		{utils.NetworksScreenKey.Label(), app.networks.GetTitle()},
		{utils.SystemScreenKey.Label(), app.system.GetTitle()},
	}
	app.menu = newMenu(menuItems)
	app.pages.AddPage(app.pods.GetTitle(), app.pods, true, false)
	app.pages.AddPage(app.containers.GetTitle(), app.containers, true, false)
	app.pages.AddPage(app.images.GetTitle(), app.images, true, false)
	app.pages.AddPage(app.volumes.GetTitle(), app.volumes, true, false)
	app.pages.AddPage(app.networks.GetTitle(), app.networks, true, false)
	app.pages.AddPage(app.system.GetTitle(), app.system, true, false)
	app.pages.AddPage(app.connection.GetTitle(), app.connection, true, false)
	app.pages.AddPage(app.help.GetTitle(), app.help, true, false)

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

		if !app.frontScreenHasActiveDialog() {
			// previous and next screen keys
			switch event.Rune() {
			case utils.NextScreenKey.Rune():
				// next screen
				app.switchToNextScreen()
				return nil

			case utils.PreviousScreenKey.Rune():
				// previous screen
				app.switchToPreviousScreen()
				return nil
			}

			// normal page key switch
			switch event.Key() {
			case utils.HelpScreenKey.EventKey():
				// help page
				app.switchToScreen(app.help.GetTitle())
				return nil

			case utils.PodsScreenKey.EventKey():
				//pods page
				app.switchToScreen(app.pods.GetTitle())
				return nil

			case utils.ContainersScreenKey.EventKey():
				// containers page
				app.switchToScreen(app.containers.GetTitle())
				return nil

			case utils.VolumesScreenKey.EventKey():
				// volumes page
				app.switchToScreen(app.volumes.GetTitle())
				return nil

			case utils.ImagesScreenKey.EventKey():
				// images page
				app.switchToScreen(app.images.GetTitle())
				return nil

			case utils.NetworksScreenKey.EventKey():
				// networks page
				app.switchToScreen(app.networks.GetTitle())
				return nil

			case utils.SystemScreenKey.EventKey():
				// system page
				app.switchToScreen(app.system.GetTitle())
				return nil
			}
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
