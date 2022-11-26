package cntdialogs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v4/pkg/channel"
	"github.com/gdamore/tcell/v2"
	"github.com/hinshun/vt10x"
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
	layout             *tview.Flex
	cntInfo            *tview.InputField
	terminalScreen     *tview.Box
	form               *tview.Form
	cancelHandler      func()
	fastRefreshHandler func()
	display            bool
	state              state
	containerID        string
	sessionID          string
	streamDoneChan     chan bool
	streamOutputWriter channel.WriteCloser
	streamInputBuffer  *bufio.Writer
	vtWriter           *io.PipeWriter
	vtReader           *io.PipeReader
	vtReaderBuf        **bufio.Reader
	vtTerminal         vt10x.Terminal
	detachKeys         detachKeys
	focusElement       int
	ttyWidth           int
	ttyHeight          int
}

// NewContainerExecTerminalDialog returns new container exec terminal dialog
func NewContainerExecTerminalDialog() *ContainerExecTerminalDialog {
	dialog := &ContainerExecTerminalDialog{
		Box:     tview.NewBox(),
		cntInfo: tview.NewInputField(),
		state: state{
			isRunning: false,
		},
		terminalScreen: tview.NewBox(),
		streamDoneChan: make(chan bool, 2),
		display:        false,
		detachKeys: detachKeys{
			keyString: "ctrl-p,ctrl-q,ctrl-p",
			tcellKeys: []tcell.Key{
				tcell.KeyCtrlP, tcell.KeyCtrlQ, tcell.KeyCtrlP,
			},
		},
	}
	bgColor := style.DialogBgColor

	// label
	cntInfoLabel := "CONTAINER ID:"
	dialog.cntInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.cntInfo.SetLabel("[::b]" + cntInfoLabel)
	dialog.cntInfo.SetLabelWidth(len(cntInfoLabel) + 1)
	dialog.cntInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.cntInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// terminal screen
	terminalBgColor := style.TerminalBgColor
	terminalBorderColor := style.TerminalBorderColor
	dialog.terminalScreen.SetBackgroundColor(terminalBgColor)
	dialog.terminalScreen.SetBorder(true)
	dialog.terminalScreen.SetBorderColor(terminalBorderColor)

	// form fields
	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	//label and terminal layout
	middleLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	middleLayout.SetBackgroundColor(bgColor)
	middleLayout.SetBorder(false)
	//middleLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	middleLayout.AddItem(dialog.cntInfo, 1, 0, true)
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
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetTitle("PODMAN CONTAINER EXEC")

	// main layout
	dialog.layout.AddItem(layout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive
func (d *ContainerExecTerminalDialog) Display() {
	// clear terminal content
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
	d.sendDetach()
	d.state.stop()
	d.SetExecInfo("", "", "")
	d.ttyHeight = 0
	d.ttyWidth = 0
	d.display = false
	// close OS pipes attached to session and VT terminal
	d.vtWriter.Close()
	d.streamOutputWriter.Close()
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
				d.cancelHandler()
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
	dX := x + 1
	dY := y + 1
	dWidth := width - 2
	dHeight := height - 2

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

	terminalFgColor := style.TerminalFgColor
	terminalBgColor := style.TerminalBgColor
	terminalStyle := tcell.StyleDefault.Foreground(terminalFgColor).Background(terminalBgColor)

	x, y, width, height = d.terminalScreen.GetInnerRect()
	if width != d.ttyWidth || height != d.ttyHeight {
		d.setTtySize(width, height)
	}

	// set terminal background
	for trow := 0; trow < height; trow++ {
		for tcol := 0; tcol < width; tcol++ {
			tview.PrintJoinedSemigraphics(screen, x+tcol, y+trow, rune(0), terminalStyle)
		}
	}
	content, cursor := d.vtContent()
	contentLines := strings.Split(content, "\n")
	for row := 0; row < len(contentLines); row++ {
		tview.PrintSimple(screen, contentLines[row], x, y+row)
	}
	cursorX := x + cursor.X
	cursorY := y + cursor.Y
	if cursor.Y < height && cursor.X < width {
		tview.Print(screen, "▉", cursorX, cursorY, 1, tview.AlignCenter, terminalFgColor)
	}

}

// SetCancelFunc sets form close button selected function
func (d *ContainerExecTerminalDialog) SetCancelFunc(handler func()) *ContainerExecTerminalDialog {
	d.cancelHandler = handler
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

	d.vtTerminal.Resize(d.ttyWidth, d.ttyHeight)

	go containers.ResizeExecTty(d.sessionID, d.ttyHeight, d.ttyWidth)
}

// PrepareForExec prepare for session exec e.g. setting input, output streams
func (d *ContainerExecTerminalDialog) PrepareForExec(id string, name string, opts *containers.ExecOption) error {
	log.Debug().Msg("container exec terminal dialog: prepare for exec")
	d.vtTerminal = vt10x.New()
	d.containerID = id

	c := make(chan []byte, 1000)
	d.streamOutputWriter = channel.NewWriter(c)
	termInputReader, termInputWriter, err := os.Pipe()

	if err != nil {
		return err
	}
	execInputStream := bufio.NewReader(termInputReader)
	opts.OutputStream = d.streamOutputWriter
	opts.InputStream = execInputStream
	opts.DetachKeys = d.detachKeys.string()

	terminalInputBuffer := bufio.NewWriter(termInputWriter)
	d.streamInputBuffer = terminalInputBuffer

	d.setTtySize(opts.TtyWidth, opts.TtyHeight)
	return nil
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
		case data := <-d.streamOutputWriter.Chan():
			if d.state.isStopped() {
				log.Debug().Msg("container exec terminal dialog: reader stream stopped")
				return
			}
			d.addToOutput(string(data))
			d.fastRefreshHandler()
		}

	}
}

