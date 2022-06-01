package dialogs

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("confirm dialog", Ordered, func() {
	var app *tview.Application
	var screen tcell.SimulationScreen
	var confirmDialog *ConfirmDialog
	var runApp func()

	BeforeAll(func() {
		app = tview.NewApplication()
		confirmDialog = NewConfirmDialog()
		screen = tcell.NewSimulationScreen("UTF-8")
		err := screen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := app.SetScreen(screen).SetRoot(confirmDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		confirmDialog.Display()
		Expect(confirmDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		app.SetFocus(confirmDialog)
		Expect(confirmDialog.HasFocus()).To(Equal(true))
	})

	It("set title", func() {
		title := "podman"
		confirmDialog.SetTitle(title)
		Expect(confirmDialog.layout.GetTitle()).To(Equal(strings.ToUpper(title)))
	})

	It("set confirm text message", func() {
		confirmMsg := "test confirm message line01\ntest confirm message line02"
		confirmMsgWants := "\ntest confirm message line01\ntest confirm message line02"
		confirmDialog.SetText(confirmMsg)
		Expect(confirmDialog.textview.GetText(true)).To(Equal(confirmMsgWants))
	})

	It("enter button selected", func() {
		enterButton := "initial"
		enterButtonWants := "enter selected"
		enterFunc := func() {
			enterButton = enterButtonWants
		}
		confirmDialog.SetSelectedFunc(enterFunc)
		confirmDialog.Display()
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		Expect(enterButton).To(Equal(enterButtonWants))
	})

	It("cancel button selected", func() {
		cancelButton := "initial"
		cancelButtonWants := "cancel selected"
		cancelFunc := func() {
			cancelButton = cancelButtonWants
		}
		confirmDialog.SetCancelFunc(cancelFunc)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		app.Draw()
		Expect(cancelButton).To(Equal(cancelButtonWants))
		cancelButton = "initial"
		confirmDialog.Display()
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		Expect(cancelButton).To(Equal(cancelButtonWants))
	})

	It("hide", func() {
		confirmDialog.Hide()
		Expect(confirmDialog.IsDisplay()).To(Equal(false))
		Expect(confirmDialog.message).To(Equal(""))
	})

	AfterAll(func() {
		app.Stop()
	})

})
