package tvxwidgets

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// PercentageModeGauge represents percentage mode gauge permitive.
type PercentageModeGauge struct {
	*tview.Box
	// maxValue value
	maxValue int
	// value is current value
	value int
	// pgBgColor: progress block background color
	pgBgColor tcell.Color
}

// NewPercentageModeGauge returns new percentage mode gauge permitive.
func NewPercentageModeGauge() *PercentageModeGauge {
	gauge := &PercentageModeGauge{
		Box:       tview.NewBox(),
		value:     0,
		pgBgColor: tcell.ColorBlue,
	}

	return gauge
}

// Draw draws this primitive onto the screen.
func (g *PercentageModeGauge) Draw(screen tcell.Screen) {
	g.Box.DrawForSubclass(screen, g)

	if g.maxValue == 0 {
		return
	}

	x, y, width, height := g.Box.GetInnerRect()
	pcWidth := 3
	pc := g.value * gaugeMaxPc / g.maxValue
	pcString := fmt.Sprintf("%d%%", pc)
	tW := width - pcWidth
	tX := x + (tW / emptySpaceParts)
	tY := y + height/emptySpaceParts
	prgBlock := g.progressBlock(width)
	style := tcell.StyleDefault.Background(g.pgBgColor).Foreground(tview.Styles.PrimaryTextColor)

	for i := 0; i < height; i++ {
		for j := 0; j < prgBlock; j++ {
			screen.SetContent(x+j, y+i, ' ', nil, style)
		}
	}

	// print percentage in middle of box

	pcRune := []rune(pcString)
	for j := 0; j < len(pcRune); j++ {
		style = tcell.StyleDefault.Background(tview.Styles.PrimitiveBackgroundColor).Foreground(tview.Styles.PrimaryTextColor)
		if x+prgBlock >= tX+j {
			style = tcell.StyleDefault.Background(g.pgBgColor).Foreground(tview.Styles.PrimaryTextColor)
		}

		for i := 0; i < height; i++ {
			screen.SetContent(tX+j, y+i, ' ', nil, style)
		}
		screen.SetContent(tX+j, tY, pcRune[j], nil, style)
	}
}

// SetTitle sets title for this primitive.
func (g *PercentageModeGauge) SetTitle(title string) {
	g.Box.SetTitle(title)
}

// Focus is called when this primitive receives focus.
func (g *PercentageModeGauge) Focus(delegate func(p tview.Primitive)) {
}

// HasFocus returns whether or not this primitive has focus.
func (g *PercentageModeGauge) HasFocus() bool {
	return g.Box.HasFocus()
}

// GetRect return primitive current rect.
func (g *PercentageModeGauge) GetRect() (int, int, int, int) {
	return g.Box.GetRect()
}

// SetRect sets rect for this primitive.
func (g *PercentageModeGauge) SetRect(x, y, width, height int) {
	g.Box.SetRect(x, y, width, height)
}

// SetPgBgColor sets progress block background color.
func (g *PercentageModeGauge) SetPgBgColor(color tcell.Color) {
	g.pgBgColor = color
}

// SetValue update the gauge progress.
func (g *PercentageModeGauge) SetValue(value int) {
	if value <= g.maxValue {
		g.value = value
	}
}

// GetValue returns current gauge value.
func (g *PercentageModeGauge) GetValue() int {
	return g.value
}

// SetMaxValue set maximum allows value for the gauge.
func (g *PercentageModeGauge) SetMaxValue(value int) {
	if value > 0 {
		g.maxValue = value
	}
}

// GetMaxValue returns maximum allows value for the gauge.
func (g *PercentageModeGauge) GetMaxValue() int {
	return g.maxValue
}

// Reset resets the gauge counter (set to 0).
func (g *PercentageModeGauge) Reset() {
	g.value = 0
}

func (g *PercentageModeGauge) progressBlock(max int) int {
	if g.maxValue == 0 {
		return g.maxValue
	}

	pc := g.value * gaugeMaxPc / g.maxValue
	value := pc * max / gaugeMaxPc

	return value
}
