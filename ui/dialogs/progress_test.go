package dialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("progress bar dialog", Ordered, func() {
	var app *tview.Application
	var screen tcell.SimulationScreen
	var progressDialog *ProgressDialog
	var runApp func()

	BeforeAll(func() {
		app = tview.NewApplication()
		progressDialog = NewProgressDialog()
		screen = tcell.NewSimulationScreen("UTF-8")
		err := screen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := app.SetScreen(screen).SetRoot(progressDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		progressDialog.Display()
		Expect(progressDialog.IsDisplay()).To(Equal(true))
		Expect(progressDialog.counterValue).To(Equal(0))
	})

	It("set focus", func() {
		app.SetFocus(progressDialog.Box)
		Expect(progressDialog.HasFocus()).To(Equal(true))
	})

	It("set title", func() {
		title := "podman"
		progressDialog.SetTitle(title)
		Expect(progressDialog.Box.GetTitle()).To(Equal(title))
	})

	It("set rect", func() {
		x := 0
		y := 0
		width := 50
		height := 20
		// wants
		wWants := prgMinWidth
		hWants := 3 // progress bar has fixed height 3 rows
		spaceWidth := (width - wWants) / 2
		spaceHeight := (height - hWants) / 2
		xWants := x + spaceWidth
		yWants := y + spaceHeight

		// set rects
		progressDialog.SetRect(x, y, width, height)
		// get rects
		x1, y1, w1, h1 := progressDialog.Box.GetRect()
		Expect(x1).To(Equal(xWants))
		Expect(y1).To(Equal(yWants))
		Expect(w1).To(Equal(wWants))
		Expect(h1).To(Equal(hWants))
	})

	It("progress value", func() {
		maxValue := 10
		progressDialog.counterValue = 0
		progressDialog.tickStr(maxValue)
		Expect(progressDialog.counterValue).To(Equal(1))
		progressDialog.tickStr(maxValue)
		Expect(progressDialog.counterValue).To(Equal(2))
		for i := 2; i <= maxValue-4; i++ {
			progressDialog.tickStr(maxValue)
		}
		Expect(progressDialog.counterValue).To(Equal(0))
	})

	It("hide", func() {
		progressDialog.Hide()
		Expect(progressDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		app.Stop()
	})

})
