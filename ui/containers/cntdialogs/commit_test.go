package cntdialogs

import (
	"fmt"

	"github.com/containers/buildah/define"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("container commit", Ordered, func() {
	var containerCommitApp *tview.Application
	var containerCommitScreen tcell.SimulationScreen
	var commitDialog *ContainerCommitDialog
	var runApp func()

	BeforeAll(func() {
		containerCommitApp = tview.NewApplication()
		commitDialog = NewContainerCommitDialog()
		containerCommitScreen = tcell.NewSimulationScreen("UTF-8")
		err := containerCommitScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := containerCommitApp.SetScreen(containerCommitScreen).SetRoot(commitDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		commitDialog.Display()
		Expect(commitDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		containerCommitApp.SetFocus(commitDialog)
		Expect(commitDialog.HasFocus()).To(Equal(true))
	})

	It("set container info", func() {
		cntID := "cntID"
		cntName := "cntName"
		cntInfoWants := fmt.Sprintf("%s (%s)", cntID, cntName)
		commitDialog.SetContainerInfo(cntID, cntName)
		Expect(commitDialog.cntInfo.GetText()).To(Equal(cntInfoWants))
	})

	It("cancel button selected", func() {
		cancelFunc := func() {
			commitDialog.Hide()
		}
		commitDialog.Hide()
		containerCommitApp.Draw()
		commitDialog.SetCancelFunc(cancelFunc)
		commitDialog.Display()
		containerCommitApp.Draw()
		containerCommitApp.SetFocus(commitDialog.form)
		containerCommitApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		containerCommitApp.Draw()
		Expect(commitDialog.IsDisplay()).To(Equal(false))
	})

	It("commit button selected", func() {
		commitButton := "initial"
		commitButtonWants := "commit selected"
		commitFunc := func() {
			commitButton = commitButtonWants
		}
		commitDialog.Hide()
		containerCommitApp.Draw()
		commitDialog.SetCommitFunc(commitFunc)
		commitDialog.Display()
		containerCommitApp.Draw()
		containerCommitApp.SetFocus(commitDialog.form)
		containerCommitApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		containerCommitApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		containerCommitApp.Draw()
		Expect(commitButton).To(Equal(commitButtonWants))
	})

	It("set commit handler", func() {
		commit := "initial"
		commitWants := "commit"
		commitHandler := func() {
			commit = commitWants
		}
		commitDialog.SetCommitFunc(commitHandler)
		commitDialog.commitHandler()
		Expect(commit).To(Equal(commitWants))
	})

	It("set cancel handler", func() {
		cancel := "initial"
		cancelWants := "cancel"
		cancelHandler := func() {
			cancel = cancelWants
		}
		commitDialog.SetCancelFunc(cancelHandler)
		commitDialog.cancelHandler()
		Expect(cancel).To(Equal(cancelWants))
	})

	It("get commit options", func() {
		opts := struct {
			Image   string
			Author  string
			Change  []string
			Format  string
			Squash  bool
			Pause   bool
			Message string
		}{
			Image:   "a",                // (256,97,0)
			Author:  "b",                // (256,98,0)
			Change:  []string{"c", "d"}, // (256,99,0) (256,100,0)
			Format:  define.Dockerv2ImageManifest,
			Squash:  true,
			Pause:   true,
			Message: "e", // (256,101,0)
		}

		commitDialog.Hide()
		containerCommitApp.Draw()
		commitDialog.Display()
		containerCommitApp.SetFocus(commitDialog)
		containerCommitApp.Draw()
		// image input field
		containerCommitApp.QueueEvent(tcell.NewEventKey(256, 97, tcell.ModNone))
		containerCommitApp.Draw()
		// author input field
		commitDialog.setFocusElement()
		containerCommitApp.SetFocus(commitDialog)
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(256, 98, tcell.ModNone))
		containerCommitApp.Draw()
		// change input field
		commitDialog.setFocusElement()
		containerCommitApp.SetFocus(commitDialog)
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone))
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(256, 32, tcell.ModNone)) // space
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(256, 100, tcell.ModNone))
		containerCommitApp.Draw()
		// format dropdown
		commitDialog.setFocusElement()
		containerCommitApp.SetFocus(commitDialog)
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		containerCommitApp.Draw()
		// squash checkbox
		commitDialog.setFocusElement()
		containerCommitApp.SetFocus(commitDialog)
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		containerCommitApp.Draw()
		// pause checkbox
		commitDialog.setFocusElement()
		containerCommitApp.SetFocus(commitDialog)
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		containerCommitApp.Draw()
		// message input field
		commitDialog.setFocusElement()
		containerCommitApp.SetFocus(commitDialog)
		containerCommitApp.Draw()
		containerCommitApp.QueueEvent(tcell.NewEventKey(256, 101, tcell.ModNone))
		containerCommitApp.Draw()

		// get and check commit options
		commitOpts := commitDialog.GetContainerCommitOptions()
		Expect(commitOpts.Image).To(Equal(opts.Image))
		Expect(commitOpts.Author).To(Equal(opts.Author))
		Expect(len(commitOpts.Changes)).To(Equal(len(opts.Change)))
		Expect(commitOpts.Changes[0]).To(Equal(opts.Change[0]))
		Expect(commitOpts.Changes[1]).To(Equal(opts.Change[1]))
		Expect(commitOpts.Format).To(Equal(opts.Format))
		Expect(commitOpts.Squash).To(Equal(opts.Squash))
		Expect(commitOpts.Pause).To(Equal(opts.Squash))
		Expect(commitOpts.Message).To(Equal(opts.Message))

	})

	AfterAll(func() {
		containerCommitApp.Stop()
	})

})
