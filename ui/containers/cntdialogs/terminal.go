package cntdialogs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	vt100 "github.com/navidys/vtterm"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	defaultTtyHeight = 20
	defaultTtyWidth  = 100
)

const (
	terminalScreenFieldFocus = 0 + iota
	terminalFormFieldFocus
)

// ContainerExecTerminalDialog represents container execterminal dialog primitive
type ContainerExecTerminalDialog struct {
	*tview.Box
	layout              *tview.Flex
	label               *tview.TextView
	terminalScreen      *tview.Box
	form                *tview.Form
	doneHandler         func()
	fastRefreshHandler  func()
	display             bool
	state               state
	containerID         string
	sessionID           string
	streamDoneChan      chan bool
	streamOutputChannel *chan string
	streamInputBuffer   *bufio.Writer
	vtWriter            *io.PipeWriter
	vtReader            *io.PipeReader
	vtReaderBuf         **bufio.Reader
	vtTerminal          vt
	focusElement        int
	ttyWidth            int
	ttyHeight           int
}

type vt struct {
	*vt100.VT100
	sync.Locker
}

// NewContainerExecTerminalDialog returns new container exec terminal dialog
func NewContainerExecTerminalDialog() *ContainerExecTerminalDialog {
	dialog := &ContainerExecTerminalDialog{
		Box:   tview.NewBox(),
		label: tview.NewTextView(),
		state: state{
			isRunning: false,
		},
		terminalScreen: tview.NewBox(),
		streamDoneChan: make(chan bool, 2),
		display:        false,
	}

	bgColor := utils.Styles.ContainerExecTerminalDialog.BgColor
	fgColor := utils.Styles.ContainerExecTerminalDialog.FgColor

	// label
	dialog.label.SetDynamicColors(true)
	dialog.label.SetBackgroundColor(bgColor)
	dialog.label.SetTextColor(fgColor)
	dialog.label.SetBorder(false)

	// terminal screen
	terminalBgColor := utils.Styles.ContainerExecTerminalDialog.Terminal.BgColor
	dialog.terminalScreen.SetBackgroundColor(terminalBgColor)
	dialog.terminalScreen.SetBorder(true)

	// form fields
	dialog.form = tview.NewForm().
		AddButton("Close", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(utils.Styles.ButtonPrimitive.BgColor)

	//label and terminal layout
	middleLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	middleLayout.SetBackgroundColor(bgColor)
	middleLayout.SetBorder(false)
	middleLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	middleLayout.AddItem(dialog.label, 1, 0, true)
	middleLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	middleLayout.AddItem(dialog.terminalScreen, 0, 1, true)

	// main dialog layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.SetBackgroundColor(bgColor)
	layout.SetBorder(false)

	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(middleLayout, 0, 1, false)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetTitle("PODMAN CONTAINER EXEC")

	// main layout
	dialog.layout.AddItem(layout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive
func (d *ContainerExecTerminalDialog) Display() {
	// OS pipes attached to session and VT terminal
	// it will read from exec session and write to VT termianl for decode
	vReaderRaw, vWriter := io.Pipe()
	vReader := bufio.NewReader(vReaderRaw)
	d.vtWriter = vWriter
	d.vtReader = vReaderRaw
	d.vtReaderBuf = &vReader

	d.focusElement = terminalScreenFieldFocus
	d.state.start()
	go d.startVTreader()
	go d.startTerminalOutputStreamReader()
	d.display = true
	d.fastRefreshHandler()
}

// IsDisplay returns true if primitive is shown
func (d *ContainerExecTerminalDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ContainerExecTerminalDialog) Hide() {
	d.state.stop()
	d.streamDoneChan <- true
	d.setExecInfo("", "", "")
	d.ttyHeight = 0
	d.ttyWidth = 0
	d.display = false
	// close OS pipes attached to session and VT terminal
	d.vtWriter.Close()

}

// HasFocus returns whether or not this primitive has focus
func (d *ContainerExecTerminalDialog) HasFocus() bool {
	if d.terminalScreen.HasFocus() || d.form.HasFocus() {
		return true
	}
	return d.Box.HasFocus() || d.layout.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ContainerExecTerminalDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	// terminal screen field focus
	case terminalScreenFieldFocus:
		d.terminalScreen.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = terminalFormFieldFocus
				d.Focus(delegate)
				return nil
			}
			d.addKeyToTerminal(event)
			return event
		})
		delegate(d.terminalScreen)
	// form field focus
	case terminalFormFieldFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = terminalScreenFieldFocus
				d.Focus(delegate)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.doneHandler()
				return nil
			}
			return event
		})
		delegate(d.form)
	}

}

