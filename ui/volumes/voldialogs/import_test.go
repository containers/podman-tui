package voldialogs

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("volume import", Ordered, func() {
	var volImportDialogApp *tview.Application
	var volImportDialogScreen tcell.SimulationScreen
	var volImportDialog *VolumeImportDialog
	var runApp func()

	BeforeAll(func() {
		volImportDialogApp = tview.NewApplication()
		volImportDialog = NewVolumeImportDialog()
		volImportDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := volImportDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := volImportDialogApp.SetScreen(volImportDialogScreen).SetRoot(volImportDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		volImportDialog.Display()
		Expect(volImportDialog.IsDisplay()).To(Equal(true))
		Expect(volImportDialog.focusElement).To(Equal(volumeImportSourceFieldFocus))
	})

	It("initdata", func() {
		volImportDialog.source.SetText("source")
		volImportDialog.volume.SetText("volume")

		volImportDialog.initData()

		Expect(volImportDialog.source.GetText()).To(Equal(""))
		Expect(volImportDialog.volume.GetText()).To(Equal(""))
		Expect(volImportDialog.focusElement).To(Equal(volumeImportSourceFieldFocus))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}

		volImportDialog.SetCancelFunc(cancelFunc)
		volImportDialog.focusElement = volumeImportFormFieldFocus
		volImportDialogApp.SetFocus(volImportDialog)
		volImportDialogApp.Draw()
		volImportDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		volImportDialogApp.Draw()
		volImportDialogApp.Draw()

		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("import button selected", func() {
		importWants := "import selected"
		importAction := "import init"
		importFunc := func() {
			importAction = importWants
		}

		volImportDialog.SetImportFunc(importFunc)
		volImportDialog.focusElement = volumeImportFormFieldFocus
		volImportDialogApp.SetFocus(volImportDialog)
		volImportDialogApp.Draw()
		volImportDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		volImportDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		volImportDialogApp.Draw()

		Expect(importAction).To(Equal(importWants))
	})

	It("import source", func() {
		source := "source.tar"
		volImportDialog.Hide()
		volImportDialog.Display()
		volImportDialogApp.SetFocus(volImportDialog)
		volImportDialogApp.Draw()

		outputEvents := utils.StringToEventKey(source)
		for i := 0; i < len(outputEvents); i++ {
			volImportDialogApp.QueueEvent(outputEvents[i])
			volImportDialogApp.SetFocus(volImportDialog)
			volImportDialogApp.Draw()
		}

		importSource := volImportDialog.VolumeImportSource()

		Expect(importSource).To(Equal(source))
	})

	It("hide", func() {
		volImportDialog.Hide()
		Expect(volImportDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		volImportDialogApp.Stop()
	})

})
