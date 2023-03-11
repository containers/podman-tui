package cntdialogs

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("container checkpoint", Ordered, func() {
	var checkpointDialogApp *tview.Application
	var checkpointDialogScreen tcell.SimulationScreen
	var checkpointDialog *ContainerCheckpointDialog
	var runApp func()

	BeforeAll(func() {
		checkpointDialogApp = tview.NewApplication()
		checkpointDialog = NewContainerCheckpointDialog()
		checkpointDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := checkpointDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := checkpointDialogApp.SetScreen(checkpointDialogScreen).SetRoot(checkpointDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		checkpointDialog.Display()
		Expect(checkpointDialog.IsDisplay()).To(Equal(true))
		Expect(checkpointDialog.focusElement).To(Equal(cntCheckpointImageFocus))
	})

	It("set focus", func() {
		checkpointDialogApp.SetFocus(checkpointDialog)
		Expect(checkpointDialog.HasFocus()).To(Equal(true))
	})

	It("set container info", func() {
		cntID := "cntID"
		cntName := "cntName"
		cntInfoWants := fmt.Sprintf("%s (%s)", cntID, cntName)
		checkpointDialog.SetContainerInfo(cntID, cntName)
		Expect(strings.TrimSpace(checkpointDialog.containerInfo.GetText())).To(Equal(cntInfoWants))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		checkpointDialog.SetCancelFunc(cancelFunc)
		checkpointDialog.focusElement = cntCheckpointFormFocus
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("checkpoint button selected", func() {
		checkpointWants := "checkpoint selected"
		checkpointAction := "checkpoint init"
		cancelFunc := func() {
			checkpointAction = checkpointWants
		}
		checkpointDialog.SetCheckpointFunc(cancelFunc)
		checkpointDialog.focusElement = cntCheckpointFormFocus
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()
		Expect(checkpointAction).To(Equal(checkpointWants))
	})

	It("get checkpoint options", func() {
		opts := struct {
			createImage    string
			export         string
			filelock       bool
			ignoreRootsFS  bool
			tcpEstablished bool
			keep           bool
			printStats     bool
			preCheckpoint  bool
			leaveRunning   bool
			withPrevious   bool
		}{
			createImage:    "a", // (256,97,0)
			export:         "c", // (256,99,0)
			filelock:       true,
			ignoreRootsFS:  true,
			tcpEstablished: true,
			keep:           true,
			printStats:     true,
			leaveRunning:   true,
			withPrevious:   true,
			preCheckpoint:  true,
		}

		checkpointDialog.Hide()
		checkpointDialogApp.Draw()
		checkpointDialog.Display()
		checkpointDialogApp.Draw()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()

		// image
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(256, 97, tcell.ModNone))
		checkpointDialogApp.Draw()

		// export
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone))
		checkpointDialogApp.Draw()

		// create file lock
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()

		// ignore rootfs
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()

		// tcp established
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()

		// pre checkpoint
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()

		// print stats
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()

		// keep
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()

		// leave running
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()

		// withPrevious
		checkpointDialog.setFocusElement()
		checkpointDialogApp.SetFocus(checkpointDialog)
		checkpointDialogApp.Draw()
		checkpointDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		checkpointDialogApp.Draw()

		checkpointOptions := checkpointDialog.GetCheckpointOptions()
		Expect(checkpointOptions.CreateImage).To(Equal(opts.createImage))
		Expect(checkpointOptions.Export).To(Equal(opts.export))
		Expect(checkpointOptions.FileLocks).To(Equal(opts.filelock))
		Expect(checkpointOptions.IgnoreRootFs).To(Equal(opts.ignoreRootsFS))
		Expect(checkpointOptions.TCPEstablished).To(Equal(opts.tcpEstablished))
		Expect(checkpointOptions.Keep).To(Equal(opts.keep))
		Expect(checkpointOptions.PrintStats).To(Equal(opts.printStats))
		Expect(checkpointOptions.PreCheckpoint).To(Equal(opts.preCheckpoint))
		Expect(checkpointOptions.LeaveRunning).To(Equal(opts.leaveRunning))
		Expect(checkpointOptions.WithPrevious).To(Equal(opts.withPrevious))

	})

	It("hide", func() {
		checkpointDialog.Hide()
		Expect(checkpointDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		checkpointDialogApp.Stop()
	})
})
