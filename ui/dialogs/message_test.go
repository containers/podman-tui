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
	var app *tview.Application
	var screen tcell.SimulationScreen
	var messageDialog *MessageDialog
	var messageText string = "this is a test message"
	var runApp func()

	BeforeAll(func() {
		app = tview.NewApplication()
		messageDialog = NewMessageDialog(messageText)
		screen = tcell.NewSimulationScreen("UTF-8")
		err := screen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := app.SetScreen(screen).SetRoot(messageDialog, true).Run(); err != nil {
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
		app.SetFocus(messageDialog)
		Expect(messageDialog.HasFocus()).To(Equal(true))
	})

	It("set title", func() {
		title := "podman"
		messageDialog.SetTitle(title)
		Expect(messageDialog.layout.GetTitle()).To(Equal(strings.ToUpper(title)))
	})

	It("set text", func() {
		messageDialog.SetText(messageText)
		Expect(messageDialog.textview.GetText(true)).To(Equal(messageText))
	})

	It("set rect", func() {
		x := 0
		y := 0
		width := 50
		height := 20
		// wants
		wWants := 40
		hWants := 8
		xWants := 5 // 0 + (50-40)/2
		yWants := 6 // 0 + (20-8)/2

		messageDialog.SetRect(x, y, width, height)
		x1, y1, w1, h1 := messageDialog.Box.GetRect()
		Expect(x1).To(Equal(xWants))
		Expect(y1).To(Equal(yWants))
		Expect(w1).To(Equal(wWants))
		Expect(h1).To(Equal(hWants))
	})

	It("enter button selected", func() {
		enterButton := "initial"
		enterButtonWants := "enter selected"
		enterFunc := func() {
			enterButton = enterButtonWants
		}
		messageDialog.SetSelectedFunc(enterFunc)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		Expect(enterButton).To(Equal(enterButtonWants))
	})

	It("cancel button selected", func() {
		cancelButton := "initial"
		cancelButtonWants := "cancel selected"
		cancelFunc := func() {
			cancelButton = cancelButtonWants
		}
		messageDialog.SetCancelFunc(cancelFunc)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		app.Draw()
		Expect(cancelButton).To(Equal(cancelButtonWants))
	})

	It("hide", func() {
		messageDialog.Hide()
		Expect(messageDialog.IsDisplay()).To(Equal(false))
		Expect(messageDialog.textview.GetText(true)).To(Equal(""))
	})

	AfterAll(func() {
		app.Stop()
	})

})
