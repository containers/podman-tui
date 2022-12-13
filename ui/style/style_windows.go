//go:build windows

package style

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// HeavyGreenCheckMark unicode
	HeavyGreenCheckMark = "[green::]\u25CF[-::]"
	// HeavyRedCrossMark unicode
	HeavyRedCrossMark = "[red::]\u25CF[-::]"
	// ProgressBar cell
	ProgressBarCell = "\u2593"
)

var (
	// infobar
	InfoBarItemFgColor = tcell.ColorGray
	// main views
	FgColor              = tview.Styles.PrimaryTextColor
	BgColor              = tview.Styles.PrimitiveBackgroundColor
	BorderColor          = tcell.ColorPink
	MenuBgColor          = tcell.ColorPink
	HelpHeaderFgColor    = tcell.ColorPink
	PageHeaderBgColor    = tcell.ColorPink
	PageHeaderFgColor    = tview.Styles.PrimaryTextColor
	RunningStatusFgColor = tcell.ColorLime
	PausedStatusFgColor  = tcell.ColorYellow

	// dialogs
	DialogBgColor            = tview.Styles.PrimitiveBackgroundColor
	DialogFgColor            = tview.Styles.PrimaryTextColor
	DialogBorderColor        = tcell.ColorPink
	DialogSubBoxBorderColor  = tcell.ColorGray
	ErrorDialogBgColor       = tcell.ColorRed
	ErrorDialogButtonBgColor = tcell.ColorPink
	// terminal
	TerminalBgColor     = tview.Styles.PrimitiveBackgroundColor
	TerminalFgColor     = tview.Styles.PrimaryTextColor
	TerminalBorderColor = tview.Styles.PrimitiveBackgroundColor
	// table header
	TableHeaderBgColor = tcell.ColorPink
	TableHeaderFgColor = tview.Styles.PrimaryTextColor
	// progress bar
	PrgBgColor       = tview.Styles.PrimaryTextColor
	PrgBarColor      = tcell.ColorFuchsia
	PrgBarEmptyColor = tcell.ColorWhite
	PrgBarOKColor    = tcell.ColorLime
	PrgBarWarnColor  = tcell.ColorYellow
	PrgBarCritColor  = tcell.ColorRed
	// dropdown
	DropDownUnselected = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	DropDownSelected   = tcell.StyleDefault.Background(tcell.ColorPurple).Foreground(tview.Styles.PrimaryTextColor)
	// other primitives
	InputFieldBgColor = tcell.ColorGray
	ButtonBgColor     = tcell.ColorPink
)

// GetColorName returns convert tcell color to its name
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
