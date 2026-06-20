package tvxwidgets

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	barChartYAxisLabelWidth = 2
	barGap                  = 2
	barWidth                = 3
)

// BarChartItem represents a single bar in bar chart.
type BarChartItem struct {
	label string
	value int
	color tcell.Color
}

// BarChart represents bar chart primitive.
type BarChart struct {
	*tview.Box

	// bar items
	bars []BarChartItem
	// maximum value of bars
	maxVal int
	// barGap gap between two bars
	barGap int
	// barWidth width of bars
	barWidth int
	// hasBorder true if primitive has border
	hasBorder      bool
	axesColor      tcell.Color
	axesLabelColor tcell.Color
}

// NewBarChart returns a new bar chart primitive.
func NewBarChart() *BarChart {
	chart := &BarChart{
		Box:            tview.NewBox(),
		barGap:         barGap,
		barWidth:       barWidth,
		axesColor:      tcell.ColorDimGray,
		axesLabelColor: tcell.ColorDimGray,
	}

	return chart
}

// Focus is called when this primitive receives focus.
func (c *BarChart) Focus(delegate func(p tview.Primitive)) {
	delegate(c.Box)
}

// HasFocus returns whether or not this primitive has focus.
func (c *BarChart) HasFocus() bool {
	return c.Box.HasFocus()
}

// Draw draws this primitive onto the screen.
func (c *BarChart) Draw(screen tcell.Screen) { //nolint:funlen,cyclop
	c.DrawForSubclass(screen, c)

	x, y, width, height := c.GetInnerRect()

	maxValY := y + 1
	xAxisStartY := y + height - 2 //nolint:mnd
	barStartY := y + height - 3   //nolint:mnd
	borderPadding := 0

	if c.hasBorder {
		borderPadding = 1
	}
	// set max value if not set
	c.initMaxValue()
	maxValueSr := strconv.Itoa(c.maxVal)
	maxValLenght := len(maxValueSr) + 1

	if maxValLenght < barChartYAxisLabelWidth {
		maxValLenght = barChartYAxisLabelWidth
	}

	axesStyle := tcell.StyleDefault.Background(c.GetBackgroundColor()).Foreground(c.axesColor)
	axesLabelStyle := tcell.StyleDefault.Background(c.GetBackgroundColor()).Foreground(c.axesLabelColor)

	// draw Y axis line
	drawLine(screen,
		x+maxValLenght,
		y+borderPadding,
		height-borderPadding-1,
		verticalLine, axesStyle)

	// draw X axis line
	drawLine(screen,
		x+maxValLenght+1,
		xAxisStartY,
		width-borderPadding-maxValLenght-1,
		horizontalLine, axesStyle)

	tview.PrintJoinedSemigraphics(screen,
		x+maxValLenght,
		xAxisStartY,
		tview.BoxDrawingsLightUpAndRight, axesStyle)

	tview.PrintJoinedSemigraphics(screen, x+maxValLenght-1, xAxisStartY, '0', axesLabelStyle)

	mxValRune := []rune(maxValueSr)
	for i := range mxValRune {
		tview.PrintJoinedSemigraphics(screen, x+borderPadding+i, maxValY, mxValRune[i], axesLabelStyle)
	}

	// draw bars
	startX := x + maxValLenght + c.barGap
	labelY := y + height - 1
	valueMaxHeight := barStartY - maxValY

	for _, item := range c.bars {
		if startX > x+width {
			return
		}
		// set labels
		r := []rune(item.label)
		for j := range r {
			tview.PrintJoinedSemigraphics(screen, startX+j, labelY, r[j], axesLabelStyle)
		}
		// bar style
		bStyle := tcell.StyleDefault.Background(c.GetBackgroundColor()).Foreground(item.color)
		barHeight := c.getHeight(valueMaxHeight, item.value)

		for k := range barHeight {
			for l := range c.barWidth {
				tview.PrintJoinedSemigraphics(screen, startX+l, barStartY-k, fullBlockRune, bStyle)
			}
		}
		// bar value
		vSt := strconv.Itoa(item.value)
		vRune := []rune(vSt)

		for i := range vRune {
			tview.PrintJoinedSemigraphics(screen, startX+i, barStartY-barHeight, vRune[i], bStyle)
		}

		// calculate next startX for next bar
		rWidth := len(r)
		if rWidth < c.barWidth {
			rWidth = c.barWidth
		}

		startX = startX + c.barGap + rWidth
	}
}

// SetBorder sets border for this primitive.
func (c *BarChart) SetBorder(status bool) {
	c.hasBorder = status
	c.Box.SetBorder(status)
}

// GetRect return primitive current rect.
func (c *BarChart) GetRect() (int, int, int, int) {
	return c.Box.GetRect()
}

// SetRect sets rect for this primitive.
func (c *BarChart) SetRect(x, y, width, height int) {
	c.Box.SetRect(x, y, width, height)
}

// SetMaxValue sets maximum value of bars.
func (c *BarChart) SetMaxValue(maxValue int) {
	c.maxVal = maxValue
}

// SetAxesColor sets axes x and y lines color.
func (c *BarChart) SetAxesColor(color tcell.Color) {
	c.axesColor = color
}

// SetAxesLabelColor sets axes x and y label color.
func (c *BarChart) SetAxesLabelColor(color tcell.Color) {
	c.axesLabelColor = color
}

// AddBar adds new bar item to the bar chart primitive.
func (c *BarChart) AddBar(label string, value int, color tcell.Color) {
	c.bars = append(c.bars, BarChartItem{
		label: label,
		value: value,
		color: color,
	})
}

// RemoveBar removes a bar item from the bar chart.
func (c *BarChart) RemoveBar(label string) {
	bars := c.bars[:0]

	for _, barItem := range c.bars {
		if barItem.label != label {
			bars = append(bars, barItem)
		}
	}

	c.bars = bars
}

// SetBarValue sets bar values.
func (c *BarChart) SetBarValue(name string, value int) {
	for i := range c.bars {
		if c.bars[i].label == name {
			c.bars[i].value = value
		}
	}
}

func (c *BarChart) getHeight(maxHeight int, value int) int {
	if value >= c.maxVal {
		return maxHeight
	}

	height := (value * maxHeight) / c.maxVal

	return height
}

func (c *BarChart) initMaxValue() {
	// set max value if not set
	if c.maxVal == 0 {
		for _, b := range c.bars {
			if b.value > c.maxVal {
				c.maxVal = b.value
			}
		}
	}
}
