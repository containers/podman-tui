package images

import (
	"sort"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// List returns list of images information
func List() ([]ImageListReporter, error) {
	log.Debug().Msg("pdcs: podman image ls")

	conn, err := registry.GetConnection()
	if err != nil {
		return nil, err
	}
	response, err := images.List(conn, new(images.ListOptions).WithAll(true))
	if err != nil {
		return nil, err
	}
	imgs, err := sortImages(response)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("pdcs: %v", imgs)
	return imgs, nil
}

//ImageListReporter image list report
type ImageListReporter struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
	entities.ImageSummary
}

func sortImages(imageS []*entities.ImageSummary) ([]ImageListReporter, error) {
	imgs := make([]ImageListReporter, 0, len(imageS))
	var err error
	for _, e := range imageS {
		var h ImageListReporter
		if len(e.RepoTags) > 0 {
			tagged := []ImageListReporter{}
			untagged := []ImageListReporter{}
			for _, tag := range e.RepoTags {
				h.ImageSummary = *e
				h.Repository, h.Tag, err = tokenRepoTag(tag)
				if err != nil {
					return nil, errors.Wrapf(err, "error parsing repository tag %q:", tag)
				}
				if h.Tag == "<none>" {
					untagged = append(untagged, h)
				} else {
					tagged = append(tagged, h)
				}
			}
			// Note: we only want to display "<none>" if we
			// couldn't find any tagged name in RepoTags.
			if len(tagged) > 0 {
				imgs = append(imgs, tagged...)
			} else {
				imgs = append(imgs, untagged[0])
			}
		} else {
			h.ImageSummary = *e
			h.Repository = "<none>"
			h.Tag = "<none>"
			imgs = append(imgs, h)
		}
	}

	sort.Slice(imgs, sortFunc(imgs))
	return imgs, err
}

func sortFunc(data []ImageListReporter) func(i, j int) bool {
	return func(i, j int) bool {
		return data[i].Repository < data[j].Repository
	}
}

func tokenRepoTag(ref string) (string, string, error) {
	if ref == "<none>:<none>" {
		return "<none>", "<none>", nil
	}

	repo, err := reference.Parse(ref)
	if err != nil {
		return "<none>", "<none>", err
	}

	named, ok := repo.(reference.Named)
	if !ok {
		return ref, "<none>", nil
	}
	name := named.Name()
	if name == "" {
		name = "<none>"
	}

	tagged, ok := repo.(reference.Tagged)
	if !ok {
		return name, "<none>", nil
	}
	tag := tagged.Tag()
	if tag == "" {
		tag = "<none>"
	}

	return name, tag, nil
}
