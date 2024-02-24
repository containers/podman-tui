package cntdialogs

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("container create", Ordered, func() {
	var createDialogApp *tview.Application
	var createDialogScreen tcell.SimulationScreen
	var createDialog *ContainerCreateDialog
	var runApp func()

	BeforeAll(func() {
		createDialogApp = tview.NewApplication()
		createDialog = NewContainerCreateDialog()
		createDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := createDialogScreen.Init()
		if err != nil {
			panic(err)
		}

		runApp = func() {
			if err := createDialogApp.SetScreen(createDialogScreen).SetRoot(createDialog, true).Run(); err != nil {
				panic(err)
			}
		}

		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		createDialog.Display()
		createDialogApp.Draw()
		Expect(createDialog.IsDisplay()).To(Equal(true))
		Expect(createDialog.focusElement).To(Equal(createCategoryPagesFocus))
	})

	It("set focus", func() {
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		Expect(createDialog.HasFocus()).To(Equal(true))
	})

	It("dropdown has focus", func() {
		createDialog.focusElement = createContainerImageFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		Expect(createDialog.dropdownHasFocus()).To(Equal(true))

		createDialog.focusElement = createcontainerPodFieldFocis
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		Expect(createDialog.dropdownHasFocus()).To(Equal(true))

		createDialog.focusElement = createContainerNetworkFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		Expect(createDialog.dropdownHasFocus()).To(Equal(true))

		createDialog.focusElement = createContainerImageVolumeFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		Expect(createDialog.dropdownHasFocus()).To(Equal(true))
	})

	It("hide", func() {
		createDialog.Hide()
		Expect(createDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		createDialogApp.Stop()
	})
})
