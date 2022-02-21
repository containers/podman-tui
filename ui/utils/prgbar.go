package utils

import (
	"fmt"
)

const (
	prgCell  = "▉"
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
			progressCell = progressCell + prgCell
		}
	}
	return progressCell + fmt.Sprintf("%6.2f%%", percentage)
}

func getBarColor(value int) string {

	barCell := ""
	barColor := ""

	if value < prgWarn {
		barColor = GetColorName(Styles.InfoBar.ProgressBar.BarOKColor)
	} else if value < prgCrit {
		barColor = GetColorName(Styles.InfoBar.ProgressBar.BarWarnColor)
	} else {
		barColor = GetColorName(Styles.InfoBar.ProgressBar.BarCritColor)
	}
	barCell = fmt.Sprintf("[%s::]%s[%s::]", barColor, prgCell, GetColorName(Styles.InfoBar.ProgressBar.BarEmptyColor))
	return barCell
}
