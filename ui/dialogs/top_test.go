package dialogs

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("top dialog", Ordered, func() {
	var app *tview.Application
	var screen tcell.SimulationScreen
	var topDialog *TopDialog
	var runApp func()

	BeforeAll(func() {
		app = tview.NewApplication()
		topDialog = NewTopDialog()
		screen = tcell.NewSimulationScreen("UTF-8")
		err := screen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := app.SetScreen(screen).SetRoot(topDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		topDialog.Display()
		Expect(topDialog.IsDisplay()).To(Equal(true))
	})

	It("set title", func() {
		title := "podman"
		topDialog.SetTitle(title)
		Expect(topDialog.layout.GetTitle()).To(Equal(strings.ToUpper(title)))
	})

	It("set focus", func() {
		app.SetFocus(topDialog)
		Expect(topDialog.HasFocus()).To(Equal(true))
	})

	It("enter button selected", func() {
		enterButton := "initial"
		enterButtonWants := "enter selected"
		enterFunc := func() {
			enterButton = enterButtonWants
		}
		topDialog.SetCancelFunc(enterFunc)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		Expect(enterButton).To(Equal(enterButtonWants))
	})

	It("hide", func() {
		topDialog.Hide()
		Expect(topDialog.IsDisplay()).To(Equal(false))
	})

	It("input handler", func() {
		topContent := [][]string{
			{"header", "header", "header", "header", "header", "header", "header", "header"},
			{"row01", "r01_pid", "r01_ppid", "r01_cpu", "r01_elapsed", "r01_tty", "r01_time", "r01_command"},
			{"row02", "r02_pid", "r02_ppid", "r02_cpu", "r02_elapsed", "r02_tty", "r02_time", "r02_command"},
			{"row03", "r03_pid", "r03_ppid", "r03_cpu", "r03_elapsed", "r03_tty", "r03_time", "r03_command"},
		}
		topDialog.Display()
		topDialog.UpdateResults(TopPodInfo, "", "", topContent)
		app.Draw()
		row := 1
		Expect(topDialog.table.GetCell(row, 0).Text).To(Equal(topContent[row][0]))
		app.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		app.Draw()
		currentRow, _ := topDialog.table.GetSelection()
		Expect(currentRow).To(Equal(row + 1))
		app.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		app.Draw()
		currentRow, _ = topDialog.table.GetSelection()
		Expect(currentRow).To(Equal(row + 2))
	})

	AfterAll(func() {
		app.Stop()
	})

})
