package vterm

import (
	"bufio"
	"fmt"
	"io"
	"strings"

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
	vtermDialogScreenFieldFocus = 0 + iota
	vtermDialogFormFieldFocus
)

const (
	sessionModeNone = 0 + iota
	sessionModeAttach
	sessionModeExec
)

// VtermDialog implements virtual terminal that can be used during
// exec, attach, run activity.
type VtermDialog struct {
	*tview.Box
	layout                *tview.Flex
	form                  *tview.Form
	containerInfo         *tview.InputField
	termScreen            *tview.Box
	detachKeys            termDetachKeys
	sessionMode           int
	sessionStdin          *bufio.Reader
	sessionStdout         Writer
	execSessionStdout     channel.WriteCloser
	sessionStdinWriter    *bufio.Writer
	vtTerminal            vt10x.Terminal
	vtTermBuffer          *bufio.Reader
	vtTermPipeWriter      *io.PipeWriter
	vtTermPipeReader      *io.PipeReader
	display               bool
	sessionOutputDoneChan chan bool
	containerID           string
	sessionID             string
	focusElement          int
	init                  bool
	ttyWidth              int
	ttyHeight             int
	cancelHandler         func()
	fastRefreshHandler    func()
}

// NewVtermDialog returns new VtermDialog primitive.
func NewVtermDialog() *VtermDialog {
	dialog := &VtermDialog{
		Box:                   tview.NewBox(),
		layout:                tview.NewFlex(),
		form:                  tview.NewForm(),
		containerInfo:         tview.NewInputField(),
		termScreen:            tview.NewBox(),
		sessionOutputDoneChan: make(chan bool),
		detachKeys: termDetachKeys{
			keyString: "ctrl-p,ctrl-q,ctrl-p",
			tcellKeys: []tcell.Key{
				tcell.KeyCtrlP, tcell.KeyCtrlQ, tcell.KeyCtrlP,
			},
		},
		display: false,
	}

	dialog.initLayoutUI()
	return dialog
}

func (d *VtermDialog) initLayoutUI() {
	bgColor := style.DialogBgColor
	borderColor := style.DialogBorderColor
	terminalBgColor := style.TerminalBgColor
	terminalBorderColor := style.TerminalBorderColor

	// container information field
	// label
	cntInfoLabel := "CONTAINER ID:"
	d.containerInfo.SetBackgroundColor(bgColor)
	d.containerInfo.SetLabel("[::b]" + cntInfoLabel)
	d.containerInfo.SetLabelWidth(len(cntInfoLabel) + 1)
	d.containerInfo.SetFieldBackgroundColor(bgColor)
	d.containerInfo.SetLabelStyle(tcell.StyleDefault.
		Background(borderColor).
		Foreground(style.DialogFgColor))

	// form fields
	d.form.AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)
	d.form.SetBackgroundColor(bgColor)
	d.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// terminal screen
	d.termScreen.SetBackgroundColor(terminalBgColor)
	d.termScreen.SetBorder(true)
	d.termScreen.SetBorderColor(terminalBorderColor)

	// layout setup
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	termLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	layout.SetBackgroundColor(bgColor)
	layout.SetBorder(false)
	layout.AddItem(d.containerInfo, 1, 0, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	layout.AddItem(d.termScreen, 0, 1, true)

	termLayout.SetBackgroundColor(bgColor)
	termLayout.SetBorder(false)
	termLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	termLayout.AddItem(layout, 0, 1, false)
	termLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	// main layout
	d.layout.SetDirection(tview.FlexRow)
	d.layout.SetBorder(true)
	d.layout.SetBorderColor(borderColor)
	d.layout.SetBackgroundColor(bgColor)
	d.layout.SetTitle("CONTAINER TERMINAL")

	d.layout.AddItem(termLayout, 0, 1, true)
	d.layout.AddItem(d.form, dialogs.DialogFormHeight, 0, true)
}

func (d *VtermDialog) initChannelsCommon() {
	d.sessionOutputDoneChan = make(chan bool, 2)
	d.vtTerminal = vt10x.New()
	sessionStdinPipeIn, sessionStdinPipeOut := io.Pipe()
	d.sessionStdin = bufio.NewReader(sessionStdinPipeIn)
	d.sessionStdinWriter = bufio.NewWriter(sessionStdinPipeOut)

	d.vtTermPipeReader, d.vtTermPipeWriter = io.Pipe()
	d.vtTermBuffer = bufio.NewReader(d.vtTermPipeReader)
}