// InputHandler returns input handler function for this primitive
func (d *ContainerExecTerminalDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container exec terminal dialog: event %v received", event)
		// terminal screen field
		if d.terminalScreen.HasFocus() {
			if terminalScreenHandler := d.terminalScreen.InputHandler(); terminalScreenHandler != nil {
				terminalScreenHandler(event, setFocus)
				return
			}
		}
		// form primitive
		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *ContainerExecTerminalDialog) SetRect(x, y, width, height int) {
	dX := x + dialogs.DialogPadding
	dY := y + dialogs.DialogPadding - 1
	dWidth := width - (2 * dialogs.DialogPadding)
	dHeight := height - (2 * (dialogs.DialogPadding - 1))

	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// Draw draws this primitive onto the screen.
func (d *ContainerExecTerminalDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)

	terminalFgColor := utils.Styles.ContainerExecTerminalDialog.Terminal.FgColor
	terminalBgColor := utils.Styles.ContainerExecTerminalDialog.Terminal.BgColor
	terminalStyle := tcell.StyleDefault.Foreground(terminalFgColor).Background(terminalBgColor)

	x, y, width, height = d.terminalScreen.GetInnerRect()

	// set terminal background
	for trow := 0; trow < height; trow++ {
		for tcol := 0; tcol < width; tcol++ {
			tview.PrintJoinedSemigraphics(screen, x+tcol, y+trow, rune(0), terminalStyle)
		}
	}

	content, cursor := d.vtContent()
	for row := 0; row < len(content); row++ {
		for col := 0; col < len(content[row]); col++ {
			if (col >= width) || (row >= height) {
				continue
			}
			screen.SetContent(x+col, y+row, content[row][col], nil, terminalStyle)
		}
	}
	cursorX := x + cursor.X
	cursorY := y + cursor.Y
	//cursorRune := []rune("▉")
	if cursor.Y < height && cursor.X < width {
		tview.Print(screen, "▉", cursorX, cursorY, 1, tview.AlignCenter, terminalFgColor)
	}

	//screen.SetContent(x+cursorX, y+cursorY, rune(0), cursorRune, terminalStyle)
}

// SetDoneFunc sets form close button selected function
func (d *ContainerExecTerminalDialog) SetDoneFunc(handler func()) *ContainerExecTerminalDialog {
	d.doneHandler = handler
	closeButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	closeButton.SetSelectedFunc(handler)
	return d
}

// SetFastRefreshHandler sets fast refresh handler
// fast refresh is used to print exec output as fast as possible
func (d *ContainerExecTerminalDialog) SetFastRefreshHandler(handler func()) {
	d.fastRefreshHandler = handler
}

func (d *ContainerExecTerminalDialog) setTtySize(width int, height int) {
	if width <= 0 {
		d.ttyWidth = defaultTtyWidth
	}
	if height <= 0 {
		d.ttyHeight = defaultTtyHeight
	}
	d.ttyWidth = width
	d.ttyHeight = height

	d.vtTerminal = vt{
		vt100.NewVT100(d.ttyHeight, d.ttyWidth),
		new(sync.Mutex),
	}
	go containers.ResizeExecTty(d.sessionID, d.ttyHeight, d.ttyWidth)
}

// PrepareForExec prepare for session exec e.g. setting input, output streams
func (d *ContainerExecTerminalDialog) PrepareForExec(id string, name string, opts *containers.ExecOption) (string, error) {
	d.containerID = id

	outputStream := utils.NewStreamChannel(100)
	termInputReader, termInputWriter, err := os.Pipe()

	if err != nil {
		return "", err
	}
	var execOutputStream io.WriteCloser = outputStream
	execInputStream := bufio.NewReader(termInputReader)
	opts.OutputStream = execOutputStream
	opts.InputStream = execInputStream
	terminalInputBuffer := bufio.NewWriter(termInputWriter)

	execSessionID, err := containers.NewExecSession(d.containerID, *opts)
	if err != nil {
		return "", err
	}
	d.setExecInfo(d.containerID, name, execSessionID)

	d.streamInputBuffer = terminalInputBuffer
	d.streamOutputChannel = outputStream.Channel()

	d.setTtySize(opts.TtyWidth, opts.TtyHeight)
	return execSessionID, nil
}