func (d *ContainerExecTerminalDialog) addToOutput(data string) {
	d.vtWriter.Write([]byte(data))
}

// vtContent returns current content of vt10x terminal
func (d *ContainerExecTerminalDialog) vtContent() (string, vt10x.Cursor) {
	var content string
	var cursor vt10x.Cursor
	content = d.vtTerminal.String()
	cursor = d.vtTerminal.Cursor()
	return content, cursor
}

func (d *ContainerExecTerminalDialog) startVTreader() {
	log.Debug().Msg("container exec terminal dialog: terminal reader started")
	for {
		err := d.vtTerminal.Parse(*d.vtReaderBuf)
		if err != nil {
			if err != io.EOF {
				log.Error().Msgf("container exec terminal dialog: terminal reader error %v", err)
			}
			break
		}
	}
	log.Debug().Msg("container exec terminal dialog: terminal reader exited")
	d.vtReader.Close()
}

// SetExecInfo sets container exec terminal information
// container ID , name and session ID
func (d *ContainerExecTerminalDialog) SetExecInfo(id string, name string, sessionID string) {
	d.containerID = id
	d.sessionID = sessionID
	if len(sessionID) > utils.IDLength {
		sessionID = sessionID[0:utils.IDLength]
	}

	screenLabel := fmt.Sprintf("SESSION (%s)", sessionID)
	d.terminalScreen.SetTitle(screenLabel)

	containerInfo := fmt.Sprintf("%12s (%s)", id, name)
	d.cntInfo.SetText(containerInfo)
}

func (d *ContainerExecTerminalDialog) sendDetach() {
	log.Debug().Msg("container exec terminal dialog: sending detached keys")
	keys := d.detachKeys.keys()
	for i := 0; i < len(keys); i++ {
		d.streamInputBuffer.WriteRune(rune(keys[i]))
	}
	d.streamInputBuffer.Flush()
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

type detachKeys struct {
	keyString string
	tcellKeys []tcell.Key
}

func (dkey detachKeys) string() string {
	return dkey.keyString
}

func (dkey detachKeys) keys() []tcell.Key {
	return dkey.tcellKeys
}
