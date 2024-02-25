package imgdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("image save", Ordered, func() {
	var saveDialogApp *tview.Application
	var saveDialogScreen tcell.SimulationScreen
	var saveDialog *ImageSaveDialog
	var runApp func()

	BeforeAll(func() {
		saveDialogApp = tview.NewApplication()
		saveDialog = NewImageSaveDialog()
		saveDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := saveDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := saveDialogApp.SetScreen(saveDialogScreen).SetRoot(saveDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		saveDialog.Display()
		saveDialogApp.Draw()
		Expect(saveDialog.IsDisplay()).To(Equal(true))
		Expect(saveDialog.focusElement).To(Equal(imageSaveOutputFocus))
	})

	It("set focus", func() {
		saveDialogApp.SetFocus(saveDialog)
		saveDialogApp.Draw()
		Expect(saveDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		saveDialog.SetCancelFunc(cancelFunc)
		saveDialog.focusElement = imageSaveFormFocus
		saveDialogApp.SetFocus(saveDialog)
		saveDialogApp.Draw()
		saveDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		saveDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("save options", func() {
		saveDialog.focusElement = imageSaveOutputFocus
		saveDialogApp.SetFocus(saveDialog)
		saveDialogApp.Draw()
		saveDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone)) // (256,99,0) c character
		saveDialogApp.Draw()

		opts, err := saveDialog.ImageSaveOptions()
		Expect(err).To(BeNil())
		Expect(opts.Output).To(Equal("c"))
	})

	It("hide", func() {
		saveDialog.Hide()
		Expect(saveDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		saveDialogApp.Stop()
	})
})
