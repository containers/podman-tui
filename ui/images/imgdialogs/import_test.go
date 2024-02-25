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

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		importDialog.SetCancelFunc(cancelFunc)
		importDialog.focusElement = imageImportFormFocus
		importDialogApp.SetFocus(importDialog)
		importDialogApp.Draw()
		importDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		importDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("import button selected", func() {
		importWants := "import selected"
		importAction := "import init"
		importFunc := func() {
			importAction = importWants
		}
		importDialog.SetImportFunc(importFunc)
		importDialog.focusElement = imageImportFormFocus
		importDialogApp.SetFocus(importDialog)
		importDialogApp.Draw()
		importDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTAB, 0, tcell.ModNone))
		importDialogApp.Draw()
		importDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		importDialogApp.Draw()
		Expect(importAction).To(Equal(importWants))
	})

	It("import options", func() {
		importDialog.focusElement = imageImportPathFocus
		importDialogApp.SetFocus(importDialog)
		importDialogApp.Draw()
		importDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone)) // (256,99,0) c character
		importDialogApp.Draw()

		opts, err := importDialog.ImageImportOptions()
		Expect(err).To(BeNil())
		Expect(opts.Source).To(Equal("c"))
	})

	It("hide", func() {
		importDialog.Hide()
		Expect(importDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		importDialogApp.Stop()
	})
})
