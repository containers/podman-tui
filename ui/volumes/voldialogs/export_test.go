package voldialogs

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("volume export", Ordered, func() {
	var volExportDialogApp *tview.Application
	var volExportDialogScreen tcell.SimulationScreen
	var volExportDialog *VolumeExportDialog
	var runApp func()

	BeforeAll(func() {
		volExportDialogApp = tview.NewApplication()
		volExportDialog = NewVolumeExportDialog()
		volExportDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := volExportDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := volExportDialogApp.SetScreen(volExportDialogScreen).SetRoot(volExportDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		volExportDialog.Display()
		Expect(volExportDialog.IsDisplay()).To(Equal(true))
		Expect(volExportDialog.focusElement).To(Equal(volumeExportOutputFieldFocus))
	})

	It("initdata", func() {
		volExportDialog.output.SetText("output")
		volExportDialog.initData()
		Expect(volExportDialog.output.GetText()).To(Equal(""))
		Expect(volExportDialog.focusElement).To(Equal(volumeExportOutputFieldFocus))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}

		volExportDialog.SetCancelFunc(cancelFunc)
		volExportDialog.focusElement = volumeExportFormFieldFocus
		volExportDialogApp.SetFocus(volExportDialog)
		volExportDialogApp.Draw()
		volExportDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		volExportDialogApp.Draw()
		volExportDialogApp.Draw()

		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("export button selected", func() {
		exportWants := "export selected"
		exportAction := "export init"
		exportFunc := func() {
			exportAction = exportWants
		}

		volExportDialog.SetExportFunc(exportFunc)
		volExportDialog.focusElement = volumeExportFormFieldFocus
		volExportDialogApp.SetFocus(volExportDialog)
		volExportDialogApp.Draw()
		volExportDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		volExportDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		volExportDialogApp.Draw()

		Expect(exportAction).To(Equal(exportWants))
	})

	It("export output", func() {
		output := "output.tar"
		volExportDialog.Hide()
		volExportDialog.Display()
		volExportDialogApp.SetFocus(volExportDialog)
		volExportDialogApp.Draw()

		outputEvents := utils.StringToEventKey(output)
		for i := 0; i < len(outputEvents); i++ {
			volExportDialogApp.QueueEvent(outputEvents[i])
			volExportDialogApp.SetFocus(volExportDialog)
			volExportDialogApp.Draw()
		}

		exportOutput := volExportDialog.VolumeExportOutput()

		Expect(exportOutput).To(Equal(output))
	})

	It("hide", func() {
		volExportDialog.Hide()
		Expect(volExportDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		volExportDialogApp.Stop()
	})

})
