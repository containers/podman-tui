package tvxwidgets

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type drawLineMode int

const (
	horizontalLine drawLineMode = iota
	verticalLine
)

const (
	// gauge cell.
	prgCell = "▉"
	// form height.
	dialogFormHeight = 3
	// gauge warning percentage.
	gaugeWarnPc = 60.00
	// gauge critical percentage.
	gaugeCritPc = 85.00
	// gauge min percentage.
	gaugeMinPc = 0.00
	// gauge max percentage.
	gaugeMaxPc = 100
	// dialog padding.
	dialogPadding = 2
	// empty space parts.
	emptySpaceParts   = 2
	brailleOffsetRune = '\u2800'
	dotRune           = '\u25CF'
)

var (
	brailleRune = [4][2]rune{ // nolint:gochecknoglobals
		{'\u0001', '\u0008'},
		{'\u0002', '\u0010'},
		{'\u0004', '\u0020'},
		{'\u0040', '\u0080'},
	}

	barsRune = [...]rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'} // nolint:gochecknoglobals
)

// getColorName returns convert tcell color to its name.
func getColorName(color tcell.Color) string {
	for name, c := range tcell.ColorNames {
		if c == color {
			return name
		}
	}

	return ""
}

// getMessageWidth returns width size for dialogs based on messages.
func getMessageWidth(message string) int {
	var messageWidth int
	for _, msg := range strings.Split(message, "\n") {
		if len(msg) > messageWidth {
			messageWidth = len(msg)
		}
	}

	return messageWidth
}

// returns max values in 2D float64 slices.
func getMaxFloat64From2dSlice(slices [][]float64) float64 {
	if len(slices) == 0 {
		return 0
	}

	var max float64

	for _, slice := range slices {
		for _, val := range slice {
			if val > max {
				max = val
			}
		}
	}

	return max
}

// returns max values in float64 slices.
func getMaxFloat64FromSlice(slice []float64) float64 {
	if len(slice) == 0 {
		return 0
	}

	var max float64

	for _, val := range slice {
		if val > max {
			max = val
		}
	}

	return max
}

func absInt(x int) int {
	if x >= 0 {
		return x
	}

	return -x
}

func drawLine(screen tcell.Screen, x int, y int, length int, mode drawLineMode, style tcell.Style) {
	if mode == horizontalLine {
		for i := 0; i < length; i++ {
			tview.PrintJoinedSemigraphics(screen, x+i, y, tview.BoxDrawingsLightTripleDashHorizontal, style)
		}
	} else if mode == verticalLine {
		for i := 0; i < length; i++ {
			tview.PrintJoinedSemigraphics(screen, x, y+i, tview.BoxDrawingsLightTripleDashVertical, style)
		}
	}
}
