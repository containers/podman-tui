package dialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("error dialog", Ordered, func() {
	var errorDialogApp *tview.Application
	var errorDialogScreen tcell.SimulationScreen
	var errorDialog *ErrorDialog
	var runApp func()

	BeforeAll(func() {
		errorDialogApp = tview.NewApplication()
		errorDialog = NewErrorDialog()
		errorDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := errorDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := errorDialogApp.SetScreen(errorDialogScreen).SetRoot(errorDialog, true).Run(); err != nil {
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
		errorDialogApp.SetFocus(errorDialog)
		Expect(errorDialog.HasFocus()).To(Equal(true))
	})

	It("enter button selected", func() {
		enterButton := "initial"
		enterButtonWants := "enter selected"
		enterFunc := func() {
			enterButton = enterButtonWants
		}
		errorDialog.SetDoneFunc(enterFunc)
		errorDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		errorDialogApp.Draw()
		Expect(enterButton).To(Equal(enterButtonWants))
	})

	It("hide", func() {
		errorDialog.Hide()
		Expect(errorDialog.IsDisplay()).To(Equal(false))
		Expect(errorDialog.message).To(Equal(""))
	})

	AfterAll(func() {
		errorDialogApp.Stop()
	})

})
