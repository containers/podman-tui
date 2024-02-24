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

	It("has focus", func() {
		execDialog.focusElement = execInteractiveFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))

		execDialog.focusElement = execTtyFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))

		execDialog.focusElement = execPrivilegedFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))

		execDialog.focusElement = execWorkingDirFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))

		execDialog.focusElement = execEnvVariablesFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))

		execDialog.focusElement = execEnvFileFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))

		execDialog.focusElement = execUserFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))

		execDialog.focusElement = execDetachFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))

		execDialog.focusElement = execFormFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		Expect(execDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		execDialog.SetCancelFunc(cancelFunc)
		execDialog.focusElement = execFormFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		execDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		execDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("exec options", func() {
		execDialog.focusElement = execCommandFieldFocus
		execDialogApp.SetFocus(execDialog)
		execDialogApp.Draw()
		execDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone)) // (256,99,0) c character
		execDialogApp.Draw()

		opts := execDialog.ContainerExecOptions()
		Expect(opts.Cmd[0]).To(Equal("c"))
	})

	It("hide", func() {
		execDialog.Hide()
		Expect(execDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		execDialogApp.Stop()
	})
})