// InitChannels will init buffers and channels for attach.
func (d *VtermDialog) InitAttachChannels() (io.Reader, io.Writer) {
	log.Debug().Msg("view: container terminal dialog init channels (attach)")

	d.sessionMode = sessionModeAttach
	d.initChannelsCommon()
	c := make(chan []byte, 1000)
	d.sessionStdout = NewWriter(c)

	d.init = true

	return d.sessionStdin, d.sessionStdout
}

// InitChannels will init buffers and channels for exec.
func (d *VtermDialog) InitExecChannels() (*bufio.Reader, channel.WriteCloser) {
	log.Debug().Msg("view: container terminal dialog init channels (exec)")
	d.sessionMode = sessionModeExec

	d.initChannelsCommon()
	c := make(chan []byte, 1000)
	d.execSessionStdout = channel.NewWriter(c)

	d.init = true
	return d.sessionStdin, d.execSessionStdout
}

// Display will start display this primitive onto the screen.
func (d *VtermDialog) Display() {
	if !d.init {
		log.Error().Msg("view: container terminal dialog is not read, init first")
		return
	}

	if d.sessionMode == sessionModeAttach {
		go d.sessionOutputStreamer()
	} else if d.sessionMode == sessionModeExec {
		go d.execSessionOutputStreamer()
	}

	if d.sessionMode != sessionModeNone {
		go d.startVTBuffer()
	}

	d.display = true
}

// IsDisplay returns true if primitive is shown onto the screen.
func (d *VtermDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive onto the screen.
func (d *VtermDialog) Hide() {
	d.display = false
	d.containerID = ""
	d.sessionID = ""
	d.termScreen.SetTitle("")
	d.focusElement = vtermDialogScreenFieldFocus

	d.sessionOutputDoneChan <- true

	d.sessionMode = sessionModeNone

	d.sendDetachToSession()

	if d.sessionMode == sessionModeExec {
		d.execSessionStdout.Close()
	}

	d.vtTermPipeReader.Close()
	d.vtTermPipeWriter.Close()
}

// HasFocus returns true if terminal dialog has focus.
func (d *VtermDialog) HasFocus() bool {
	if d.form.HasFocus() || d.containerInfo.HasFocus() {
		return true
	}

	if d.termScreen.HasFocus() || d.layout.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *VtermDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	// terminal screen field focus
	case vtermDialogScreenFieldFocus:
		d.termScreen.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = vtermDialogFormFieldFocus

				d.Focus(delegate)
				return nil
			}

			return event
		})
		delegate(d.termScreen)
	// form field focus
	case vtermDialogFormFieldFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = vtermDialogScreenFieldFocus
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

