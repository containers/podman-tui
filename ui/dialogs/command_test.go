package dialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("isPrintableASCII", func() {
	It("accepts lowercase letters", func() {
		Expect(isPrintableASCII('a')).To(BeTrue())
		Expect(isPrintableASCII('z')).To(BeTrue())
	})

	It("accepts uppercase letters", func() {
		Expect(isPrintableASCII('A')).To(BeTrue())
		Expect(isPrintableASCII('Z')).To(BeTrue())
	})

	It("accepts digits", func() {
		Expect(isPrintableASCII('0')).To(BeTrue())
		Expect(isPrintableASCII('9')).To(BeTrue())
	})

	It("rejects spaces and punctuation", func() {
		Expect(isPrintableASCII(' ')).To(BeFalse())
		Expect(isPrintableASCII('-')).To(BeFalse())
		Expect(isPrintableASCII('_')).To(BeFalse())
		Expect(isPrintableASCII('.')).To(BeFalse())
		Expect(isPrintableASCII('/')).To(BeFalse())
	})

	It("rejects non-ASCII characters", func() {
		Expect(isPrintableASCII('é')).To(BeFalse())
		Expect(isPrintableASCII('ñ')).To(BeFalse())
		Expect(isPrintableASCII('中')).To(BeFalse())
	})
})

var _ = Describe("shortcut assignment", func() {
	It("assigns unique first character", func() {
		cmd := NewCommandDialog([][]string{
			{"checkpoint", "checkpoint desc"},
			{"commit", "commit desc"},
			{"create", "create desc"},
		})
		Expect(cmd.shortcuts[0]).To(Equal('c'))
		Expect(cmd.shortcuts[1]).To(Equal('o'))
		Expect(cmd.shortcuts[2]).To(Equal('r'))
	})

	It("skips duplicate characters", func() {
		cmd := NewCommandDialog([][]string{
			{"aaaa", "aaaa desc"},
			{"abbb", "abbb desc"},
			{"accc", "accc desc"},
		})
		Expect(cmd.shortcuts[0]).To(Equal('a'))
		Expect(cmd.shortcuts[1]).To(Equal('b'))
		Expect(cmd.shortcuts[2]).To(Equal('c'))
	})

	It("skips non-alphanumeric characters in command names", func() {
		cmd := NewCommandDialog([][]string{
			{"my-command", "my desc"},
			{"my command", "another desc"},
		})
		Expect(cmd.shortcuts[0]).To(Equal('m'))
		Expect(cmd.shortcuts[1]).To(Equal('y'))
	})

	It("prefers numbers when they appear first", func() {
		cmd := NewCommandDialog([][]string{
			{"123abc", "desc1"},
			{"123def", "desc2"},
		})
		Expect(cmd.shortcuts[0]).To(Equal('1'))
		Expect(cmd.shortcuts[1]).To(Equal('2'))
	})

	It("handles all lowercase letters exhausted by falling back", func() {
		opts := make([][]string, 26)
		for i := 0; i < 26; i++ {
			second := 'a' + rune(i)
			opts[i] = []string{"a" + string(second), "desc"}
		}
		cmd := NewCommandDialog(opts)
		// All 26 should be assigned unique lowercase letters
		for i := 0; i < 26; i++ {
			Expect(cmd.shortcuts[i]).To(Equal(rune('a' + i)))
		}
	})

	It("handles duplicate first letters with many items", func() {
		cmd := NewCommandDialog([][]string{
			{"cmd01", "desc01"},
			{"cmd02", "desc02"},
			{"cmd03", "desc03"},
		})
		// All start with 'c', so first gets 'c', others fall back
		Expect(cmd.shortcuts[0]).To(Equal('c'))
		Expect(cmd.shortcuts[1]).NotTo(Equal('c'))
		Expect(cmd.shortcuts[2]).NotTo(Equal('c'))
		Expect(cmd.shortcuts[1]).NotTo(Equal(cmd.shortcuts[2]))
	})

	It("handles single item", func() {
		cmd := NewCommandDialog([][]string{
			{"single", "single desc"},
		})
		Expect(cmd.shortcuts[0]).To(Equal('s'))
	})

	It("handles uppercase letters", func() {
		cmd := NewCommandDialog([][]string{
			{"START", "start desc"},
			{"STOP", "stop desc"},
		})
		Expect(cmd.shortcuts[0]).To(Equal('S'))
		Expect(cmd.shortcuts[1]).To(Equal('T'))
	})
})

