package images

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/rs/zerolog/log"
)

// ImagePushOptions image push options
type ImagePushOptions struct {
	Desitnation   string
	Compress      bool
	Format        string
	SkipTLSVerify bool
	AuthFile      string
	Username      string
	Password      string
}

// Push push a source image to a specified destination
func Push(id string, opts ImagePushOptions) error {
	log.Debug().Msgf("pdcs: podman image push %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}
	pushOptions := new(images.PushOptions)
	pushOptions.WithCompress(opts.Compress)
	pushOptions.WithFormat(opts.Format)
	pushOptions.WithSkipTLSVerify(opts.SkipTLSVerify)
	pushOptions.WithAuthfile(opts.AuthFile)
	pushOptions.WithUsername(opts.Username)
	pushOptions.WithPassword(opts.Password)

	return images.Push(conn, id, opts.Desitnation, pushOptions)
}
