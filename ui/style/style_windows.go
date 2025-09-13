//go:build windows

package style

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// HeavyGreenCheckMark unicode.
	HeavyGreenCheckMark = "[green::]\u25CF[-::]"
	// HeavyRedCrossMark unicode.
	HeavyRedCrossMark = "[red::]\u25CF[-::]"
	// ProgressBar cell.
	ProgressBarCell = "\u2593"
)

var (
	// infobar.
	InfoBarItemFgColor = tcell.ColorSilver
	// main views.
	FgColor              = tcell.ColorFloralWhite
	BgColor              = tview.Styles.PrimitiveBackgroundColor
	BorderColor          = tcell.NewRGBColor(135, 135, 175) //nolint:mnd
	HelpHeaderFgColor    = tcell.NewRGBColor(135, 135, 175) //nolint:mnd
	MenuBgColor          = tcell.ColorMediumPurple
	PageHeaderBgColor    = tcell.ColorMediumPurple
	PageHeaderFgColor    = tcell.ColorFloralWhite
	RunningStatusFgColor = tcell.NewRGBColor(95, 215, 0)  //nolint:mnd
	PausedStatusFgColor  = tcell.NewRGBColor(255, 175, 0) //nolint:mnd
	// dialogs.
	DialogBgColor            = tcell.NewRGBColor(38, 38, 38) //nolint:mnd
	DialogBorderColor        = tcell.ColorMediumPurple
	DialogFgColor            = tcell.ColorFloralWhite
	DialogSubBoxBorderColor  = tcell.ColorDimGray
	ErrorDialogBgColor       = tcell.NewRGBColor(215, 0, 0) //nolint:mnd
	ErrorDialogButtonBgColor = tcell.ColorDarkRed
	// terminal.
	TerminalFgColor     = tcell.ColorFloralWhite
	TerminalBgColor     = tcell.NewRGBColor(5, 5, 5) //nolint:mnd
	TerminalBorderColor = tcell.ColorDimGray
	// table header.
	TableHeaderBgColor = tcell.ColorMediumPurple
	TableHeaderFgColor = tcell.ColorFloralWhite
	// progress bar.
	PrgBgColor       = tcell.ColorDimGray
	PrgBarColor      = tcell.ColorDarkOrange
	PrgBarEmptyColor = tcell.ColorWhite
	PrgBarOKColor    = tcell.ColorGreen
	PrgBarWarnColor  = tcell.ColorOrange
	PrgBarCritColor  = tcell.ColorRed
	// dropdown.
	DropDownUnselected = tcell.StyleDefault.Background(tcell.ColorWhiteSmoke).Foreground(tcell.ColorBlack)
	DropDownSelected   = tcell.StyleDefault.Background(tcell.ColorLightSlateGray).Foreground(tcell.ColorWhite)
	DropDownFocused    = tcell.StyleDefault.Background(tcell.ColorWhiteSmoke).Foreground(tcell.ColorBlack)
	// other primitives.
	InputLabelStyle      = tcell.StyleDefault.Background(DialogBgColor).Foreground(DialogFgColor)
	InputFieldStyle      = tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)
	FieldBackgroundColor = tcell.ColorDarkGray
	ButtonBgColor        = tcell.ColorMediumPurple
)

// GetColorName returns convert tcell color to its name.
func GetColorName(color tcell.Color) string {
	for name, c := range tcell.ColorNames {
		if c == color {
			return name
		}
	}
	return ""
}

// GetColorHex shall returns convert tcell color to its hex useful for textview primitives,
// however, for windows nodes it will return color name.
func GetColorHex(color tcell.Color) string {
	for name, c := range tcell.ColorNames {
		if c == color {
			return name
		}
	}
	return ""
}
