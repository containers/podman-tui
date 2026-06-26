package containers

import (
	"strings"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/containers"
)

// CntCommitOptions containers commit options.
type CntCommitOptions struct {
	Author  string
	Changes []string
	Message string
	Format  string
	Pause   bool
	Squash  bool
	Image   string
}

// Commit creates an image from a container's changes.
func Commit(id string, opts CntCommitOptions) (string, error) {
	log.Debug().Msgf("pdcs: podman container commit %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return "", err
	}

	commitOpts := new(containers.CommitOptions)
	commitOpts.WithAuthor(opts.Author)
	commitOpts.WithChanges(opts.Changes)
	commitOpts.WithComment(opts.Message)
	commitOpts.WithPause(opts.Pause)
	commitOpts.WithSquash(opts.Squash)
	imageRepoTag := strings.Split(opts.Image, ":")
	commitOpts.WithRepo(imageRepoTag[0])

	if len(imageRepoTag) > 1 {
		commitOpts.WithTag(imageRepoTag[1])
	}

	response, err := containers.Commit(conn, id, commitOpts)

	return response.ID, err
}
