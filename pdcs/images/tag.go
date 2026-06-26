package images

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/images"
)

// Tag tags the specified image ID.
func Tag(id string, tag string) error {
	log.Debug().Msgf("pdcs: podman image tag %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	repo := "localhost/" + tag

	return images.Tag(conn, id, tag, repo, new(images.TagOptions))
}

// Untag tags the specified image ID.
func Untag(id string) error {
	log.Debug().Msgf("pdcs: podman image untag %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return images.Untag(conn, id, "", "", new(images.UntagOptions))
}
