package vterm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/channel"
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

const vTermDialogLabelPadding = 1

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
	alreadyDetached       bool
	alreadyDetachedLock   sync.Mutex
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
		display:         false,
		alreadyDetached: false,
	}

	dialog.initLayoutUI()

	return dialog
}

// InitChannels will init buffers and channels for attach.
func (d *VtermDialog) InitAttachChannels() (io.Reader, io.Writer) {
	log.Debug().Msg("view: container terminal dialog init channels (attach)")

	d.sessionMode = sessionModeAttach

	d.initChannelsCommon()

	c := make(chan []byte, 1000) //nolint:mnd
	d.sessionStdout = NewWriter(c)

	d.init = true

	return d.sessionStdin, d.sessionStdout
}

// InitChannels will init buffers and channels for exec.
func (d *VtermDialog) InitExecChannels() (*bufio.Reader, channel.WriteCloser) { //nolint:ireturn
	log.Debug().Msg("view: container terminal dialog init channels (exec)")

	d.sessionMode = sessionModeExec

	d.initChannelsCommon()

	c := make(chan []byte, 1000) //nolint:mnd

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

	switch d.sessionMode {
	case sessionModeAttach:
		go d.sessionOutputStreamer()
	case sessionModeExec:
		go d.execSessionOutputStreamer()
	}

	if d.sessionMode != sessionModeNone {
		go d.startVTBuffer()
	}

	d.SetAlreadyDetach(false)
	d.display = true
}

func (d *VtermDialog) SetAlreadyDetach(detached bool) {
	d.alreadyDetachedLock.Lock()
	defer d.alreadyDetachedLock.Unlock()

	d.alreadyDetached = detached
}

