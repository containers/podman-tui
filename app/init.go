package app

import "github.com/containers/podman-tui/pdcs/registry"

func (app *App) initUI() {
	connStatus, _ := app.health.ConnStatus()
	if connStatus == registry.ConnectionStatusConnected {
		app.pods.UpdateData()
		app.containers.UpdateData()
		app.networks.UpdateData()
		app.images.UpdateData()
		app.volumes.UpdateData()
		app.initInfoBar()
	}

	app.system.UpdateConnectionsData()
}

func (app *App) initInfoBar() {
	// update basic information
	hostname, kernel, ostype := app.health.GetSysInfo()
	// update memory and swap usage
	memUsage, swapUsage := app.health.GetSysUsage()
	// update podman information
	apiVer, runtime, conmonVer, buildahVer := app.health.GetPodmanInfo()

	connStatus, _ := app.health.ConnStatus()
	if connStatus == registry.ConnectionStatusConnected {
		app.infoBar.UpdateBasicInfo(hostname, kernel, ostype)
		app.infoBar.UpdateSystemUsageInfo(memUsage, swapUsage)
		app.infoBar.UpdatePodmanInfo(apiVer, runtime, conmonVer, buildahVer)

		return
	}

	app.clearInfoUIData()
}
