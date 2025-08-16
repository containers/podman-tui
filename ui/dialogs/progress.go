package dialogs

import (
	"fmt"

	"github.com/containers/podman-tui/ui/style"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	prgCell     = "▉"
	prgMinWidth = 40
)

// ProgressDialog represents progress bar permitive.
type ProgressDialog struct {
	*tview.Box

	x            int
	y            int
	width        int
	height       int
	counterValue int
	display      bool
}

// NewProgressDialog returns new progress dialog primitive.
func NewProgressDialog() *ProgressDialog {
	return &ProgressDialog{
		Box: tview.NewBox().
			SetBorder(true).
			SetBorderColor(style.BorderColor),
		display: false,
	}
}

// SetTitle sets title for this primitive.
func (d *ProgressDialog) SetTitle(title string) {
	d.Box.SetTitle(title)
}

// Draw draws this primitive onto the screen.
func (d *ProgressDialog) Draw(screen tcell.Screen) {
	if !d.display || d.height < 3 {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, _ := d.GetInnerRect()
	tickStr := d.tickStr(width)
	tview.Print(screen, tickStr, x, y, width, tview.AlignLeft, tcell.ColorYellow)
}

// SetRect set rects for this primitive.
func (d *ProgressDialog) SetRect(x, y, width, height int) {
	d.x = x
	d.y = y
	d.width = width

	if d.width > prgMinWidth {
		d.width = prgMinWidth
		spaceWidth := (width - d.width) / 2 //nolint:mnd
		d.x = x + spaceWidth
	}

	if height > 3 { //nolint:mnd
		d.height = 3
		spaceHeight := (height - d.height) / 2 //nolint:mnd
		d.y = y + spaceHeight
	}

	d.Box.SetRect(d.x, d.y, d.width, d.height)
}

// Hide Hides this primitive.
func (d *ProgressDialog) Hide() {
	d.display = false
}

// Display displays this primitive.
func (d *ProgressDialog) Display() {
	d.counterValue = 0
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ProgressDialog) IsDisplay() bool {
	return d.display
}

// Focus is called when this primitive receives focus.
func (d *ProgressDialog) Focus(delegate func(p tview.Primitive)) {}

// HasFocus returns whether or not this primitive has focus.
func (d *ProgressDialog) HasFocus() bool {
	return d.Box.HasFocus()
}

// InputHandler returns input handler function for this primitive.
func (d *ProgressDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("progress dialog: event %v received", event)
	})
}

func (d *ProgressDialog) tickStr(maxCount int) string {
	prgStr := prgCell + prgCell + prgCell + prgCell
	prgHeadStr := ""
	hWidth := 0
	prgEndStr := ""
	barColor := style.GetColorHex(style.PrgBarColor)
	counter := d.counterValue

	if counter < maxCount-4 {
		d.counterValue++
	} else {
		d.counterValue = 0
	}

	for range d.counterValue {
		prgHeadStr += fmt.Sprintf("[black::]%s", prgCell) //nolint:perfsprint
		hWidth++
	}

	for range maxCount + hWidth + 4 {
		prgEndStr += fmt.Sprintf("[black::]%s", prgCell) //nolint:perfsprint
	}

	progress := fmt.Sprintf("%s[%s::]%s%s", prgHeadStr, barColor, prgStr, prgEndStr)

	return progress
}
