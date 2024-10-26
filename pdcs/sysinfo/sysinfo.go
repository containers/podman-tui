package sysinfo

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/system"
)

// SystemInfo implements system information data.
type SystemInfo struct {
	Hostname       string
	OS             string
	Arch           string
	Kernel         string
	MemUsagePC     float64
	SwapUsagePC    float64
	Runtime        string
	APIVersion     string
	BuildahVersion string
	ConmonVersion  string
}

// SysInfo returns basic system information.
func SysInfo() (*SystemInfo, error) {
	info := &SystemInfo{}

	conn, err := registry.GetConnection()
	if err != nil {
		return info, err
	}

	response, err := system.Info(conn, nil)
	if err != nil {
		return info, err
	}

	info.Hostname = response.Host.Hostname
	info.OS = response.Host.OS
	info.Arch = response.Host.Arch
	info.Kernel = response.Host.Kernel
	memUsed := response.Host.MemTotal - response.Host.MemFree
	swapUsed := response.Host.SwapTotal - response.Host.SwapFree
	info.MemUsagePC = float64(memUsed*100) / float64(response.Host.MemTotal)    //nolint:mnd
	info.SwapUsagePC = float64(swapUsed*100) / float64(response.Host.SwapTotal) //nolint:mnd

	info.Runtime = response.Host.OCIRuntime.Version
	info.BuildahVersion = response.Host.BuildahVersion
	info.APIVersion = response.Version.APIVersion
	info.ConmonVersion = response.Host.Conmon.Version

	return info, err
}
