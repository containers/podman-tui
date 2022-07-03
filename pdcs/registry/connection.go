package registry

import (
	"context"
	"errors"

	"github.com/containers/podman/v4/pkg/bindings"
)

// ErrConnectionNotSelected implements connection is not seleceted error.
var ErrConnectionNotSelected = errors.New("system connection not selected")

// GetConnection returns connection to podman socket.
func GetConnection() (context.Context, error) {
	if !ConnectionIsSet() {
		return nil, ErrConnectionNotSelected
	}

	if pdcsRegistry.connContext == nil {
		var (
			conn context.Context
			err  error
		)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		conn, err = bindings.NewConnectionWithIdentity(ctx, ConnectionURI(), ConnectionIdentity())

		if err != nil {
			cancel()

			return nil, err
		}

		pdcsRegistry.connContext = &conn
		pdcsRegistry.connContextCancel = cancel

		return conn, nil
	}

	return *pdcsRegistry.connContext, nil
}
