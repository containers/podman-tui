package imgdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("image build", Ordered, func() {
	var buildDialogApp *tview.Application
	var buildDialogScreen tcell.SimulationScreen
	var buildDialog *ImageBuildDialog
	var runApp func()

	BeforeAll(func() {
		buildDialogApp = tview.NewApplication()
		buildDialog = NewImageBuildDialog()
		buildDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := buildDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := buildDialogApp.SetScreen(buildDialogScreen).SetRoot(buildDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		buildDialog.Display()
		buildDialogApp.Draw()
		Expect(buildDialog.IsDisplay()).To(Equal(true))
		Expect(buildDialog.focusElement).To(Equal(buildDialogContextDirectoryPathFieldFocus))
	})

	It("set focus", func() {
		buildDialogApp.SetFocus(buildDialog)
		buildDialogApp.Draw()
		Expect(buildDialog.HasFocus()).To(Equal(true))
	})

	It("hide", func() {
		buildDialog.Hide()
		Expect(buildDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		buildDialogApp.Stop()
	})
})
