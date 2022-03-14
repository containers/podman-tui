package app

import (
	"os"
	"time"

	"github.com/containers/podman-tui/config"
	"github.com/containers/podman-tui/pdcs/registry"
	health "github.com/containers/podman-tui/system"
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
	menu            *tview.TextView
	health          *health.Engine
	help            *help.Help
	currentPage     string
	needInitUI      bool
	fastRefreshChan chan bool
	config          *config.Config
}

// NewApp returns new app
func NewApp(name string, version string) *App {
	log.Debug().Msg("app: new application")

	// create application UI
	app := App{
		Application:     tview.NewApplication(),
		pages:           tview.NewPages(),
		needInitUI:      false,
		fastRefreshChan: make(chan bool, 10),
	}

	var err error
	app.config, err = config.NewConfig()
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}

	app.health = health.NewEngine(refreshInterval)

	app.infoBar = infobar.NewInfoBar()

	app.pods = pods.NewPods()
	app.containers = containers.NewContainers()
	app.volumes = volumes.NewVolumes()
	app.images = images.NewImages()
	app.networks = networks.NewNetworks()
	app.system = system.NewSystem()
	app.system.SetConnectionListFunc(app.config.ServicesConnections)
	app.system.SetConnectionSetDefaultFunc(func(name string) error {
		err := app.config.SetDefaultService(name)
		app.system.UpdateConnectionsData()
		return err
	})
	app.system.SetConnectionConnectFunc(app.health.Connect)
	app.system.SetConnectionDisconnectFunc(app.health.Disconnect)
	app.system.SetConnectionAddFunc(app.config.Add)
	app.system.SetConnectionRemoveFunc(app.config.Remove)

	app.help = help.NewHelp(name, version)

	// set refresh channel for container page
	// its required for container exec dialog
	app.containers.SetFastRefreshChannel(app.fastRefreshChan)

	// menu items
	var menuItems = [][]string{
		{utils.HelpScreenKey.Label(), app.help.GetTitle()},
		{utils.SystemScreenKey.Label(), app.system.GetTitle()},
		{utils.PodsScreenKey.Label(), app.pods.GetTitle()},
		{utils.ContainersScreenKey.Label(), app.containers.GetTitle()},
		{utils.VolumesScreenKey.Label(), app.volumes.GetTitle()},
		{utils.ImagesScreenKey.Label(), app.images.GetTitle()},
		{utils.NetworksScreenKey.Label(), app.networks.GetTitle()},
	}
	app.menu = newMenu(menuItems)
	app.pages.AddPage(app.help.GetTitle(), app.help, true, false)
	app.pages.AddPage(app.system.GetTitle(), app.system, true, false)
	app.pages.AddPage(app.pods.GetTitle(), app.pods, true, false)
	app.pages.AddPage(app.containers.GetTitle(), app.containers, true, false)
	app.pages.AddPage(app.images.GetTitle(), app.images, true, false)
	app.pages.AddPage(app.volumes.GetTitle(), app.volumes, true, false)
	app.pages.AddPage(app.networks.GetTitle(), app.networks, true, false)

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
		if event.Key() == utils.AppExitKey.Key {
			log.Info().Msg("app: stop")
			app.Stop()
			os.Exit(0)
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

			case utils.SystemScreenKey.EventKey():
				// system page
				app.switchToScreen(app.system.GetTitle())
				return nil

			case utils.PodsScreenKey.EventKey():
				// pods page
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
			}
		}

		// if connection is not OK do not display command dialog
		// except for system view and app exit key
		if app.currentPage != app.system.GetTitle() {
			connStatus, _ := app.health.ConnStatus()
			if connStatus != registry.ConnectionStatusConnected {
				return nil
			}
		}
		return event
	})
	app.currentPage = app.system.GetTitle()
	app.pages.SwitchToPage(app.system.GetTitle())

	// start refresh loop
	go app.refresh()

	// start fast refresh loop
	go app.fastRefresh()

	if err := app.SetRoot(flex, true).SetFocus(app.system).EnableMouse(false).Run(); err != nil {
		return err
	}
	return nil
}