func (d *VtermDialog) IsAlreadyDetach() bool {
	var alreadyDetached bool

	d.alreadyDetachedLock.Lock()

	alreadyDetached = d.alreadyDetached

	d.alreadyDetachedLock.Unlock()

	return alreadyDetached
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

	if !d.IsAlreadyDetach() {
		d.sendDetachToSession()
	}

	if d.sessionMode == sessionModeExec {
		err := d.execSessionStdout.Close()
		if err != nil {
			log.Error().Msgf("failed to close vterm exec stdout session: %s", err.Error())
		}
	}

	err := d.vtTermPipeReader.Close()
	if err != nil {
		log.Error().Msgf("failed to close vterm pipe reader: %s", err.Error())
	}

	err = d.vtTermPipeWriter.Close()
	if err != nil {
		log.Error().Msgf("failed to close vterm pipe writer: %s", err.Error())
	}
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

				if !d.IsAlreadyDetach() {
					d.writeToSession(event)
				}

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
	dWidth := width - 2   //nolint:mnd
	dHeight := height - 2 //nolint:mnd

	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// Draw draws this primitive onto the screen.
func (d *VtermDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()
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
	for trow := range height {
		for tcol := range width {
			tview.PrintJoinedSemigraphics(screen, x+tcol, y+trow, rune(0), terminalStyle)
		}
	}

	content, cursor := d.vtContent()

	contentLines := strings.Split(content, "\n")
	for row := range contentLines {
		tview.PrintSimple(screen, contentLines[row], x, y+row)
	}

	cursorX := x + cursor.X
	cursorY := y + cursor.Y

	if cursor.Y < height && cursor.X < width {
		tview.Print(screen, "▉", cursorX, cursorY, 1, tview.AlignCenter, terminalFgColor)
	}
}

// SetCancelFunc sets form close button selected function.
func (d *VtermDialog) SetCancelFunc(handler func()) *VtermDialog {
	d.cancelHandler = handler
	// closeButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	// closeButton.SetSelectedFunc(handler)

	return d
}

// SetContainerInfo sets container's ID and NAME to the terminal header.
func (d *VtermDialog) SetContainerInfo(id string, name string) {
	d.containerID = id
	containerInfo := id[0:12]

	if name != "" {
		containerInfo = fmt.Sprintf("%s (%s)", containerInfo, name)
		containerInfo = utils.LabelWidthLeftPadding(containerInfo, vTermDialogLabelPadding)
	}

	d.containerInfo.SetText(containerInfo)
}

// SetSessionID sets attach sessions's ID.
func (d *VtermDialog) SetSessionID(id string) {
	d.sessionID = id

	if len(id) > utils.IDLength {
		id = id[0:utils.IDLength]
	}

	sessionIDLabel := fmt.Sprintf("TERMINAL SESSION (%s)", id)
	d.termScreen.SetTitle(sessionIDLabel)
}

// DetachKeys returns the detach keys used to detach from the session.
func (d *VtermDialog) DetachKeys() string {
	return d.detachKeys.string()
}

// SetFastRefreshHandler sets fast refresh handler
// fast refresh is used to print the outputs as fast as possible.
func (d *VtermDialog) SetFastRefreshHandler(handler func()) {
	d.fastRefreshHandler = handler
}

func (d *VtermDialog) writeToSession(event *tcell.EventKey) {
	switch event.Key() { //nolint:exhaustive
	case tcell.KeyUp:
		d.writeToStdinSessionWriter(rune(27)) //nolint:mnd
		d.writeToStdinSessionWriter(rune(91)) //nolint:mnd
		d.writeToStdinSessionWriter(rune(65)) //nolint:mnd
	case tcell.KeyDown:
		d.writeToStdinSessionWriter(rune(27)) //nolint:mnd
		d.writeToStdinSessionWriter(rune(91)) //nolint:mnd
		d.writeToStdinSessionWriter(rune(66)) //nolint:mnd
	case tcell.KeyRight:
		d.writeToStdinSessionWriter(rune(27)) //nolint:mnd
		d.writeToStdinSessionWriter(rune(91)) //nolint:mnd
		d.writeToStdinSessionWriter(rune(67)) //nolint:mnd
	case tcell.KeyLeft:
		d.writeToStdinSessionWriter(rune(27)) //nolint:mnd
		d.writeToStdinSessionWriter(rune(91)) //nolint:mnd
		d.writeToStdinSessionWriter(rune(68)) //nolint:mnd
	case tcell.KeyEsc:
		d.writeToStdinSessionWriter(rune(27)) //nolint:mnd
	default:
		d.writeToStdinSessionWriter(event.Rune())
	}

	err := d.sessionStdinWriter.Flush()
	if err != nil {
		log.Error().Msgf("failed to flush vterm stdin writer session: %s", err.Error())
	}
}

func (d *VtermDialog) sendDetachToSession() {
	log.Debug().Msg("view: container terminal dialog sending detached keys")

	keys := d.detachKeys.keys()

	for i := range keys {
		d.writeToStdinSessionWriter(rune(keys[i]))
	}

	err := d.sessionStdinWriter.Flush()
	if err != nil {
		log.Error().Msgf("failed to flush vterm stdin writer session: %s", err.Error())
	}
}

// vtContent returns current content and cursor location of vt10x terminal.
func (d *VtermDialog) vtContent() (string, vt10x.Cursor) {
	var (
		content string
		cursor  vt10x.Cursor
	)

	content = d.vtTerminal.String()
	cursor = d.vtTerminal.Cursor()

	return content, cursor
}

func (d *VtermDialog) startVTBuffer() {
	log.Debug().Msg("view: vterm dialog vt buffer reader started")

	for {
		err := d.vtTerminal.Parse(d.vtTermBuffer)
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) {
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
			_, err := d.vtTermPipeWriter.Write(data)
			if err != nil {
				log.Error().Msgf("failed to write %s to vterm pipe writer: %s", data, err.Error())
			}

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

			_, err := d.vtTermPipeWriter.Write([]byte(dataString))
			if err != nil {
				log.Error().Msgf("failed to write %s to vterm pipe writer: %s", []byte(dataString), err.Error())
			}

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

func (d *VtermDialog) initLayoutUI() {
	bgColor := style.DialogBgColor
	borderColor := style.DialogBorderColor
	terminalBgColor := style.TerminalBgColor
	terminalBorderColor := style.TerminalBorderColor

	// container information field
	// label
	d.containerInfo.SetBackgroundColor(bgColor)
	d.containerInfo.SetLabel("[::b]" + utils.ContainerIDLabel)
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
	// layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
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
	// d.layout.SetTitle("CONTAINER TERMINAL")

	d.layout.AddItem(termLayout, 0, 1, true)
	d.layout.AddItem(d.form, dialogs.DialogFormHeight, 0, true)
}

func (d *VtermDialog) initChannelsCommon() {
	d.sessionOutputDoneChan = make(chan bool, 2) //nolint:mnd
	d.vtTerminal = vt10x.New()
	sessionStdinPipeIn, sessionStdinPipeOut := io.Pipe()
	d.sessionStdin = bufio.NewReader(sessionStdinPipeIn)
	d.sessionStdinWriter = bufio.NewWriter(sessionStdinPipeOut)

	d.vtTermPipeReader, d.vtTermPipeWriter = io.Pipe()
	d.vtTermBuffer = bufio.NewReader(d.vtTermPipeReader)
}

func (d *VtermDialog) setTTYSize(width int, height int) {
	if width < 0 || height < 0 {
		return
	}

	d.ttyWidth = width
	d.ttyHeight = height

	d.vtTerminal.Resize(width, height)

	switch d.sessionMode {
	case sessionModeAttach:
		go containers.ResizeContainerTTY(d.containerID, width, height) //nolint:errcheck
	case sessionModeExec:
		go containers.ResizeExecTty(d.sessionID, d.ttyHeight, d.ttyWidth)
	}
}

func (d *VtermDialog) writeToStdinSessionWriter(val rune) {
	_, err := d.sessionStdinWriter.WriteRune(val)
	if err != nil {
		log.Error().Msgf("failed to write value %d to vterm session writer rune: %s", val, err.Error())
	}
}
