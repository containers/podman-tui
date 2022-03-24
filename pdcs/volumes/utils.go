package volumes

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/volumes"
	"github.com/rs/zerolog/log"
)

// VolumeDest returns a volume destination (mountpoint) path
func VolumeDest(name string) (string, error) {
	log.Debug().Msgf("pdcs: podman volume destination query: %s", name)
	var dest string
	conn, err := registry.GetConnection()
	if err != nil {
		return dest, err
	}
	queryFilter := make(map[string][]string)
	queryFilter["name"] = []string{
		name,
	}
	response, err := volumes.List(conn, new(volumes.ListOptions).WithFilters(queryFilter))
	if err != nil {
		return dest, err
	}
	for _, resp := range response {
		if resp.Name == name {
			dest = resp.Mountpoint
		}
	}
	if dest == "" {
		return "", fmt.Errorf("volume destination cannot be empty")
	}
	return dest, nil
}
