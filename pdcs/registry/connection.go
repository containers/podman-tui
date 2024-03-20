package registry

import (
	"context"
	"errors"
	"net/url"
	"os"

	"github.com/containers/podman/v5/pkg/bindings"
)

// ErrConnectionNotSelected implements connection is not selected error.
var ErrConnectionNotSelected = errors.New("system connection not selected")

// GetConnection returns connection to podman socket.
func GetConnection() (context.Context, error) {
	if !ConnectionIsSet() {
		return nil, ErrConnectionNotSelected
	}

	if pdcsRegistry.connContext == nil {
		var (
			conn       context.Context
			err        error
			passPhrase string
		)

		dest := ConnectionURI()

		connURI, err := url.Parse(dest)
		if err != nil {
			return nil, err
		}

		if v, found := os.LookupEnv("CONTAINER_PASSPHRASE"); found {
			passPhrase = v
		}

		connURI.User = url.UserPassword(connURI.User.String(), passPhrase)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		conn, err = bindings.NewConnectionWithIdentity(ctx, connURI.String(), ConnectionIdentity(), false)
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
