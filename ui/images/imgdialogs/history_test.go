package imgdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("image history", Ordered, func() {
	var historyDialogApp *tview.Application
	var historyDialogScreen tcell.SimulationScreen
	var historyDialog *ImageHistoryDialog
	var runApp func()

	BeforeAll(func() {
		historyDialogApp = tview.NewApplication()
		historyDialog = NewImageHistoryDialog()
		historyDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := historyDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := historyDialogApp.SetScreen(historyDialogScreen).SetRoot(historyDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		historyDialog.Display()
		historyDialogApp.Draw()
		Expect(historyDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		historyDialogApp.SetFocus(historyDialog)
		historyDialogApp.Draw()
		Expect(historyDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"

		cancelFunc := func() {
			cancelAction = cancelWants
		}

		historyDialog.SetCancelFunc(cancelFunc)
		historyDialogApp.Draw()
		historyDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		historyDialogApp.Draw()
		historyDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		historyDialogApp.Draw()
		Expect(cancelWants).To(Equal(cancelAction))
	})

	It("hide", func() {
		historyDialog.Hide()
		Expect(historyDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		historyDialogApp.Stop()
	})
})
