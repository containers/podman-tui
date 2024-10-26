package images

import (
	"os"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/channel"
	"github.com/containers/podman/v5/pkg/errorhandling"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ImageSaveOptions image save options.
type ImageSaveOptions struct {
	Output                      string
	Compressed                  bool
	OciAcceptUncompressedLayers bool
	Format                      string
}

// Save saves an image on the local machine.
func Save(imageID string, opts ImageSaveOptions) error { //nolint:cyclop
	log.Debug().Msgf("pdcs: podman image save %v", opts)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	info, err := os.Stat(opts.Output)
	if err == nil {
		if info.Mode().IsRegular() {
			return errors.Errorf("%q already exists as a regular file", opts.Output)
		}
	}

	// performing basic check for output.
	var outputErrors []error

	outputFile, err := os.Create(opts.Output)
	if err != nil {
		log.Info().Msgf("%v", err)

		return err
	}
	defer outputFile.Close()

	cancelChan := make(chan bool, 1)
	writerChan := make(chan []byte, 1024) //nolint:mnd
	outputWriter := channel.NewWriter(writerChan)

	writeOutputFunc := func() {
		log.Debug().Msgf("pdcs: podman image save starting writer")

		for {
			select {
			case <-cancelChan:
				close(writerChan)
				close(cancelChan)
				log.Debug().Msgf("pdcs: podman image save writer stopped")

				return
			case data := <-writerChan:
				if _, err := outputFile.Write(data); err != nil {
					outputErrors = append(outputErrors, err)
				}
			}
		}
	}
	go writeOutputFunc()

	var saveOpts images.ExportOptions

	saveOpts.WithCompress(opts.Compressed)
	saveOpts.WithFormat(opts.Format)
	saveOpts.WithOciAcceptUncompressedLayers(opts.OciAcceptUncompressedLayers)

	err = images.Export(conn, []string{imageID}, outputWriter, &saveOpts)
	cancelChan <- true

	if err != nil {
		return err
	}

	if len(outputErrors) > 0 {
		return errorhandling.JoinErrors(outputErrors)
	}

	return nil
}
