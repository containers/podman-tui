package imgdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("image build progress", Ordered, func() {
	var buildDialogProgressApp *tview.Application
	var buildProgressDialogScreen tcell.SimulationScreen
	var buildProgressDialog *ImageBuildProgressDialog
	var runApp func()

	BeforeAll(func() {
		buildDialogProgressApp = tview.NewApplication()
		buildProgressDialog = NewImageBuildProgressDialog()
		buildProgressDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := buildProgressDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := buildDialogProgressApp.SetScreen(buildProgressDialogScreen).SetRoot(buildProgressDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		buildProgressDialog.Display()
		buildDialogProgressApp.Draw()
		Expect(buildProgressDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		buildDialogProgressApp.SetFocus(buildProgressDialog)
		buildDialogProgressApp.Draw()
		Expect(buildProgressDialog.HasFocus()).To(Equal(true))
	})

	It("hide", func() {
		buildProgressDialog.Hide()
		Expect(buildProgressDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		buildDialogProgressApp.Stop()
	})
})
