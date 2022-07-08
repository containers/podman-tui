package sysinfo

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/system"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/docker/go-units"
	"github.com/rs/zerolog/log"
)

// DfSummary implements df summary report.
type DfSummary struct {
	rType       string
	total       int
	active      int
	size        int64
	reclaimable int64
}

// DiskUsage returns information about image, container, and volume disk
// consumption.
func DiskUsage() ([]*DfSummary, error) {
	log.Debug().Msgf("pdcs: podman system disk usage")

	conn, err := registry.GetConnection()
	if err != nil {
		return nil, err
	}

	dfRawReport, err := system.DiskUsage(conn, new(system.DiskOptions))
	if err != nil {
		return nil, err
	}

	dfreport := prepDfSummary(dfRawReport)

	return dfreport, nil
}

func prepDfSummary(reports *entities.SystemDfReport) []*DfSummary { // nolint:funlen
	var (
		dfSummaries       []*DfSummary
		active            int
		size, reclaimable int64
	)

	// Images
	for _, i := range reports.Images {
		if i.Containers > 0 {
			active++
		}

		size += i.Size

		if i.Containers < 1 {
			reclaimable += i.Size
		}
	}

	imageSummary := DfSummary{
		rType:       "Images",
		total:       len(reports.Images),
		active:      active,
		size:        size,
		reclaimable: reclaimable,
	}
	dfSummaries = append(dfSummaries, &imageSummary)

	// Containers
	var (
		conActive               int
		conSize, conReclaimable int64
	)

	for _, c := range reports.Containers {
		if c.Status == "running" {
			conActive++
		} else {
			conReclaimable += c.RWSize
		}

		conSize += c.RWSize
	}

	containerSummary := DfSummary{
		rType:       "Containers",
		total:       len(reports.Containers),
		active:      conActive,
		size:        conSize,
		reclaimable: conReclaimable,
	}
	dfSummaries = append(dfSummaries, &containerSummary)

	// Volumes
	var (
		activeVolumes                   int
		volumesSize, volumesReclaimable int64
	)

	for _, v := range reports.Volumes {
		activeVolumes += v.Links
		volumesSize += v.Size
		volumesReclaimable += v.ReclaimableSize
	}

	volumeSummary := DfSummary{
		rType:       "Local Volumes",
		total:       len(reports.Volumes),
		active:      activeVolumes,
		size:        volumesSize,
		reclaimable: volumesReclaimable,
	}
	dfSummaries = append(dfSummaries, &volumeSummary)

	return dfSummaries
}

// Type returns df summary report type: Images, Containers or Local Volumes.
func (dfsum *DfSummary) Type() string {
	return dfsum.rType
}

// Total returns total value of df summary.
func (dfsum *DfSummary) Total() string {
	return fmt.Sprintf("%d", dfsum.total)
}

// Active returns active value of df summary.
func (dfsum *DfSummary) Active() string {
	return fmt.Sprintf("%d", dfsum.active)
}

// Size returns size value of df summary.
func (dfsum *DfSummary) Size() string {
	return units.HumanSize(float64(dfsum.size))
}

// Reclaimable returns reclaimable value of df summary.
func (dfsum *DfSummary) Reclaimable() string {
	percent := int(float64(dfsum.reclaimable)/float64(dfsum.size)) * 100 // nolint:gomnd

	return fmt.Sprintf("%s (%d%%)", units.HumanSize(float64(dfsum.reclaimable)), percent)
}
