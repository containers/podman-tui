package images

import (
	"fmt"
	"sort"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/distribution/reference"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	noneTag string = "<none>"
)

// List returns list of images information.
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

// ImageListReporter image list report.
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

		if len(e.RepoTags) > 0 { //nolint:nestif
			tagged := []ImageListReporter{}
			untagged := []ImageListReporter{}

			for _, tag := range e.RepoTags {
				h.ImageSummary = *e

				h.Repository, h.Tag, err = tokenRepoTag(tag)
				if err != nil {
					return nil, errors.Wrapf(err, "error parsing repository tag %q", tag)
				}

				if h.Tag == noneTag {
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
			h.Repository = noneTag
			h.Tag = noneTag
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
	tagRef := fmt.Sprintf("%s:%s", noneTag, noneTag)
	if ref == tagRef {
		return noneTag, noneTag, nil
	}

	repo, err := reference.Parse(ref)
	if err != nil {
		return noneTag, noneTag, err
	}

	named, ok := repo.(reference.Named)
	if !ok {
		return ref, noneTag, nil
	}

	name := named.Name()
	if name == "" {
		name = noneTag
	}

	tagged, ok := repo.(reference.Tagged)
	if !ok {
		return name, noneTag, nil
	}

	tag := tagged.Tag()
	if tag == "" {
		tag = noneTag
	}

	return name, tag, nil
}
