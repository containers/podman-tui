package dialogs

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("message dialog", Ordered, func() {
	var messageDialogApp *tview.Application
	var messageDialogScreen tcell.SimulationScreen
	var messageDialog *MessageDialog
	var messageText string = "this is a test message"
	var runApp func()

	BeforeAll(func() {
		messageDialogApp = tview.NewApplication()
		messageDialog = NewMessageDialog(messageText)
		messageDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := messageDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := messageDialogApp.SetScreen(messageDialogScreen).SetRoot(messageDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		messageDialog.Display()
		Expect(messageDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		messageDialogApp.SetFocus(messageDialog)
		Expect(messageDialog.HasFocus()).To(Equal(true))
	})

	It("set title", func() {
		title := "podman"
		messageDialog.SetTitle(title)
		Expect(messageDialog.layout.GetTitle()).To(Equal(strings.ToUpper(title)))
	})

	It("set text", func() {
		messageDialog.SetText(MessageContainerInfo, "ID (NAME)", messageText)
		Expect(messageDialog.textview.GetText(true)).To(Equal(messageText))
	})

	It("set rect", func() {
		x := 0
		y := 0
		width := 50
		height := 20
		// wants
		wWants := 40
		hWants := 11
		xWants := 5 // 0 + (50-40)/2
		yWants := 4 // 0 + (20-8)/2 - 2

		messageDialog.SetRect(x, y, width, height)
		x1, y1, w1, h1 := messageDialog.Box.GetRect()
		Expect(x1).To(Equal(xWants))
		Expect(y1).To(Equal(yWants))
		Expect(w1).To(Equal(wWants))
		Expect(h1).To(Equal(hWants))
	})

	It("cancel button selected", func() {
		cancelButton := "initial"
		cancelButtonWants := "cancel selected"
		cancelFunc := func() {
			cancelButton = cancelButtonWants
		}
		messageDialog.SetCancelFunc(cancelFunc)
		messageDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		messageDialogApp.Draw()
		Expect(cancelButton).To(Equal(cancelButtonWants))
	})

	It("hide", func() {
		messageDialog.Hide()
		Expect(messageDialog.IsDisplay()).To(Equal(false))
		Expect(messageDialog.textview.GetText(true)).To(Equal(""))
	})

	AfterAll(func() {
		messageDialogApp.Stop()
	})

})
