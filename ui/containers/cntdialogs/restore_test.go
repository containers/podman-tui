package cntdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("container restore", Ordered, func() {
	var restoreDialogApp *tview.Application
	var restoreDialogScreen tcell.SimulationScreen
	var restoreDialog *ContainerRestoreDialog
	var runApp func()

	BeforeAll(func() {
		restoreDialogApp = tview.NewApplication()
		restoreDialog = NewContainerRestoreDialog()
		restoreDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := restoreDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := restoreDialogApp.SetScreen(restoreDialogScreen).SetRoot(restoreDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		restoreDialog.Display()
		Expect(restoreDialog.IsDisplay()).To(Equal(true))
		Expect(restoreDialog.focusElement).To(Equal(cntRestoreContainersFocus))
	})

	It("set focus", func() {
		restoreDialogApp.SetFocus(restoreDialog)
		//restoreDialogApp.Draw()
		Expect(restoreDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		restoreDialog.SetCancelFunc(cancelFunc)
		restoreDialog.focusElement = cntRestoreFormFocus
		restoreDialogApp.SetFocus(restoreDialog)
		restoreDialogApp.Draw()
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		restoreDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("restore button selected", func() {
		restoreWants := "restore selected"
		restoreAction := "restore init"
		cancelFunc := func() {
			restoreAction = restoreWants
		}
		restoreDialog.SetRestoreFunc(cancelFunc)
		restoreDialog.focusElement = cntRestoreFormFocus
		restoreDialogApp.SetFocus(restoreDialog)
		restoreDialogApp.Draw()
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		restoreDialogApp.Draw()
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		restoreDialogApp.Draw()
		Expect(restoreAction).To(Equal(restoreWants))
	})

	It("restore set containers", func() {
		cntList := [][]string{
			{"aaaaaaaaaaaa", "cnt01"},
			{"bbbbbbbbbbbb", "cnt02"},
		}

		restoreDialog.SetContainers(cntList)
		restoreDialog.Display()
		restoreDialogApp.Draw()

		optionCount := restoreDialog.containers.GetOptionCount()
		Expect(optionCount).To(Equal(len(cntList) + 1))
	})

	It("restore set pods", func() {
		podList := [][]string{
			{"aaaaaaaaaaaa", "pod01"},
			{"bbbbbbbbbbbb", "pod02"},
		}

		restoreDialog.SetPods(podList)
		restoreDialog.Display()
		restoreDialogApp.Draw()

		optionCount := restoreDialog.pods.GetOptionCount()
		Expect(optionCount).To(Equal(len(podList) + 1))
	})

	It("restore get options", func() {
		cntList := [][]string{
			{"aaaaaaaaaaaa", "cnt01"},
			{"bbbbbbbbbbbb", "cnt02"},
		}
		podList := [][]string{
			{"aaaaaaaaaaaa", "pod01"},
			{"bbbbbbbbbbbb", "pod02"},
		}
		restoreDialog.SetContainers(cntList)
		restoreDialog.SetPods(podList)
		restoreDialog.Display()
		restoreDialogApp.Draw()
		restoreDialogApp.SetFocus(restoreDialog)
		restoreDialogApp.Draw()
		// container
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		restoreDialogApp.Draw()
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		restoreDialogApp.Draw()
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		restoreDialogApp.Draw()

		// pod
		restoreDialog.setFocusElement()
		restoreDialogApp.SetFocus(restoreDialog)
		restoreDialogApp.Draw()
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		restoreDialogApp.Draw()
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		restoreDialogApp.Draw()
		restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		restoreDialogApp.Draw()

		// name
		restoreDialog.setFocusElement()
		restoreDialogApp.SetFocus(restoreDialog)
		restoreDialogApp.Draw()
		// publish
		restoreDialog.setFocusElement()
		restoreDialogApp.SetFocus(restoreDialog)
		restoreDialogApp.Draw()
		// Import
		restoreDialog.setFocusElement()
		restoreDialogApp.SetFocus(restoreDialog)
		restoreDialogApp.Draw()

		for i := 0; i < 8; i++ {
			restoreDialog.setFocusElement()
			restoreDialogApp.SetFocus(restoreDialog)
			restoreDialogApp.Draw()
			restoreDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			restoreDialogApp.Draw()
		}

		restoreOpts := restoreDialog.GetRestoreOptions()

		Expect(restoreOpts.ContainerID).To(Equal(cntList[0][0]))
		Expect(restoreOpts.PodID).To(Equal(podList[0][0]))
		Expect(restoreOpts.Keep).To(Equal(true))
		Expect(restoreOpts.IgnoreStaticIP).To(Equal(true))
		Expect(restoreOpts.IgnoreStaticMAC).To(Equal(true))
		Expect(restoreOpts.FileLocks).To(Equal(true))
		Expect(restoreOpts.PrintStats).To(Equal(true))
		Expect(restoreOpts.TCPEstablished).To(Equal(true))
		Expect(restoreOpts.IgnoreVolumes).To(Equal(true))
		Expect(restoreOpts.IgnoreRootfs).To(Equal(true))

	})

	It("hide", func() {
		restoreDialog.Hide()
		Expect(restoreDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		restoreDialogApp.Stop()
	})
})
