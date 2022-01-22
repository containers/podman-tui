package system

import (
	"time"

	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/rs/zerolog/log"
)

const (
	messageBufferSize = 100
)

// Engine implements connetion and system info check
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
			connOK:     false,
			prevStatus: false,
		},
		refreshInterval: refreshInterval,
		sysEvents: podmanEvents{
			messageBufferSize: messageBufferSize,
		},
		sysinfo: systemInfo{},
	}
	//health.updateSysInfo()
	return health
}

// Start starts health checkers
func (engine *Engine) Start() {
	// check init connection
	var err error
	engine.sysinfo.info, err = sysinfo.SysInfo()
	if err != nil {
		log.Error().Msgf("health check: initial connection status_nok: %s", err.Error())
		engine.conn.setStatus(false, err.Error())
	} else {
		engine.startEventStreamer()
	}

	go engine.checkConnHealth()
}
