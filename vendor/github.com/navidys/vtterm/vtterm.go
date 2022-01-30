package vtterm

import (
	"fmt"
	"image/color"
)

// Intensity represents cell intensity
type Intensity int

const (
	// Normal intensity
	Normal Intensity = 0
	// Bright intensity
	Bright = 1
	// Dim intensity
	Dim = 2
	// TODO(jaguilar): Should this be in a subpackage, since the names are pretty collide-y?
)

var (
	// DefaultColor technically RGBAs are supposed to be premultiplied. But CSS doesn't expect them
	// that way, so we won't do it in this file.
	DefaultColor = color.RGBA{0, 0, 0, 0}
	// Black has 255 alpha, so it will compare negatively with DefaultColor.
	Black = color.RGBA{0, 0, 0, 255}
	// Red color
	Red = color.RGBA{255, 0, 0, 255}
	// Green color
	Green = color.RGBA{0, 255, 0, 255}
	// Yellow color
	Yellow = color.RGBA{255, 255, 0, 255}
	// Blue color
	Blue = color.RGBA{0, 0, 255, 255}
	// Magenta color
	Magenta = color.RGBA{255, 0, 255, 255}
	// Cyan color
	Cyan = color.RGBA{0, 255, 255, 255}
	// White color
	White = color.RGBA{255, 255, 255, 255}
)

// Format represents the display format of text on a terminal.
type Format struct {
	// Fg is the foreground color.
	Fg color.RGBA
	// Bg is the background color.
	Bg color.RGBA
	// Intensity is the text intensity (bright, normal, dim).
	Intensity Intensity
	// Various text properties.
	Underscore, Conceal, Negative, Blink, Inverse bool
}

// Cursor represents both the position and text type of the cursor.
type Cursor struct {
	// Y and X are the coordinates.
	Y, X int

	// F is the format that will be displayed.
	F Format
}

// VT100 represents a simplified, raw VT100 terminal.
type VT100 struct {
	// Height and Width are the dimensions of the terminal.
	Height, Width int

	// Content is the text in the terminal.
	Content [][]rune

	// Format is the display properties of each cell.
	Format [][]Format

	// Cursor is the current state of the cursor.
	Cursor Cursor

	// savedCursor is the state of the cursor last time save() was called.
	savedCursor Cursor
}

// NewVT100 creates a new VT100 object with the specified dimensions. y and x
// must both be greater than zero.
//
// Each cell is set to contain a ' ' rune, and all formats are left as the
// default.
func NewVT100(y, x int) *VT100 {
	if y == 0 || x == 0 {
		panic(fmt.Errorf("invalid dim (%d, %d)", y, x))
	}

	v := &VT100{
		Height:  y,
		Width:   x,
		Content: make([][]rune, y),
		Format:  make([][]Format, y),
	}

	for row := 0; row < y; row++ {
		v.Content[row] = make([]rune, x)
		v.Format[row] = make([]Format, x)

		for col := 0; col < x; col++ {
			v.clear(row, col)
		}
	}
	return v
}

// Process handles a single ANSI terminal command, updating the terminal
// appropriately.
//
// One special kind of error that this can return is an UnsupportedError. It's
// probably best to check for these and skip, because they are likely recoverable.
// Support errors are exported as expvars, so it is possibly not necessary to log
// them. If you want to check what's failed, start a debug http server and examine
// the vt100-unsupported-commands field in /debug/vars.
func (v *VT100) Process(c Command) error {
	return c.display(v)
}

// put puts r onto the current cursor's position, then advances the cursor.
func (v *VT100) put(r rune) {
	if v.Cursor.Y >= v.Height || v.Cursor.X >= v.Width {
		v.advance()
	}
	v.Content[v.Cursor.Y][v.Cursor.X] = r
	v.Format[v.Cursor.Y][v.Cursor.X] = v.Cursor.F
	v.advance()
}

// advance advances the cursor, wrapping to the next line if need be.
func (v *VT100) advance() {
	v.Cursor.X++
	if v.Cursor.X >= v.Width {
		v.Cursor.X = 0
		v.Cursor.Y++
	}
	if v.Cursor.Y >= v.Height {
		v.scrollDown()
	}
}

// scrollDown scrolls down the content one line
func (v *VT100) scrollDown() {
	for i := 0; i < len(v.Content)-1; i++ {
		for j := 0; j < v.Width; j++ {
			v.Content[i][j] = v.Content[i+1][j]
			v.Format[i][j] = v.Format[i+1][j]
		}
	}
	v.Cursor.Y--
	v.Cursor.X = 0
	v.eraseRegion(v.Cursor.Y, 0, v.Cursor.Y, v.Width-1)
}

// home moves the cursor to the coordinates y x. If y x are out of bounds, v.Err
// is set.
func (v *VT100) home(y, x int) {
	v.Cursor.Y, v.Cursor.X = y, x
}

// eraseDirection is the logical direction in which an erase command happens,
// from the cursor. For both erase commands, forward is 0, backward is 1,
// and everything is 2.
type eraseDirection int

const (
	// From the cursor to the end, inclusive.
	eraseForward eraseDirection = iota

	// From the beginning to the cursor, inclusive.
	eraseBack

	// Everything.
	eraseAll
)

// eraseColumns erases columns from the current line.
func (v *VT100) eraseColumns(d eraseDirection) {
	y, x := v.Cursor.Y, v.Cursor.X // Aliases for simplicity.
	switch d {
	case eraseBack:
		v.eraseRegion(y, 0, y, x)
	case eraseForward:
		v.eraseRegion(y, x, y, v.Width-1)
	case eraseAll:
		v.eraseRegion(y, 0, y, v.Width-1)
	}
}

// eraseLines erases lines from the current terminal. Note that
// no matter what is selected, the entire current line is erased.
func (v *VT100) eraseLines(d eraseDirection) {
	y := v.Cursor.Y // Alias for simplicity.
	switch d {
	case eraseBack:
		v.eraseRegion(0, 0, y, v.Width-1)
	case eraseForward:
		v.eraseRegion(y, 0, v.Height-1, v.Width-1)
	case eraseAll:
		v.eraseRegion(0, 0, v.Height-1, v.Width-1)
	}
}

func (v *VT100) eraseRegion(y1, x1, y2, x2 int) {
	// Do not sanitize or bounds-check these coordinates, since they come from the
	// programmer (me). We should panic if any of them are out of bounds.
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}

	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			v.clear(y, x)
		}
	}
}

func (v *VT100) clear(y, x int) {
	v.Content[y][x] = ' '
	v.Format[y][x] = Format{}
}

func (v *VT100) backspace() {
	v.Cursor.X--
	if v.Cursor.X < 0 {
		if v.Cursor.Y == 0 {
			v.Cursor.X = 0
		} else {
			v.Cursor.Y--
			v.Cursor.X = v.Width - 1
		}
	}
}

func (v *VT100) save() {
	v.savedCursor = v.Cursor
}

func (v *VT100) unsave() {
	v.Cursor = v.savedCursor
}
