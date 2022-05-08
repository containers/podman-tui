package dialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("error dialog", Ordered, func() {
	var app *tview.Application
	var screen tcell.SimulationScreen
	var errorDialog *ErrorDialog
	var runApp func()

	BeforeAll(func() {
		app = tview.NewApplication()
		errorDialog = NewErrorDialog()
		screen = tcell.NewSimulationScreen("UTF-8")
		err := screen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := app.SetScreen(screen).SetRoot(errorDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		errorDialog.Display()
		Expect(errorDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		app.SetFocus(errorDialog)
		Expect(errorDialog.HasFocus()).To(Equal(true))
	})

	It("enter button selected", func() {
		enterButton := "inital"
		enterButtonWants := "enter selected"
		enterFunc := func() {
			enterButton = enterButtonWants
		}
		errorDialog.SetDoneFunc(enterFunc)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		Expect(enterButton).To(Equal(enterButtonWants))
	})

	It("hide", func() {
		errorDialog.Hide()
		Expect(errorDialog.IsDisplay()).To(Equal(false))
		Expect(errorDialog.message).To(Equal(""))
	})

	AfterAll(func() {
		app.Stop()
	})

})
