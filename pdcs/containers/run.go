package containers

import (
	"io"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// RunInitAttach will init container for run and attach.
func RunInitAttach(cntID string, stdin io.Reader, stdout io.Writer, attachReady chan bool, detachKey string) error {
	log.Debug().Msgf("pdcs: podman container run init attach %s", cntID)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	attachOptions := new(containers.AttachOptions)
	attachOptions.WithDetachKeys(detachKey)

	err = containers.ContainerInit(conn, cntID, new(containers.InitOptions))
	if err != nil {
		return err
	}

	err = containers.Attach(conn, cntID, stdin, stdout, stdout, attachReady, attachOptions)
	if err != nil {
		return err
	}

	return nil
}
