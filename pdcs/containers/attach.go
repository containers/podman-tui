package containers

import (
	"errors"
	"io"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

var errAttachCntRunning = errors.New("you can only attach to running container")

// Attach will attach to a running container.
func Attach(id string, stdin io.Reader, stdout io.Writer, attachReady chan bool, detachKey string) error {
	log.Debug().Msgf("pdcs: podman container attach %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	cntInSpecData, err := containers.Inspect(conn, id, new(containers.InspectOptions))
	if err != nil {
		return err
	}

	if !cntInSpecData.State.Running {
		return errAttachCntRunning
	}

	attachOptions := new(containers.AttachOptions)
	attachOptions.WithDetachKeys(detachKey)

	return containers.Attach(conn, id, stdin, stdout, stdout, attachReady, attachOptions)
}

// ResizeContainerTTY resize the attached container tty size.
func ResizeContainerTTY(id string, width int, height int) error {
	log.Debug().Msgf("pdcs: podman container %s attach resize tty width=%d,height=%d", id, width, height)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	resizeOptions := new(containers.ResizeTTYOptions)
	resizeOptions.WithWidth(height)
	resizeOptions.WithHeight(height)
	resizeOptions.WithRunning(true)

	return containers.ResizeContainerTTY(conn, id, resizeOptions)
}
