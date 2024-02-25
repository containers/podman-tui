package imgdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("image search", Ordered, func() {
	var searchDialogApp *tview.Application
	var searchDialogScreen tcell.SimulationScreen
	var searchDialog *ImageSearchDialog
	var runApp func()

	BeforeAll(func() {
		searchDialogApp = tview.NewApplication()
		searchDialog = NewImageSearchDialog()
		searchDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := searchDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := searchDialogApp.SetScreen(searchDialogScreen).SetRoot(searchDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		searchDialog.Display()
		searchDialogApp.Draw()
		Expect(searchDialog.IsDisplay()).To(Equal(true))
		Expect(searchDialog.focusElement).To(Equal(sInputElement))
	})

	It("set focus", func() {
		searchDialogApp.SetFocus(searchDialog)
		searchDialogApp.Draw()
		Expect(searchDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"

		cancelFunc := func() {
			cancelAction = cancelWants
		}

		searchDialog.SetCancelFunc(cancelFunc)
		searchDialog.focusElement = sInputElement
		searchDialogApp.Draw()
		searchDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		searchDialogApp.Draw()
		searchDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		searchDialogApp.Draw()
		searchDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		searchDialogApp.Draw()
		searchDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		searchDialogApp.Draw()
		Expect(cancelWants).To(Equal(cancelAction))
	})

	It("search options", func() {
		searchDialog.focusElement = sInputElement
		searchDialogApp.SetFocus(searchDialog)
		searchDialogApp.Draw()
		searchDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone)) // (256,99,0) c character
		searchDialogApp.Draw()

		opts := searchDialog.GetSearchText()
		Expect(opts).To(Equal("c"))
	})

	It("hide", func() {
		searchDialog.Hide()
		Expect(searchDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		searchDialogApp.Stop()
	})
})
