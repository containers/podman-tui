package style

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// ThemeDefault is podman-tui's own color theme.
	ThemeDefault = "default"
	// ThemeTerminal uses only the terminal's default foreground/background and
	// its 16 ANSI colors, so the interface follows the terminal's color scheme
	// (light or dark), including a scheme the terminal changes at runtime.
	ThemeTerminal = "terminal"
)

// ErrUnknownTheme is returned by SetTheme for a theme name it does not know.
var ErrUnknownTheme = errors.New("unknown theme")

var terminalTheme bool

// ansiColorNames are the tcell/tview names of the 16 ANSI colors, by index.
var ansiColorNames = [...]string{
	"black", "maroon", "green", "olive", "navy", "purple", "teal", "silver",
	"gray", "red", "lime", "yellow", "blue", "fuchsia", "aqua", "white",
}

// SetTheme applies the named color theme. It must be called before the
// interface is created.
func SetTheme(name string) error {
	switch name {
	case "", ThemeDefault:
		return nil
	case ThemeTerminal:
		setTerminalTheme()

		return nil
	}

	return fmt.Errorf("%w %q (available: %s, %s)", ErrUnknownTheme, name, ThemeDefault, ThemeTerminal)
}

func setTerminalTheme() {
	// focused and selected elements: reverse video in the terminal's own colors
	highlight := tcell.StyleDefault.Reverse(true)

	terminalTheme = true

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.PrimaryTextColor = tcell.ColorDefault
	tview.Styles.TitleColor = tcell.ColorDefault

	InfoBarItemFgColor = tcell.ColorSilver
	FgColor = tcell.ColorDefault
	BgColor = tcell.ColorDefault
	BorderColor = tcell.ColorBlue
	HelpHeaderFgColor = tcell.ColorBlue
	MenuBgColor = tcell.ColorPurple
	PageHeaderBgColor = tcell.ColorPurple
	PageHeaderFgColor = tcell.ColorBlack
	RunningStatusFgColor = tcell.ColorGreen
	PausedStatusFgColor = tcell.ColorOlive
	DialogBgColor = tcell.ColorDefault
	DialogBorderColor = tcell.ColorPurple
	DialogFgColor = tcell.ColorDefault
	DialogSubBoxBorderColor = tcell.ColorDefault
	ErrorDialogBgColor = tcell.ColorDefault
	ErrorDialogButtonBgColor = tcell.ColorMaroon
	TerminalFgColor = tcell.ColorDefault
	TerminalBgColor = tcell.ColorDefault
	TerminalBorderColor = tcell.ColorDefault
	TableHeaderBgColor = tcell.ColorPurple
	TableHeaderFgColor = tcell.ColorBlack
	PrgBgColor = tcell.ColorGray
	PrgBarColor = tcell.ColorTeal
	PrgBarEmptyColor = tcell.ColorDefault
	PrgBarOKColor = tcell.ColorGreen
	PrgBarWarnColor = tcell.ColorOlive
	PrgBarCritColor = tcell.ColorMaroon
	DropDownUnselected = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorDefault)
	DropDownSelected = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorBlack)
	DropDownFocused = highlight
	InputLabelStyle = tcell.StyleDefault.Background(DialogBgColor).Foreground(DialogFgColor)
	InputFieldStyle = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(DialogFgColor)
	FieldBackgroundColor = tcell.ColorGray
	ButtonBgColor = tcell.ColorPurple
	ButtonFgColor = tcell.ColorBlack
	TableSelectedStyle = highlight
	DialogSelectedStyle = highlight
	ButtonActivatedStyle = highlight
	CheckboxActivatedStyle = highlight
	ErrorButtonActivated = highlight.Foreground(ErrorDialogButtonBgColor)
}

// terminalColorTag names a color for a text tag in the terminal theme so that
// the terminal's own palette resolves it: "-" for the default color and the
// ANSI name for the 16 ANSI colors. Any other color is not handled.
func terminalColorTag(color tcell.Color) (string, bool) {
	switch {
	case !terminalTheme:
		return "", false
	case color == tcell.ColorDefault:
		return "-", true
	case color >= tcell.ColorBlack && color <= tcell.ColorWhite:
		return ansiColorNames[color-tcell.ColorBlack], true
	}

	return "", false
}
