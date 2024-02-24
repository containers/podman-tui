package cntdialogs

import (
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("container stat", Ordered, func() {
	var statDialogApp *tview.Application
	var statDialogScreen tcell.SimulationScreen
	var statDialog *ContainerStatsDialog
	var runApp func()

	BeforeAll(func() {
		statDialogApp = tview.NewApplication()
		statDialog = NewContainerStatsDialog()
		statDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := statDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := statDialogApp.SetScreen(statDialogScreen).SetRoot(statDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		statChannel := make(chan entities.ContainerStatsReport)
		statStream := true
		statDialog.SetStatsChannel(&statChannel)
		statDialog.SetStatsStream(&statStream)
		statDialog.Display()
		statDialogApp.Draw()
		Expect(statDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		statDialogApp.SetFocus(statDialog)
		statDialogApp.Draw()
		Expect(statDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		statDialog.SetDoneFunc(cancelFunc)
		statDialogApp.SetFocus(statDialog)
		statDialogApp.Draw()
		statDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		statDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("hide", func() {
		statDialog.Hide()
		Expect(statDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		statDialogApp.Stop()
	})
})
