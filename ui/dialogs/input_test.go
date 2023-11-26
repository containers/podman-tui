package dialogs

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("input dialog", Ordered, func() {
	var inputDialogApp *tview.Application
	var inputDialogScreen tcell.SimulationScreen
	var inputDialog *SimpleInputDialog
	var runApp func()

	BeforeAll(func() {
		inputDialogApp = tview.NewApplication()
		inputDialog = NewSimpleInputDialog("")
		inputDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := inputDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := inputDialogApp.SetScreen(inputDialogScreen).SetRoot(inputDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		inputDialog.Display()
		Expect(inputDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		inputDialogApp.SetFocus(inputDialog)
		Expect(inputDialog.HasFocus()).To(Equal(true))
	})

	It("set title", func() {
		title := "podman"
		inputDialog.SetTitle(title)
		Expect(inputDialog.layout.GetTitle()).To(Equal(strings.ToUpper(title)))
	})

	It("set label", func() {
		label := "label01"
		inputDialog.SetLabel(label)
		Expect(inputDialog.input.GetLabel()).To(Equal(label + ": "))
	})

	It("set select button label", func() {
		buttonLabel := "Add"
		inputDialog.SetSelectButtonLabel(buttonLabel)
		selectButton := inputDialog.form.GetButton(inputDialog.form.GetButtonCount() - 1)
		Expect(selectButton.GetLabel()).To(Equal(buttonLabel))
	})

	It("set labyout height", func() {
		hasDesc := true
		inputDialog.setLayout(hasDesc)
		Expect(inputDialog.height).To(Equal(siDialogHeight))

		hasDesc = false
		inputDialog.setLayout(hasDesc)
		Expect(inputDialog.height).To(Equal(siDialogHeight - 3))
	})

	It("set description", func() {
		description := "test description"
		inputDialog.SetDescription(description)
		wantedDesc := fmt.Sprintf("\n%s", description)
		Expect(inputDialog.textview.GetText(true)).To(Equal(wantedDesc))
		Expect(inputDialog.height).To(Equal(siDialogHeight))

		description = ""
		inputDialog.SetDescription(description)
		wantedDesc = "\n"
		Expect(inputDialog.textview.GetText(true)).To(Equal(wantedDesc))
		Expect(inputDialog.height).To(Equal(siDialogHeight - 3))
	})

	It("set and get input", func() {
		inputText := "podman"
		inputDialog.SetInputText(inputText)
		Expect(inputDialog.GetInputText()).To(Equal(inputText))
	})

	It("enter button selected", func() {
		enterButton := "initial"
		enterButtonWants := "enter selected"
		enterFunc := func() {
			enterButton = enterButtonWants
		}
		inputDialog.SetSelectedFunc(enterFunc)
		inputDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		inputDialogApp.Draw()
		Expect(enterButton).To(Equal(enterButtonWants))
	})

	It("cancel button selected", func() {
		cancelButton := "initial"
		cancelButtonWants := "cancel selected"
		cancelFunc := func() {
			cancelButton = cancelButtonWants
		}
		inputDialog.SetCancelFunc(cancelFunc)
		inputDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
		inputDialogApp.Draw()
		Expect(cancelButton).To(Equal(cancelButtonWants))
	})

	It("hide", func() {
		inputDialog.Hide()
		Expect(inputDialog.IsDisplay()).To(Equal(false))
		Expect(inputDialog.input.GetText()).To(Equal(""))
	})

	It("input handler", func() {
		selectButton := inputDialog.form.GetButton(inputDialog.form.GetButtonCount() - 1)
		cancelButton := inputDialog.form.GetButton(inputDialog.form.GetButtonCount() - 2)
		inputDialog.Display()
		inputDialogApp.Draw()
		Expect(inputDialog.input.HasFocus()).To(Equal(true))
		inputDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		inputDialogApp.Draw()
		Expect(cancelButton.HasFocus()).To(Equal(true))
		inputDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		inputDialogApp.Draw()
		Expect(selectButton.HasFocus()).To(Equal(true))
		inputDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		inputDialogApp.Draw()
		Expect(inputDialog.input.HasFocus()).To(Equal(true))
	})

	AfterAll(func() {
		inputDialogApp.Stop()
	})
})
