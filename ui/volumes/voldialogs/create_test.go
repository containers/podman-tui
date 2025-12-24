package voldialogs

import (
	"fmt"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("volume create", Ordered, func() {
	var volCreateDialogApp *tview.Application
	var volCreateDialogScreen tcell.SimulationScreen
	var volCreateDialog *VolumeCreateDialog
	var runApp func()

	BeforeAll(func() {
		volCreateDialogApp = tview.NewApplication()
		volCreateDialog = NewVolumeCreateDialog()
		volCreateDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := volCreateDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := volCreateDialogApp.SetScreen(volCreateDialogScreen).SetRoot(volCreateDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		volCreateDialog.Display()
		Expect(volCreateDialog.IsDisplay()).To(Equal(true))
		Expect(volCreateDialog.focusElement).To(Equal(volumeCreateNameFieldFocus))
	})

	It("initdata", func() {
		volCreateDialog.volumeNameField.SetText("sample")
		volCreateDialog.volumeLabelField.SetText("sample")
		volCreateDialog.volumeDriverField.SetText("sample")
		volCreateDialog.volumeDriverOptionsField.SetText("sample")
		volCreateDialog.volumeUIDField.SetText("sample")
		volCreateDialog.volumeGIDField.SetText("sample")

		volCreateDialog.initData()

		Expect(volCreateDialog.volumeNameField.GetText()).To(Equal(""))
		Expect(volCreateDialog.volumeLabelField.GetText()).To(Equal(""))
		Expect(volCreateDialog.volumeDriverField.GetText()).To(Equal(""))
		Expect(volCreateDialog.volumeDriverOptionsField.GetText()).To(Equal(""))
		Expect(volCreateDialog.volumeUIDField.GetText()).To(Equal(""))
		Expect(volCreateDialog.volumeGIDField.GetText()).To(Equal(""))
	})

	It("set focus", func() {
		volCreateDialogApp.SetFocus(volCreateDialog)
		Expect(volCreateDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		volCreateDialog.SetCancelFunc(cancelFunc)
		volCreateDialog.focusElement = volumeCreateFormFocus
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		volCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		volCreateDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("create button selected", func() {
		createWants := "create selected"
		createAction := "create init"
		createFunc := func() {
			createAction = createWants
		}

		volCreateDialog.SetCreateFunc(createFunc)
		volCreateDialogApp.SetFocus(volCreateDialog.form)
		volCreateDialogApp.Draw()
		volCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		volCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		volCreateDialogApp.Draw()

		Expect(createAction).To(Equal(createWants))
	})

	It("hide", func() {
		volCreateDialog.Hide()
		Expect(volCreateDialog.IsDisplay()).To(Equal(false))
	})

	It("next focus element", func() {
		volCreateDialog.Hide()
		volCreateDialogApp.Draw()
		volCreateDialog.Display()

		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		Expect(volCreateDialog.volumeNameField.HasFocus()).To(Equal(true))

		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		Expect(volCreateDialog.volumeLabelField.HasFocus()).To(Equal(true))

		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		Expect(volCreateDialog.volumeDriverField.HasFocus()).To(Equal(true))

		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		Expect(volCreateDialog.volumeDriverOptionsField.HasFocus()).To(Equal(true))

		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		Expect(volCreateDialog.volumeUIDField.HasFocus()).To(Equal(true))

		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		Expect(volCreateDialog.volumeGIDField.HasFocus()).To(Equal(true))
	})

	It("create options", func() {
		volName := "testvol"
		volUID := "1000"
		volGID := "1000"

		volLabel := struct {
			key   string
			value string
		}{key: "labelkey", value: "labelvalue"}

		volLabelStr := fmt.Sprintf("%s=%s", volLabel.key, volLabel.value)
		volDriver := "testdriver"
		volOption := struct {
			key   string
			value string
		}{key: "optionkey", value: "optionvalue"}

		volOptionStr := fmt.Sprintf("%s=%s", volOption.key, volOption.value)

		volCreateDialog.Hide()
		volCreateDialogApp.Draw()
		volCreateDialog.Display()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		// enter volume name
		volNameEvents := utils.StringToEventKey(volName)
		for i := 0; i < len(volNameEvents); i++ {
			volCreateDialogApp.QueueEvent(volNameEvents[i])
			volCreateDialogApp.SetFocus(volCreateDialog)
			volCreateDialogApp.Draw()
		}
		// enter volume label
		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		volLabelEvents := utils.StringToEventKey(volLabelStr)
		for i := 0; i < len(volLabelEvents); i++ {
			volCreateDialogApp.QueueEvent(volLabelEvents[i])
			volCreateDialogApp.SetFocus(volCreateDialog)
			volCreateDialogApp.Draw()
		}

		// enter volume driver
		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		volDriverEvents := utils.StringToEventKey(volDriver)
		for i := 0; i < len(volDriverEvents); i++ {
			volCreateDialogApp.QueueEvent(volDriverEvents[i])
			volCreateDialogApp.SetFocus(volCreateDialog)
			volCreateDialogApp.Draw()
		}

		// enter volume options
		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		volOptionEvents := utils.StringToEventKey(volOptionStr)
		for i := 0; i < len(volOptionEvents); i++ {
			volCreateDialogApp.QueueEvent(volOptionEvents[i])
			volCreateDialogApp.SetFocus(volCreateDialog)
			volCreateDialogApp.Draw()
		}

		// enter volume uid
		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		volUIDEvent := utils.StringToEventKey(volUID)
		for i := 0; i < len(volUIDEvent); i++ {
			volCreateDialogApp.QueueEvent(volUIDEvent[i])
			volCreateDialogApp.SetFocus(volCreateDialog)
			volCreateDialogApp.Draw()
		}

		// enter volume gid
		volCreateDialog.nextFocus()
		volCreateDialogApp.SetFocus(volCreateDialog)
		volCreateDialogApp.Draw()
		volGIDEvent := utils.StringToEventKey(volGID)
		for i := 0; i < len(volGIDEvent); i++ {
			volCreateDialogApp.QueueEvent(volGIDEvent[i])
			volCreateDialogApp.SetFocus(volCreateDialog)
			volCreateDialogApp.Draw()
		}

		volCreateOptions, err := volCreateDialog.VolumeCreateOptions()
		Expect(err).To(BeNil())
		Expect(volCreateOptions.Name).To(Equal(volName))
		volLabelValue := volCreateOptions.Labels[volLabel.key]
		Expect(volLabelValue).To(Equal(volLabel.value))
		Expect(volCreateOptions.Driver).To(Equal(volDriver))
		volOptionValue := volCreateOptions.DriverOptions[volOption.key]
		Expect(volOptionValue).To(Equal(volOption.value))

	})

	AfterAll(func() {
		volCreateDialogApp.Stop()
	})

})
