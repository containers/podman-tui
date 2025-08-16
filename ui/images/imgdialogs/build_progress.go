package imgdialogs

import (
	"sync"
	"time"

	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/channel"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	buildPrgDialogMaxWidth = 100
	buildPrgDialogHeight   = 20
)

// ImageBuildProgressDialog implements build progress dialog primitive.
type ImageBuildProgressDialog struct {
	*tview.Box

	layout             *tview.Flex
	output             *tview.TextView
	progressBar        *tvxwidgets.ActivityModeGauge
	display            bool
	cancelChan         chan bool
	writerChan         chan []byte
	mu                 sync.Mutex
	fastRefreshHandler func()
}

// NewImageBuildProgressDialog returns new build progress dialog.
func NewImageBuildProgressDialog() *ImageBuildProgressDialog {
	buildPrgDialog := &ImageBuildProgressDialog{
		Box:         tview.NewBox(),
		layout:      tview.NewFlex().SetDirection(tview.FlexRow),
		output:      tview.NewTextView(),
		progressBar: tvxwidgets.NewActivityModeGauge(),
	}

	bgColor := style.DialogBgColor
	outputBgColor := style.TerminalBgColor
	outputFgColor := style.TerminalFgColor
	buildPrgBorderColor := style.DialogSubBoxBorderColor
	prgCellColor := style.PrgBarColor

	// progressbar
	buildPrgDialog.progressBar.SetBorder(true)
	buildPrgDialog.progressBar.SetBorderColor(buildPrgBorderColor)
	buildPrgDialog.progressBar.SetPgBgColor(prgCellColor)

	// output
	buildPrgDialog.output.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	buildPrgDialog.output.SetBackgroundColor(outputBgColor)
	buildPrgDialog.output.SetTextColor(outputFgColor)

	// layout
	buildPrgDialog.layout.SetBackgroundColor(bgColor)
	buildPrgDialog.layout.SetBorder(true)
	buildPrgDialog.layout.SetBorderColor(style.DialogBorderColor)
	buildPrgDialog.layout.SetTitle("PODMAN IMAGE BUILD")
	buildPrgDialog.layout.AddItem(buildPrgDialog.progressBar, 3, 0, false) //nolint:mnd

	outputLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	outputLayout.AddItem(utils.EmptyBoxSpace(outputBgColor), 1, 0, false)
	outputLayout.AddItem(buildPrgDialog.output, 0, 1, false)
	outputLayout.AddItem(utils.EmptyBoxSpace(outputBgColor), 1, 0, false)
	buildPrgDialog.layout.AddItem(outputLayout, 0, 1, true)

	return buildPrgDialog
}

// Display displays this primitive.
func (d *ImageBuildProgressDialog) Display() {
	d.display = true
	d.cancelChan = make(chan bool)
	d.writerChan = make(chan []byte, 100) //nolint:mnd

	go d.outputReaderLoop()
}

// IsDisplay returns true if primitive is shown.
func (d *ImageBuildProgressDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ImageBuildProgressDialog) Hide() {
	d.display = false
	d.cancelChan <- true

	close(d.writerChan)

	d.output.SetText("")
	d.progressBar.Reset()
}

// HasFocus returns whether or not this primitive has focus.
func (d *ImageBuildProgressDialog) HasFocus() bool {
	if d.layout.HasFocus() || d.output.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ImageBuildProgressDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.layout)
}

// InputHandler returns input handler function for this primitive.
func (d *ImageBuildProgressDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("image build progress dialog: event %v received", event)
	})
}

// SetRect set rects for this primitive.
func (d *ImageBuildProgressDialog) SetRect(x, y, width, height int) {
	if width > buildPrgDialogMaxWidth {
		emptySpace := (width - buildPrgDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = buildPrgDialogMaxWidth
	}

	if height > buildPrgDialogHeight {
		emptySpace := (height - buildPrgDialogHeight) / 2 //nolint:mnd
		y += emptySpace
		height = buildPrgDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ImageBuildProgressDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// LogWriter returns output log writer.
func (d *ImageBuildProgressDialog) LogWriter() channel.WriteCloser { //nolint:ireturn
	return channel.NewWriter(d.writerChan)
}

// SetFastRefreshHandler sets fast refresh handler
// fast refresh is used to print image build output as fast as possible.
func (d *ImageBuildProgressDialog) SetFastRefreshHandler(handler func()) {
	d.fastRefreshHandler = handler
}

func (d *ImageBuildProgressDialog) outputReaderLoop() {
	tick := time.NewTicker(utils.RefreshInterval)

	log.Debug().Msg("image build progress dialog: output reader started")

	for {
		select {
		case <-tick.C:
			d.progressBar.Pulse()
		case <-d.cancelChan:
			log.Debug().Msg("image build progress dialog: output reader stopped")
			close(d.cancelChan)
			tick.Stop()

			return
		case data := <-d.writerChan:
			d.mu.Lock()

			_, err := d.output.Write(data)
			if err != nil {
				log.Error().Msgf("failed to write data to output: %s", err.Error())
			}

			d.mu.Unlock()
			d.fastRefreshHandler()
		}
	}
}
