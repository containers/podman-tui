package imgdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("image import", Ordered, func() {
	var importDialogApp *tview.Application
	var importDialogScreen tcell.SimulationScreen
	var importDialog *ImageImportDialog
	var runApp func()

	BeforeAll(func() {
		importDialogApp = tview.NewApplication()
		importDialog = NewImageImportDialog()
		importDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := importDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := importDialogApp.SetScreen(importDialogScreen).SetRoot(importDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		importDialog.Display()
		importDialogApp.Draw()
		Expect(importDialog.IsDisplay()).To(Equal(true))
		Expect(importDialog.focusElement).To(Equal(imageImportPathFocus))
	})

	It("set focus", func() {
		importDialogApp.SetFocus(importDialog)
		importDialogApp.Draw()
		Expect(importDialog.HasFocus()).To(Equal(true))
	})

	It("hide", func() {
		importDialog.Hide()
		Expect(importDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		importDialogApp.Stop()
	})
})
