package sysdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("system connection add", Ordered, func() {
	var connAddDialogApp *tview.Application
	var connAddDialogScreen tcell.SimulationScreen
	var connAddDialog *AddConnectionDialog
	var runApp func()

	BeforeAll(func() {
		connAddDialogApp = tview.NewApplication()
		connAddDialog = NewAddConnectionDialog()
		connAddDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := connAddDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := connAddDialogApp.SetScreen(connAddDialogScreen).SetRoot(connAddDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		connAddDialog.Display()
		connAddDialogApp.Draw()
		Expect(connAddDialog.IsDisplay()).To(Equal(true))
		Expect(connAddDialog.focusElement).To(Equal(connNameFieldFocus))
	})

	It("set focus", func() {
		connAddDialogApp.SetFocus(connAddDialog)
		connAddDialogApp.Draw()
		Expect(connAddDialog.HasFocus()).To(Equal(true))
	})

	It("has focus", func() {
		connAddDialog.focusElement = connNameFieldFocus
		connAddDialogApp.SetFocus(connAddDialog)
		connAddDialogApp.Draw()
		Expect(connAddDialog.HasFocus()).To(Equal(true))

		connAddDialog.focusElement = connURIFieldFocus
		connAddDialogApp.SetFocus(connAddDialog)
		connAddDialogApp.Draw()
		Expect(connAddDialog.HasFocus()).To(Equal(true))

		connAddDialog.focusElement = connIdentityFieldFocus
		connAddDialogApp.SetFocus(connAddDialog)
		connAddDialogApp.Draw()
		Expect(connAddDialog.HasFocus()).To(Equal(true))
	})

	It("add button selected", func() {
		addWants := "add selected"
		addAction := "add init"

		addFunc := func() {
			addAction = addWants
		}

		connAddDialog.SetAddFunc(addFunc)
		connAddDialog.focusElement = connNameFieldFocus
		connAddDialogApp.SetFocus(connAddDialog)
		connAddDialogApp.Draw()
		connAddDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		connAddDialogApp.Draw()
		connAddDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		connAddDialogApp.Draw()
		connAddDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		connAddDialogApp.Draw()
		connAddDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		connAddDialogApp.Draw()
		connAddDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		connAddDialogApp.Draw()
		Expect(addWants).To(Equal(addAction))
	})

	It("hide", func() {
		connAddDialog.Hide()
		Expect(connAddDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		connAddDialogApp.Stop()
	})
})
