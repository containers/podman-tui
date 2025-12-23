package volumes

import (
	"context"
	"os"
	"path/filepath"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/rs/zerolog/log"
)

// Export export the specified volume.
func Export(id string, output string) error {
	var (
		conn         context.Context
		err          error
		outputWriter *os.File
	)

	log.Debug().Msgf("pdcs: podman volume export %s %s", id, output)

	conn, err = registry.GetConnection()
	if err != nil {
		return err
	}

	outputFile := filepath.Clean(output)

	outputWriter, err = os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, utils.DefaultPermission)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := outputWriter.Close()
		if err == nil {
			err = closeErr
		}
	}()

	err = volumes.Export(conn, id, outputWriter)
	if err != nil {
		return nil
	}

	return err
}
