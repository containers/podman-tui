package sysdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("system df", Ordered, func() {
	var dfDialogApp *tview.Application
	var dfDialogScreen tcell.SimulationScreen
	var dfDialog *DfDialog
	var runApp func()

	BeforeAll(func() {
		dfDialogApp = tview.NewApplication()
		dfDialog = NewDfDialog()
		dfDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := dfDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := dfDialogApp.SetScreen(dfDialogScreen).SetRoot(dfDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		dfDialog.Display()
		dfDialogApp.Draw()
		Expect(dfDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		dfDialogApp.SetFocus(dfDialog)
		dfDialogApp.Draw()
		Expect(dfDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		dfDialog.SetCancelFunc(cancelFunc)
		dfDialogApp.SetFocus(dfDialog)
		dfDialogApp.Draw()
		dfDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		dfDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("hide", func() {
		dfDialog.Hide()
		Expect(dfDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		dfDialogApp.Stop()
	})
})
