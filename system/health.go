package system

import (
	"fmt"
	"time"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/rs/zerolog/log"
)

const (
	messageBufferSize = 100
)

// Engine implements connection and system info check
type Engine struct {
	refreshInterval time.Duration
	sysinfo         systemInfo
	sysEvents       podmanEvents
	conn            apiConn
}

// NewEngine returns new health checker
func NewEngine(refreshInterval time.Duration) *Engine {
	health := &Engine{
		conn: apiConn{
			connStaus:  registry.ConnectionStatusDisconnected,
			prevStatus: registry.ConnectionStatusDisconnected,
		},
		refreshInterval: refreshInterval,
		sysEvents: podmanEvents{
			messageBufferSize: messageBufferSize,
		},
		sysinfo: systemInfo{},
	}
	//health.sysEvents.eventCancelChan = make(chan bool)
	//health.sysEvents.cancelChan = make(chan bool)
	//health.sysEvents.eventChan = make(chan entities.Event, 20)
	//health.updateSysInfo()
	return health
}

// Start starts health checkers
func (engine *Engine) Start() {
	engine.sysinfo.info = &sysinfo.SystemInfo{}
	go engine.healthCheckLoop()
}

// ConnStatus returns connection status
func (engine *Engine) ConnStatus() (registry.ConnStatus, string) {
	return engine.conn.ConnStatus()
}

// Connect sets engine connection
func (engine *Engine) Connect(connection registry.Connection) {
	log.Debug().Msgf("health: connect to %v", connection)
	if registry.ConnectionIsSet() {
		engine.Disconnect()
	}
	registry.SetConnection(connection)
}

// Disconnect disconnects engine and unsets the connection
func (engine *Engine) Disconnect() {
	if !registry.ConnectionIsSet() {
		return
	}
	log.Debug().Msgf("health: disconnect")
	registry.SetConnectionStatus(registry.ConnectionStatusDisconnected)
	engine.conn.setStatus(registry.ConnectionStatusDisconnected, "")
	registry.UnsetConnection()
}

func (engine *Engine) healthCheckLoop() {
	tick := time.NewTicker(engine.refreshInterval)
	for {
		select {
		case <-tick.C:
			engine.healthCheck()
		}
	}
}

func (engine *Engine) healthCheck() {
	info, err := sysinfo.SysInfo()
	status := true
	if err != nil {
		status = false
		if err == registry.ErrConnectionNotSelected {
			engine.conn.setStatus(registry.ConnectionStatusDisconnected, "")
		} else {
			engine.conn.setStatus(registry.ConnectionStatusConnectionError, fmt.Sprintf("%v", err))
			log.Error().Msgf("health check: %v", err)
		}
		if registry.ConnectionIsSet() {
			registry.SetConnectionStatus(registry.ConnectionStatusConnectionError)
		} else {
			registry.SetConnectionStatus(registry.ConnectionStatusDisconnected)
		}

		if engine.conn.previousStatus() == registry.ConnectionStatusDisconnected {
			engine.clearSysInfoData()
		}
		return
	}
	// starting event streaming process after reconnecting
	if status && !engine.EventStatus() {
		if engine.conn.previousStatus() == registry.ConnectionStatusConnected {
			engine.startEventStreamer()
		}
	}
	engine.conn.setStatus(registry.ConnectionStatusConnected, "")
	engine.sysinfo.mu.Lock()
	engine.sysinfo.info = info
	engine.sysinfo.mu.Unlock()
	registry.SetConnectionStatus(registry.ConnectionStatusConnected)
}
