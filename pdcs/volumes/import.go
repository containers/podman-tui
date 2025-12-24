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

// Import import the specified volume.
func Import(id string, source string) error {
	var (
		conn         context.Context
		err          error
		sourceReader *os.File
	)

	log.Debug().Msgf("pdcs: podman volume import %s %s", id, source)

	conn, err = registry.GetConnection()
	if err != nil {
		return err
	}

	sourceFilePath := filepath.Clean(source)

	sourceReader, err = os.OpenFile(sourceFilePath, os.O_RDONLY, utils.DefaultPermission)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := sourceReader.Close()
		if err == nil {
			err = closeErr
		}
	}()

	err = volumes.Import(conn, id, sourceReader)
	if err != nil {
		return nil
	}

	return err
}