// InputHandler returns input handler function for this primitive.
func (d *VtermDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: container terminal dialog event %v received", event)

		if d.termScreen.HasFocus() {
			if handler := d.termScreen.InputHandler(); handler != nil {
				handler(event, setFocus)

				d.writeToSession(event)
				return
			}
		}

		if d.form.HasFocus() {
			if handler := d.form.InputHandler(); handler != nil {
				handler(event, setFocus)

				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *VtermDialog) SetRect(x, y, width, height int) {
	dX := x + 1
	dY := y + 1
	dWidth := width - 2
	dHeight := height - 2

	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// Draw draws this primitive onto the screen.
func (d *VtermDialog) Draw(screen tcell.Screen) {
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

	x, y, width, height = d.termScreen.GetInnerRect()
	if width != d.ttyWidth || height != d.ttyHeight {
		d.setTTYSize(width, height)
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

func (d *VtermDialog) setTTYSize(width int, height int) {
	if width < 0 || height < 0 {
		return
	}

	d.ttyWidth = width
	d.ttyHeight = height

	d.vtTerminal.Resize(width, height)

	if d.sessionMode == sessionModeAttach {
		go containers.ResizeContainerTTY(d.containerID, width, height)
	} else if d.sessionMode == sessionModeExec {
		go containers.ResizeExecTty(d.sessionID, d.ttyHeight, d.ttyWidth)
	}
}

// SetCancelFunc sets form close button selected function
func (d *VtermDialog) SetCancelFunc(handler func()) *VtermDialog {
	d.cancelHandler = handler
	// closeButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	// closeButton.SetSelectedFunc(handler)
	return d
}

// SetContainerInfo sets container's ID and NAME to the terminal header.
func (d *VtermDialog) SetContainerInfo(id string, name string) {
	d.containerID = id

	containerInfo := fmt.Sprintf("%12s (%s)", id, name)
	d.containerInfo.SetText(containerInfo)
}

// SetSessionID sets attach sessions's ID.
func (d *VtermDialog) SetSessionID(id string) {
	d.sessionID = id

	if len(id) > utils.IDLength {
		id = id[0:utils.IDLength]
	}

	sessionIDLabel := fmt.Sprintf("SESSION (%s)", id)
	d.termScreen.SetTitle(sessionIDLabel)
}

// DetachKeys returns the detach keys used to detach from the session.
func (d *VtermDialog) DetachKeys() string {
	return d.detachKeys.string()
}

// SetFastRefreshHandler sets fast refresh handler
// fast refresh is used to print the outputs as fast as possible
func (d *VtermDialog) SetFastRefreshHandler(handler func()) {
	d.fastRefreshHandler = handler
}

func (d *VtermDialog) writeToSession(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyUp:
		d.sessionStdinWriter.WriteRune(rune(27))
		d.sessionStdinWriter.WriteRune(rune(91))
		d.sessionStdinWriter.WriteRune(rune(65))
	case tcell.KeyDown:
		d.sessionStdinWriter.WriteRune(rune(27))
		d.sessionStdinWriter.WriteRune(rune(91))
		d.sessionStdinWriter.WriteRune(rune(66))
	case tcell.KeyRight:
		d.sessionStdinWriter.WriteRune(rune(27))
		d.sessionStdinWriter.WriteRune(rune(91))
		d.sessionStdinWriter.WriteRune(rune(67))
	case tcell.KeyLeft:
		d.sessionStdinWriter.WriteRune(rune(27))
		d.sessionStdinWriter.WriteRune(rune(91))
		d.sessionStdinWriter.WriteRune(rune(68))
	case tcell.KeyEsc:
		d.sessionStdinWriter.WriteRune(rune(27))
	default:
		d.sessionStdinWriter.WriteRune(event.Rune())
	}

	d.sessionStdinWriter.Flush()
}

func (d *VtermDialog) sendDetachToSession() {
	log.Debug().Msg("view: container terminal dialog sending detached keys")

	keys := d.detachKeys.keys()
	for i := 0; i < len(keys); i++ {
		d.sessionStdinWriter.WriteRune(rune(keys[i]))
	}

	d.sessionStdinWriter.Flush()
}

// vtContent returns current content and cursor location of vt10x terminal.
func (d *VtermDialog) vtContent() (string, vt10x.Cursor) {
	var content string
	var cursor vt10x.Cursor

	content = d.vtTerminal.String()
	cursor = d.vtTerminal.Cursor()

	return content, cursor
}

func (d *VtermDialog) startVTBuffer() {
	log.Debug().Msg("view: vterm dialog vt buffer reader started")
	for {
		err := d.vtTerminal.Parse(d.vtTermBuffer)
		if err != nil {
			if err != io.EOF && err != io.ErrClosedPipe {
				log.Error().Msgf("view: vterm dialog vt buffer reader error %v", err)
			}

			break
		}
	}

	log.Debug().Msg("view: vterm dialog vt buffer reader exited")
}

func (d *VtermDialog) execSessionOutputStreamer() {
	log.Debug().Msgf("view: vterm dialog exec session output streamer started")

	for {
		select {
		case <-d.sessionOutputDoneChan:
			log.Debug().Msgf("view: vterm dialog exec session output streamer exited (stop signal)")

			return
		case data := <-d.execSessionStdout.Chan():
			d.vtTermPipeWriter.Write(data)
			d.fastRefreshHandler()
		}
	}
}

func (d *VtermDialog) sessionOutputStreamer() {
	log.Debug().Msgf("view: vterm dialog session output streamer started")

	for {
		select {
		case <-d.sessionOutputDoneChan:
			log.Debug().Msgf("view: vterm dialog session output streamer exited (stop signal)")

			return
		case data := <-d.sessionStdout.Chan():
			dataString := string(data)
			dataString = strings.ReplaceAll(dataString, "\n", "\r\n")

			d.vtTermPipeWriter.Write([]byte(dataString))
			d.fastRefreshHandler()
		}
	}
}

type termDetachKeys struct {
	keyString string
	tcellKeys []tcell.Key
}

func (dkey termDetachKeys) string() string {
	return dkey.keyString
}

func (dkey termDetachKeys) keys() []tcell.Key {
	return dkey.tcellKeys
}
