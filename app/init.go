package app

func (app *App) initUI() {
	app.connection.Reset()
	connOK, _ := app.health.ConnOK()
	if connOK {
		app.pods.UpdateData()
		app.containers.UpdateData()
		app.networks.UpdateData()
		app.images.UpdateData()
		app.volumes.UpdateData()
		app.initInfoBar()
	}
}

func (app *App) initInfoBar() {
	// update basic information
	hostname, kernel, ostype := app.health.GetSysInfo()
	app.infoBar.UpdateBasicInfo(hostname, kernel, ostype)

	// udpate memory and swap usage
	memUsage, swapUsage := app.health.GetSysUsage()
	app.infoBar.UpdateSystemUsageInfo(memUsage, swapUsage)

	// update podman information
	apiVer, runtime, conmonVer, buildahVer := app.health.GetPodmanInfo()
	app.infoBar.UpdatePodmanInfo(apiVer, runtime, conmonVer, buildahVer)

}

func (app *App) clearUIData() {
	app.pods.ClearData()
	app.containers.ClearData()
	app.volumes.ClearData()
	app.networks.ClearData()
	app.images.ClearData()
}
