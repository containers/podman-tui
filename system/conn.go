package system

import (
	"sync"

	"github.com/containers/podman-tui/pdcs/registry"
)

type apiConn struct {
	mu         sync.Mutex
	connStaus  registry.ConnStatus
	prevStatus registry.ConnStatus
	message    string
}

// ConnStatus returns connection status
func (conn *apiConn) ConnStatus() (registry.ConnStatus, string) {
	var status registry.ConnStatus
	message := ""
	conn.mu.Lock()
	status = conn.connStaus
	message = conn.message
	conn.mu.Unlock()
	return status, message

}

func (conn *apiConn) previousStatus() registry.ConnStatus {
	var status registry.ConnStatus
	conn.mu.Lock()
	status = conn.prevStatus
	conn.mu.Unlock()
	return status
}

func (conn *apiConn) setStatus(status registry.ConnStatus, message string) {
	conn.mu.Lock()
	conn.prevStatus = conn.connStaus
	conn.connStaus = status
	conn.message = message
	conn.mu.Unlock()
}