// addKeyToTerminal adds user input key to terminal output
func (d *ContainerExecTerminalDialog) addKeyToTerminal(event *tcell.EventKey) {

	switch event.Key() {
	case tcell.KeyUp:
		d.streamInputBuffer.WriteRune(rune(27))
		d.streamInputBuffer.WriteRune(rune(91))
		d.streamInputBuffer.WriteRune(rune(65))
	case tcell.KeyDown:
		d.streamInputBuffer.WriteRune(rune(27))
		d.streamInputBuffer.WriteRune(rune(91))
		d.streamInputBuffer.WriteRune(rune(66))
	case tcell.KeyRight:
		d.streamInputBuffer.WriteRune(rune(27))
		d.streamInputBuffer.WriteRune(rune(91))
		d.streamInputBuffer.WriteRune(rune(67))
	case tcell.KeyLeft:
		d.streamInputBuffer.WriteRune(rune(27))
		d.streamInputBuffer.WriteRune(rune(91))
		d.streamInputBuffer.WriteRune(rune(68))
	case tcell.KeyEsc:
		d.streamInputBuffer.WriteRune(rune(27))
	default:
		d.streamInputBuffer.WriteRune(event.Rune())
	}

	d.streamInputBuffer.Flush()
}

// startTerminalOutputStreamReader reads outputs from container session
// and adds to terminal output content
func (d *ContainerExecTerminalDialog) startTerminalOutputStreamReader() {
	log.Debug().Msg("container exec terminal dialog: reader stream started")
	for {
		select {
		case <-d.streamDoneChan:
			log.Debug().Msg("container exec terminal dialog: reader stream stopped")
			return
		case data := <-*d.streamOutputChannel:
			if d.state.isStopped() {
				break
			}
			d.addToOutput(data)
			d.fastRefreshHandler()
		}

	}
}

func (d *ContainerExecTerminalDialog) addToOutput(data string) {
	d.vtWriter.Write([]byte(data))
}

func (d *ContainerExecTerminalDialog) vtContent() ([][]rune, vt100.Cursor) {
	var content [][]rune
	var cursor vt100.Cursor
	d.vtTerminal.Lock()
	content = d.vtTerminal.Content
	cursor = d.vtTerminal.Cursor
	d.vtTerminal.Unlock()
	return content, cursor
}

func (d *ContainerExecTerminalDialog) startVTreader() {
	log.Debug().Msg("container exec terminal dialog: vt100 terminal reader started")
	for {
		cmd, err := vt100.Decode(*d.vtReaderBuf)
		if err == nil {
			d.vtTerminal.Lock()
			err = d.vtTerminal.Process(cmd)
			d.vtTerminal.Unlock()
		}
		if err == nil {
			continue
		}
		if err != io.EOF {
			// TODO: fix out of bound error
			// both out of bound error and unsupported controls can be safety ignored
			//log.Error().Msgf("container exec terminal dialog: vt100: %v", err)
			continue
		}
		log.Debug().Msg("container exec terminal dialog: vt100 terminal reader exited")
		d.vtReader.Close()
		return
	}
}

func (d *ContainerExecTerminalDialog) setExecInfo(id string, name string, sessionID string) {
	fgColor := utils.Styles.ContainerExecTerminalDialog.FgColor
	bgColor := utils.Styles.ContainerExecTerminalDialog.HeaderBgColor
	headerFgColor := utils.GetColorName(fgColor)
	headerBgColor := utils.GetColorName(bgColor)
	d.containerID = id
	d.sessionID = sessionID
	if len(sessionID) > utils.IDLength {
		sessionID = sessionID[0:utils.IDLength]
	}

	label := fmt.Sprintf("[%s:%s:] CONTAINER ID: [%s:-:] %s ", headerFgColor, headerBgColor, headerFgColor, id)
	if name != "" {
		label = fmt.Sprintf("%s(%s) ", label, name)
	}
	screenLabel := fmt.Sprintf("SESSION (%s)", sessionID)
	d.terminalScreen.SetTitle(screenLabel)
	d.label.SetText(label)
}

type state struct {
	isRunning bool
	mu        sync.Mutex
}

func (s *state) isStopped() bool {
	var stopped bool
	s.mu.Lock()
	stopped = !s.isRunning
	s.mu.Unlock()
	return stopped
}

func (s *state) stop() {
	s.mu.Lock()
	s.isRunning = false
	s.mu.Unlock()
}

func (s *state) start() {
	s.mu.Lock()
	s.isRunning = true
	s.mu.Unlock()
}
