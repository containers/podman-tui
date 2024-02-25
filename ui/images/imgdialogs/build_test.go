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

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		buildDialog.SetCancelFunc(cancelFunc)
		buildDialog.focusElement = buildDialogFormFocus
		buildDialogApp.SetFocus(buildDialog)
		buildDialogApp.Draw()
		buildDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		buildDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("build button selected", func() {
		buildWants := "build selected"
		buildAction := "build init"
		buildFunc := func() {
			buildAction = buildWants
		}
		buildDialog.SetBuildFunc(buildFunc)
		buildDialog.focusElement = buildDialogFormFocus
		buildDialogApp.SetFocus(buildDialog)
		buildDialogApp.Draw()
		buildDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTAB, 0, tcell.ModNone))
		buildDialogApp.Draw()
		buildDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		buildDialogApp.Draw()
		Expect(buildAction).To(Equal(buildWants))
	})

	It("build options", func() {
		buildDialog.focusElement = buildDialogContextDirectoryPathFieldFocus
		buildDialogApp.SetFocus(buildDialog)
		buildDialogApp.Draw()
		buildDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone)) // (256,99,0) c character
		buildDialogApp.Draw()

		opts, err := buildDialog.ImageBuildOptions()
		Expect(err).To(BeNil())
		Expect(opts.BuildOptions.ContextDirectory).To(Equal("c"))
	})

	It("hide", func() {
		buildDialog.Hide()
		Expect(buildDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		buildDialogApp.Stop()
	})
})
