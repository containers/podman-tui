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

	It("hide", func() {
		saveDialog.Hide()
		Expect(saveDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		saveDialogApp.Stop()
	})
})
