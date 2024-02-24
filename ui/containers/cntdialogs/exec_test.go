package cntdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("container exec", Ordered, func() {
	var execDialogApp *tview.Application
	var execDialogScreen tcell.SimulationScreen
	var execDialog *ContainerExecDialog
	var runApp func()

	BeforeAll(func() {
		execDialogApp = tview.NewApplication()
		execDialog = NewContainerExecDialog()
		execDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := execDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := execDialogApp.SetScreen(execDialogScreen).SetRoot(execDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		execDialog.Display()
		execDialogApp.Draw()
		Expect(execDialog.IsDisplay()).To(Equal(true))
		Expect(execDialog.focusElement).To(Equal(execCommandFieldFocus))
	})

	It("set focus", func() {
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))
	})

	It("hide", func() {
		execDialog.Hide()
		Expect(execDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		execDialogApp.Stop()
	})
})
