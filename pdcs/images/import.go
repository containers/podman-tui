package images

import (
	"bufio"
	"io"
	"os"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/rs/zerolog/log"
)

// ImageImportOptions image import options.
type ImageImportOptions struct {
	Source    string
	Change    []string
	Message   string
	Reference string
	URL       bool
}

// Import creates a container image from an archive.
func Import(opts ImageImportOptions) (string, error) {
	log.Debug().Msgf("pdcs: podman image import %v", opts)

	conn, err := registry.GetConnection()
	if err != nil {
		return "", err
	}

	var reader io.Reader

	options := new(images.ImportOptions).WithMessage(opts.Message).WithReference(opts.Reference).WithChanges(opts.Change)

	if opts.URL {
		options.WithURL(opts.Source)
	} else {
		tarFile, err := os.Open(opts.Source)
		if err != nil {
			return "", err
		}

		defer tarFile.Close()

		reader = bufio.NewReader(tarFile)
	}

	report, err := images.Import(conn, reader, options)
	if err != nil {
		return "", err
	}

	return report.Id, nil
}
