package sysdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("system event", Ordered, func() {
	var eventDialogApp *tview.Application
	var eventDialogScreen tcell.SimulationScreen
	var eventDialog *EventsDialog
	var runApp func()

	BeforeAll(func() {
		eventDialogApp = tview.NewApplication()
		eventDialog = NewEventDialog()
		eventDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := eventDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := eventDialogApp.SetScreen(eventDialogScreen).SetRoot(eventDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		eventDialog.Display()
		eventDialogApp.Draw()
		Expect(eventDialog.IsDisplay()).To(Equal(true))
		Expect(eventDialog.focusElement).To(Equal(formFieldHasFocus))
	})

	It("set focus", func() {
		eventDialogApp.SetFocus(eventDialog)
		eventDialogApp.Draw()
		Expect(eventDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		eventDialog.SetCancelFunc(cancelFunc)
		eventDialogApp.SetFocus(eventDialog)
		eventDialogApp.Draw()
		eventDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		eventDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("hide", func() {
		eventDialog.Hide()
		Expect(eventDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		eventDialogApp.Stop()
	})
})
