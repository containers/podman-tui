package images

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

// ImageBuildOptions image build options.
type ImageBuildOptions struct {
	ContainerFiles []string
	BuildOptions   entities.BuildOptions
}

// Build creates an image using a containerfile reference.
func Build(opts ImageBuildOptions) (string, error) {
	log.Debug().Msgf("pdcs: podman image build %v", opts)

	conn, err := registry.GetConnection()
	if err != nil {
		return "", err
	}

	response, err := images.Build(conn, opts.ContainerFiles, opts.BuildOptions)
	if err != nil {
		return "", err
	}

	return response.ID, nil
}
