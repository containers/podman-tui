package utils

import (
	"fmt"

	"github.com/containers/podman-tui/ui/style"
)

const (
	prgWidth = 20
	prgWarn  = 13
	prgCrit  = 17
)

// ProgressUsageString return progressbar string (bars + usage percentage).
func ProgressUsageString(percentage float64) string {
	progressCell := ""
	value := int(percentage) * (prgWidth) / 100 //nolint:mnd

	for index := range prgWidth {
		if index < value {
			progressCell += getBarColor(index)
		} else {
			progressCell += style.ProgressBarCell
		}
	}

	return progressCell + fmt.Sprintf("%6.2f%%", percentage)
}

func getBarColor(value int) string {
	barCell := ""
	barColor := ""

	switch {
	case value < prgWarn:
		barColor = style.GetColorName(style.PrgBarOKColor)
	case value < prgCrit:
		barColor = style.GetColorName(style.PrgBarWarnColor)
	default:
		barColor = style.GetColorName(style.PrgBarCritColor)
	}

	barCell = fmt.Sprintf("[%s::]%s[%s::]",
		barColor, style.ProgressBarCell, style.GetColorName(style.PrgBarEmptyColor))

	return barCell
}
