package tvxwidgets

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Spartline represents a sparkline widgets.
type Sparkline struct {
	*tview.Box

	data           []float64
	dataTitle      string
	dataTitlecolor tcell.Color
	lineColor      tcell.Color
	mu             sync.Mutex
}

// NewSparkline returns a a new sparkline widget.
func NewSparkline() *Sparkline {
	return &Sparkline{
		Box: tview.NewBox(),
	}
}

// Draw draws this primitive onto the screen.
func (sl *Sparkline) Draw(screen tcell.Screen) {
	sl.Box.DrawForSubclass(screen, sl)

	x, y, width, height := sl.Box.GetInnerRect()
	barHeight := height

	// print label
	if sl.dataTitle != "" {
		tview.Print(screen, sl.dataTitle, x, y, width, tview.AlignLeft, sl.dataTitlecolor)
		barHeight--
	}

	maxVal := getMaxFloat64FromSlice(sl.data)
	if maxVal == 0 {
		return
	}

	// print lines
	for i := 0; i < len(sl.data) && i+x < x+width; i++ {
		data := sl.data[i]
		dHeight := int((data / maxVal) * float64(barHeight))

		sparkChar := barsRune[len(barsRune)-1]

		style := tcell.StyleDefault.Background(sl.GetBackgroundColor()).Foreground(sl.lineColor)

		for j := 0; j < dHeight; j++ {
			tview.PrintJoinedSemigraphics(screen, i+x, y-1+height-j, sparkChar, style)
		}

		if dHeight == 0 {
			sparkChar = barsRune[1]
			tview.PrintJoinedSemigraphics(screen, i+x, y-1+height, sparkChar, style)
		}
	}
}

// SetRect sets rect for this primitive.
func (sl *Sparkline) SetRect(x, y, width, height int) {
	sl.Box.SetRect(x, y, width, height)
}

// GetRect return primitive current rect.
func (sl *Sparkline) GetRect() (int, int, int, int) {
	return sl.Box.GetRect()
}

// HasFocus returns whether or not this primitive has focus.
func (sl *Sparkline) HasFocus() bool {
	return sl.Box.HasFocus()
}

// SetData sets sparkline data.
func (sl *Sparkline) SetData(data []float64) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.data = data
}

// SetDataTitle sets sparkline data title.
func (sl *Sparkline) SetDataTitle(title string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.dataTitle = title
}

// SetDataTitleColor sets sparkline data title color.
func (sl *Sparkline) SetDataTitleColor(color tcell.Color) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.dataTitlecolor = color
}

// SetLineColor sets sparkline line color.
func (sl *Sparkline) SetLineColor(color tcell.Color) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.lineColor = color
}
