package connection

import (
	"context"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
)

var (
	podmanConnection *context.Context
)

// GetConnection returns connection to podman socket
func GetConnection() (context.Context, error) {
	if podmanConnection == nil {
		// Get Podman socket location
		socket, err := getLocalSocket()
		if err != nil {
			return nil, err
		}
		conn, err := bindings.NewConnection(context.Background(), socket)
		if err != nil {
			return nil, err
		}
		podmanConnection = &conn
		return conn, nil
	}
	return *podmanConnection, nil
}

func getLocalSocket() (string, error) {
	var sockDir string
	var socket string
	currentUser := os.Getenv("USER")
	uid := os.Getenv("UID")

	if currentUser == "root" || uid == "0" {
		sockDir = "/run/"
	} else {
		sockDir = os.Getenv("XDG_RUNTIME_DIR")
	}

	socket = "unix:" + sockDir + "/podman/podman.sock"
	return socket, nil
}
