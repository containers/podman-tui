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

// ProgressUsageString return progressbar string (bars + usage percentage)
func ProgressUsageString(percentage float64) string {
	progressCell := ""
	value := int(int(percentage) * (prgWidth) / 100)
	for index := 0; index < prgWidth; index++ {
		if index < value {
			progressCell = progressCell + getBarColor(index)

		} else {
			progressCell = progressCell + style.ProgressBarCell
		}
	}
	return progressCell + fmt.Sprintf("%6.2f%%", percentage)
}

func getBarColor(value int) string {

	barCell := ""
	barColor := ""

	if value < prgWarn {
		barColor = style.GetColorName(style.PrgBarOKColor)
	} else if value < prgCrit {
		barColor = style.GetColorName(style.PrgBarWarnColor)
	} else {
		barColor = style.GetColorName(style.PrgBarCritColor)
	}
	barCell = fmt.Sprintf("[%s::]%s[%s::]", barColor, style.ProgressBarCell, style.GetColorName(style.PrgBarEmptyColor))
	return barCell
}
