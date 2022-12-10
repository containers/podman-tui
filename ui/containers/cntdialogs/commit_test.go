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
	var app *tview.Application
	var screen tcell.SimulationScreen
	var commitDialog *ContainerCommitDialog
	var runApp func()

	BeforeAll(func() {
		app = tview.NewApplication()
		commitDialog = NewContainerCommitDialog()
		screen = tcell.NewSimulationScreen("UTF-8")
		err := screen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := app.SetScreen(screen).SetRoot(commitDialog, true).Run(); err != nil {
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
		app.SetFocus(commitDialog)
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
		app.Draw()
		commitDialog.SetCancelFunc(cancelFunc)
		commitDialog.Display()
		app.Draw()
		app.SetFocus(commitDialog.form)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		Expect(commitDialog.IsDisplay()).To(Equal(false))
	})

	It("commit button selected", func() {
		commitButton := "initial"
		commitButtonWants := "commit selected"
		commitFunc := func() {
			commitButton = commitButtonWants
		}
		commitDialog.Hide()
		app.Draw()
		commitDialog.SetCommitFunc(commitFunc)
		commitDialog.Display()
		app.Draw()
		app.SetFocus(commitDialog.form)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
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
		app.Draw()
		commitDialog.Display()
		app.SetFocus(commitDialog)
		app.Draw()
		// image input field
		app.QueueEvent(tcell.NewEventKey(256, 97, tcell.ModNone))
		app.Draw()
		// author input field
		commitDialog.setFocusElement()
		app.SetFocus(commitDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 98, tcell.ModNone))
		app.Draw()
		// change input field
		commitDialog.setFocusElement()
		app.SetFocus(commitDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone))
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 32, tcell.ModNone)) // space
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 100, tcell.ModNone))
		app.Draw()
		// format dropdown
		commitDialog.setFocusElement()
		app.SetFocus(commitDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		// squash checkbox
		commitDialog.setFocusElement()
		app.SetFocus(commitDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		// pause checkbox
		commitDialog.setFocusElement()
		app.SetFocus(commitDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		// message input field
		commitDialog.setFocusElement()
		app.SetFocus(commitDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 101, tcell.ModNone))
		app.Draw()

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
		app.Stop()
	})

})
