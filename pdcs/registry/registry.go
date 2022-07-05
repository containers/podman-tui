package registry

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

// Connection status.
const (
	ConnectionStatusDisconnected = 0 + iota
	ConnectionStatusConnected
	ConnectionStatusConnectionError
)

var pdcsRegistry registry

// service implements podman connection service.
type registry struct {
	mu                sync.Mutex
	connection        Connection
	connContext       *context.Context
	connContextCancel func()
	connectionIsSet   bool
}

// Connection implements a system connection.
type Connection struct {
	Name     string
	Default  bool
	Status   ConnStatus
	URI      string
	Identity string
}

// ConnStatus implements Connection status.
type ConnStatus int

func init() {
	pdcsRegistry.connectionIsSet = false
}

// SetConnectionStatus sets registry Connection status.
func SetConnectionStatus(status ConnStatus) {
	if !ConnectionIsSet() {
		return
	}

	pdcsRegistry.mu.Lock()
	defer pdcsRegistry.mu.Unlock()
	pdcsRegistry.connection.Status = status
}

// SetConnection sets registry connection.
func SetConnection(connection Connection) {
	log.Debug().Msgf("pdcs: registry set connection %v", connection)
	pdcsRegistry.mu.Lock()
	defer pdcsRegistry.mu.Unlock()
	pdcsRegistry.connection = connection
	pdcsRegistry.connectionIsSet = true
}

// UnsetConnection unsets the registry loaded connection.
func UnsetConnection() {
	log.Debug().Msgf("pdcs: registry unset connection")
	pdcsRegistry.mu.Lock()
	pdcsRegistry.connectionIsSet = false
	pdcsRegistry.connection = Connection{}
	pdcsRegistry.mu.Unlock()
	CancelContext()
}

// CancelContext run the cancel function for context.
func CancelContext() {
	pdcsRegistry.mu.Lock()
	defer pdcsRegistry.mu.Unlock()

	if pdcsRegistry.connContextCancel != nil {
		log.Debug().Msgf("pdcs: registry context cancel")
		pdcsRegistry.connContextCancel()
	}

	pdcsRegistry.connContext = nil
}

// ConnectionIsSet returns true if connection is set.
func ConnectionIsSet() bool {
	pdcsRegistry.mu.Lock()
	defer pdcsRegistry.mu.Unlock()

	return pdcsRegistry.connectionIsSet
}

// ConnectionName returns selected connection name.
func ConnectionName() string {
	pdcsRegistry.mu.Lock()
	defer pdcsRegistry.mu.Unlock()

	return pdcsRegistry.connection.Name
}

// ConnectionStatus returns selected connection status.
func ConnectionStatus() ConnStatus {
	pdcsRegistry.mu.Lock()
	defer pdcsRegistry.mu.Unlock()

	return pdcsRegistry.connection.Status
}

// ConnectionURI returns selected connection url.
func ConnectionURI() string {
	pdcsRegistry.mu.Lock()
	defer pdcsRegistry.mu.Unlock()

	return pdcsRegistry.connection.URI
}

// ConnectionIdentity returns selected connection identity.
func ConnectionIdentity() string {
	pdcsRegistry.mu.Lock()
	defer pdcsRegistry.mu.Unlock()

	return pdcsRegistry.connection.Identity
}

func (connStatus ConnStatus) String() string {
	var status string

	switch connStatus {
	case ConnectionStatusConnected:
		status = "connected"
	case ConnectionStatusDisconnected:
		status = "disconnected"
	case ConnectionStatusConnectionError:
		status = "connection error"
	}

	return status
}
