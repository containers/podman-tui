package images

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/images"
)

// Tree returns a tree based representation of the given image.
func Tree(id string) (string, error) {
	log.Debug().Msgf("pdcs: podman image tree %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return "", err
	}

	tree, err := images.Tree(conn, id, new(images.TreeOptions).WithWhatRequires(true))

	return tree.Tree, err
}
