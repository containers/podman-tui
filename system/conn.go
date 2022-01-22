package system

import (
	"sync"
	"time"
)

type apiConn struct {
	mu         sync.Mutex
	connOK     bool
	prevStatus bool
	message    string
}

// ConnOK returns connetion status
func (conn *apiConn) ConnOK() (bool, string) {
	status := true
	message := ""
	conn.mu.Lock()
	status = conn.connOK
	message = conn.message
	conn.mu.Unlock()
	return status, message

}

func (conn *apiConn) previousStatus() bool {
	status := false
	conn.mu.Lock()
	status = conn.prevStatus
	conn.mu.Unlock()
	return status
}

func (conn *apiConn) setStatus(status bool, message string) {
	conn.mu.Lock()
	conn.prevStatus = conn.connOK
	conn.connOK = status
	conn.message = message
	conn.mu.Unlock()
}

func (engine *Engine) checkConnHealth() {
	tick := time.NewTicker(engine.refreshInterval)
	for {
		select {
		case <-tick.C:
			engine.updateSysInfo()
		}
	}
}

// ConnOK returns connection status
func (engine *Engine) ConnOK() (bool, string) {
	return engine.conn.ConnOK()
}
