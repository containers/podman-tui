package images

import (
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Tag tags the specified image ID
func Tag(id string, tag string) error {
	log.Debug().Msgf("pdcs: podman image tag %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	repo := "localhost/" + tag
	return images.Tag(conn, id, tag, repo, new(images.TagOptions))

}

// Untag tags the specified image ID
func Untag(id string) error {
	log.Debug().Msgf("pdcs: podman image untag %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	return images.Untag(conn, id, "", "", new(images.UntagOptions))
}