var _ = Describe("shortcut key handling", func() {
	var cmdDialogApp *tview.Application
	var cmdDialogScreen tcell.SimulationScreen
	var cmdDialog *CommandDialog
	var selectedCmd string
	var runApp func()

	BeforeEach(func() {
		cmdDialogApp = tview.NewApplication()
		cmdDialog = NewCommandDialog([][]string{
			{"checkpoint", "checkpoint desc"},
			{"commit", "commit desc"},
			{"create", "create desc"},
		})
		selectedCmd = ""
		cmdDialog.SetSelectedFunc(func() {
			selectedCmd = cmdDialog.GetSelectedItem()
		})
		cmdDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := cmdDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := cmdDialogApp.SetScreen(cmdDialogScreen).SetRoot(cmdDialog, false).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
		cmdDialog.Display()
	})

	AfterEach(func() {
		cmdDialogApp.Stop()
	})

	It("shortcut key selects and triggers handler", func() {
		// checkpoint should get 'c', commit gets 'o', create gets 'r'
		// Use 'o' to select commit
		cmdDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone))
		cmdDialogApp.Draw()
		Expect(selectedCmd).To(Equal("commit"))
	})

	It("shortcut key for first item works", func() {
		selectedCmd = ""
		cmdDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone))
		cmdDialogApp.Draw()
		Expect(selectedCmd).To(Equal("checkpoint"))
	})
})

var _ = Describe("command dialog", Ordered, func() {
	var cmdDialogApp *tview.Application
	var cmdDialogScreen tcell.SimulationScreen
	var cmdDialog *CommandDialog
	var cmdTitle [][]string
	var runApp func()

	BeforeAll(func() {
		cmdTitle = [][]string{
			{"cmd01", "cmd01 description"},
			{"cmd02", "cmd02 description"},
		}
		cmdDialogApp = tview.NewApplication()
		cmdDialog = NewCommandDialog(cmdTitle)

		cmdDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := cmdDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := cmdDialogApp.SetScreen(cmdDialogScreen).SetRoot(cmdDialog, false).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		cmdDialog.Display()
		Expect(cmdDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		cmdDialogApp.SetFocus(cmdDialog)
		cmdDialogApp.Draw()
		hasFocus := cmdDialog.HasFocus()
		Expect(hasFocus).To(Equal(true))
	})

	It("set rect", func() {
		x := 0
		y := 0
		width := 50
		height := 20
		ws := (width - cmdDialog.width) / 2
		hs := ((height - cmdDialog.height) / 2)
		yWants := y + hs
		xWants := x + ws
		cmdDialog.SetRect(x, y, width, height)
		x1, y1, w1, h1 := cmdDialog.Box.GetRect()
		Expect(x1).To(Equal(xWants))
		Expect(y1).To(Equal(yWants))
		Expect(w1).To(Equal(cmdDialog.width))
		Expect(h1).To(Equal(cmdDialog.height))
	})

	It("get total command counts", func() {
		// header + items (name:desc)
		Expect(cmdDialog.GetCommandCount()).To(Equal(3))
	})

	It("get selected item", func() {
		cmdDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		cmdDialogApp.Draw()
		Expect(cmdDialog.GetSelectedItem()).To(Equal("cmd02"))
	})

	It("command selected", func() {
		enterButton := "initial"
		enterButtonWants := "enter selected"
		enterFunc := func() {
			enterButton = enterButtonWants
		}
		cmdDialog.SetSelectedFunc(enterFunc)
		cmdDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		cmdDialogApp.Draw()
		Expect(enterButton).To(Equal(enterButtonWants))
	})

	It("cancel button selected", func() {
		cancelButton := "initial"
		cancelButtonWants := "cancel selected"
		cancelFunc := func() {
			cancelButton = cancelButtonWants
		}
		cmdDialog.SetCancelFunc(cancelFunc)
		cmdDialogApp.SetFocus(cmdDialog.form)
		cmdDialogApp.Draw()
		cmdDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		cmdDialogApp.Draw()
		Expect(cancelButton).To(Equal(cancelButtonWants))
	})

	It("hide", func() {
		cmdDialog.Hide()
		Expect(cmdDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		cmdDialogApp.Stop()
	})
})
