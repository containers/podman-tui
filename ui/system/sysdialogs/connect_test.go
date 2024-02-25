package sysdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("system connection connect", Ordered, func() {
	var connectDialogApp *tview.Application
	var connectDialogScreen tcell.SimulationScreen
	var connectDialog *ConnectDialog
	var runApp func()

	BeforeAll(func() {
		connectDialogApp = tview.NewApplication()
		connectDialog = NewConnectDialog()
		connectDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := connectDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := connectDialogApp.SetScreen(connectDialogScreen).SetRoot(connectDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		connectDialog.Display()
		connectDialogApp.Draw()
		Expect(connectDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		connectDialogApp.SetFocus(connectDialog)
		connectDialogApp.Draw()
		Expect(connectDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		connectDialog.SetCancelFunc(cancelFunc)
		connectDialogApp.SetFocus(connectDialog)
		connectDialogApp.Draw()
		connectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		connectDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("hide", func() {
		connectDialog.Hide()
		Expect(connectDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		connectDialogApp.Stop()
	})
})
