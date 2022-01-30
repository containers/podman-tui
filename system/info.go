package system

import (
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/rs/zerolog/log"
)

type systemInfo struct {
	mu   sync.Mutex
	info *sysinfo.SystemInfo
}

// GetSysInfo returns podman system information: hostname, kernel, ostype
func (engine *Engine) GetSysInfo() (string, string, string) {
	var hostname string
	var kernel string
	var ostype string
	engine.sysinfo.mu.Lock()
	hostname = engine.sysinfo.info.Hostname
	kernel = engine.sysinfo.info.Kernel
	ostype = engine.sysinfo.info.OS
	engine.sysinfo.mu.Unlock()
	return hostname, kernel, ostype
}

// GetSysUsage returns podman system memory and swap usage
func (engine *Engine) GetSysUsage() (float64, float64) {
	var memUsage float64
	var swapUsage float64
	engine.sysinfo.mu.Lock()
	memUsage = engine.sysinfo.info.MemUsagePC
	swapUsage = engine.sysinfo.info.SwapUsagePC
	engine.sysinfo.mu.Unlock()
	return memUsage, swapUsage
}

// GetPodmanInfo returns podman information: api, runtime, conmon and buildah
func (engine *Engine) GetPodmanInfo() (string, string, string, string) {
	var apiVersion string
	var conmonVersion string
	var buildahVersion string
	var runtime string
	engine.sysinfo.mu.Lock()
	apiVersion = engine.sysinfo.info.APIVersion
	conmonVersion = engine.sysinfo.info.ConmonVersion
	buildahVersion = engine.sysinfo.info.BuildahVersion
	runtime = engine.sysinfo.info.Runtime
	engine.sysinfo.mu.Unlock()
	// conmon version
	conmonVersion = strings.Split(conmonVersion, ",")[0]
	conmonVersion = strings.ReplaceAll(conmonVersion, "conmon version", "")
	conmonVersion = strings.TrimSpace(conmonVersion)

	// runtime
	runtime = strings.Split(runtime, ":")[0]
	runtime = strings.ReplaceAll(runtime, "commit", "")

	return apiVersion, runtime, conmonVersion, buildahVersion
}

func (engine *Engine) updateSysInfo() {
	info, err := sysinfo.SysInfo()
	status := true
	if err != nil {
		log.Error().Msgf("health check: %v", err)
		status = false
		engine.conn.setStatus(false, fmt.Sprintf("%v", err))
		if status != engine.conn.previousStatus() {
			engine.sysEvents.cancelChan <- true
			engine.clearSysInfoData()
		}
		return

	}
	// starting event streaming process after reconnecting
	if !engine.conn.previousStatus() && status && !engine.sysEvents.status {
		engine.startEventStreamer()
	}
	engine.conn.setStatus(true, "")
	engine.sysinfo.mu.Lock()
	engine.sysinfo.info = info
	engine.sysinfo.mu.Unlock()
}

func (engine *Engine) clearSysInfoData() {
	engine.sysinfo.mu.Lock()
	engine.sysinfo.info.Hostname = ""
	engine.sysinfo.info.Kernel = ""
	engine.sysinfo.info.OS = ""
	engine.sysinfo.info.MemUsagePC = 0.00
	engine.sysinfo.info.SwapUsagePC = 0.00
	engine.sysinfo.info.APIVersion = ""
	engine.sysinfo.info.ConmonVersion = ""
	engine.sysinfo.info.BuildahVersion = ""
	engine.sysinfo.info.Runtime = ""
	engine.sysinfo.mu.Unlock()
}
