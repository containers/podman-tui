package pods

import (
	"sort"
	"strconv"
	"strings"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/rs/zerolog/log"
)

// Stats sort options
const (
	StatSortByPodID = 0 + iota
	StatSortByContainerName
	StatSortByCPUPerc
	StatSortByMemPerc
)

// StatReporter implements pod stats metrics
type StatReporter struct {
	entities.PodStatsReport
}

// StatsOptions pod stats query option
type StatsOptions struct {
	IDs    []string
	SortBy int
}

// Stats returns resource-usage statistics of a pod.
func Stats(opts *StatsOptions) ([]StatReporter, error) {
	log.Debug().Msgf("pdcs: podman pods stats %v", *opts)
	conn, err := connection.GetConnection()
	if err != nil {
		return nil, err
	}
	statReport, err := pods.Stats(conn, opts.IDs, nil)
	if err != nil {
		return nil, err
	}

	report := sortStats(statReport, opts.SortBy)
	return report, nil
}

func sortStats(podSReport []*entities.PodStatsReport, sortBy int) []StatReporter {
	report := make([]StatReporter, 0, len(podSReport))
	for _, item := range podSReport {
		var reporterItem StatReporter
		reporterItem.Pod = item.Pod
		reporterItem.CID = item.CID
		reporterItem.Name = item.Name
		reporterItem.CPU = item.CPU
		reporterItem.MemUsage = item.MemUsage
		reporterItem.Mem = item.Mem
		reporterItem.NetIO = item.NetIO
		reporterItem.BlockIO = item.BlockIO
		reporterItem.PIDS = item.PIDS
		report = append(report, reporterItem)
	}
	sort.Slice(report, sortFunc(sortBy, report))
	return report
}

func sortFunc(key int, data []StatReporter) func(i, j int) bool {
	switch key {
	case StatSortByContainerName:
		return func(i, j int) bool {
			return data[i].CID < data[j].CID
		}
	case StatSortByCPUPerc:
		return func(i, j int) bool {
			return data[i].cpuPerc() > data[j].cpuPerc()
		}
	case StatSortByMemPerc:
		return func(i, j int) bool {
			return data[i].memPerc() > data[j].memPerc()
		}
	default:
		// case "StatSortByPodID":
		return func(i, j int) bool {
			return data[i].Pod < data[j].Pod
		}
	}
}

func (sreport StatReporter) cpuPerc() float64 {
	return percentageToFloat(sreport.CPU)
}

func (sreport StatReporter) memPerc() float64 {
	return percentageToFloat(sreport.Mem)
}

func percentageToFloat(text string) float64 {
	text = strings.ReplaceAll(text, "%", "")
	value, _ := strconv.ParseFloat(text, 64)
	return value
}
