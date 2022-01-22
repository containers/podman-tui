package dialogs

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	prgCell     = "▉"
	prgMinWidth = 40
)

// ProgressDialog represents progress bar permitive
type ProgressDialog struct {
	*tview.Box
	x            int
	y            int
	width        int
	height       int
	counterValue int
	display      bool
}

// NewProgressDialog returns new progress dialog primitive
func NewProgressDialog() *ProgressDialog {
	return &ProgressDialog{
		Box:     tview.NewBox().SetBorder(true),
		display: false,
	}
}

// SetTitle sets title for this primitive
func (d *ProgressDialog) SetTitle(title string) {
	d.Box.SetTitle(title)
}

// Draw draws this primitive onto the screen.
func (d *ProgressDialog) Draw(screen tcell.Screen) {
	if !d.display || d.height < 3 {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, _ := d.Box.GetInnerRect()
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
		spaceWidth := (width - d.width) / 2
		d.x = x + spaceWidth
	}
	if height > 3 {
		d.height = 3
		spaceHeight := (height - d.height) / 2
		d.y = y + spaceHeight
	}

	d.Box.SetRect(d.x, d.y, d.width, d.height)
}

// Hide Hides this primitive
func (d *ProgressDialog) Hide() {
	d.display = false
}

// Display displays this primitive
func (d *ProgressDialog) Display() {
	d.counterValue = 0
	d.display = true

}

// IsDisplay returns true if primitive is shown
func (d *ProgressDialog) IsDisplay() bool {
	return d.display
}

// Focus is called when this primitive receives focus
func (d *ProgressDialog) Focus(delegate func(p tview.Primitive)) {
}

// HasFocus returns whether or not this primitive has focus
func (d *ProgressDialog) HasFocus() bool {
	return d.Box.HasFocus()
}

// InputHandler returns input handler function for this primitive
func (d *ProgressDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("progress dialog: event %v received", event.Key())

	})
}

func (d *ProgressDialog) tickStr(max int) string {
	counter := d.counterValue
	if counter < max-4 {
		d.counterValue++
	} else {
		d.counterValue = 0
	}
	prgHeadStr := ""
	hWidth := 0
	prgEndStr := ""
	prgStr := ""
	for i := 0; i < d.counterValue; i++ {
		prgHeadStr = prgHeadStr + fmt.Sprintf("[black::]%s", prgCell)
		hWidth++
	}
	prgStr = prgCell + prgCell + prgCell + prgCell
	for i := 0; i < max+hWidth+4; i++ {
		prgEndStr = prgEndStr + fmt.Sprintf("[black::]%s", prgCell)
	}

	return fmt.Sprintf("%s[orange::]%s%s", prgHeadStr, prgStr, prgEndStr)
}
