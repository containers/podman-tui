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
		createDialog = NewContainerCreateDialog(ContainerCreateOnlyDialogMode)
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

		createDialog.focusElement = createcontainerPodFieldFocus
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

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		createDialog.SetCancelFunc(cancelFunc)
		createDialog.focusElement = createContainerFormFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		createDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("create button selected", func() {
		createWants := "create selected"
		createAction := "create init"
		createFunc := func() {
			createAction = createWants
		}
		createDialog.SetHandlerFunc(createFunc)
		createDialog.focusElement = createContainerFormFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTAB, 0, tcell.ModNone))
		createDialogApp.Draw()
		createDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		createDialogApp.Draw()
		Expect(createAction).To(Equal(createWants))
	})

	It("create options", func() {
		createDialog.focusElement = createContainerNameFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone)) // (256,99,0) c character
		createDialogApp.Draw()

		opts := createDialog.ContainerCreateOptions()
		Expect(opts.Name).To(Equal("c"))
	})

	It("setPortPageNextFocus", func() {
		createDialog.focusElement = createContainerPortPublishFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setPortPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerPortPublishAllFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setPortPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerPortExposeFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setPortPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerFormFocus))
	})

	It("setNetworkSettingsPageNextFocus", func() {
		createDialog.focusElement = createContainerHostnameFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setNetworkSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerIPAddrFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setNetworkSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerMacAddrFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setNetworkSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerNetworkFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setNetworkSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerFormFocus))
	})

	It("setSecurityOptionsPageNextFocus", func() {
		createDialog.focusElement = createcontainerSecLabelFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setSecurityOptionsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerApprarmorFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setSecurityOptionsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerSeccompFeildFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setSecurityOptionsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createcontainerSecMaskFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setSecurityOptionsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createcontainerSecUnmaskFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setSecurityOptionsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createcontainerSecNoNewPrivFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setSecurityOptionsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerFormFocus))
	})

	It("setUserGroupsPageNextFocus", func() {
		createDialog.focusElement = createContainerUserFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setUserGroupsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerHostUsersFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setUserGroupsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerPasswdEntryFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setUserGroupsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerGroupEntryFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setUserGroupsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerFormFocus))
	})

	It("setContainerInfoPageNextFocus", func() {
		createDialog.focusElement = createContainerNameFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCommandFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerImageFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createcontainerPodFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerLabelsFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerPrivilegedFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerRemoveFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerTimeoutFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerSecretFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setContainerInfoPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerFormFocus))
	})

	It("setEnvironmentPageNextFocus", func() {
		createDialog.focusElement = createContainerWorkDirFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setEnvironmentPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerEnvVarsFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setEnvironmentPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerEnvFileFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setEnvironmentPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerEnvMergeFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setEnvironmentPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerUnsetEnvFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setEnvironmentPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerEnvHostFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setEnvironmentPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerUnsetEnvAllFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setEnvironmentPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerUmaskFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setEnvironmentPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerFormFocus))
	})

	It("setResourceSettingsPageNextFocus", func() {
		createDialog.focusElement = createContainerMemoryFieldFocus
		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerMemoryReservatoinFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerMemorySwapFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createcontainerMemorySwappinessFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCPUsFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCPUSharesFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCPUPeriodFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCPURtPeriodFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCPUQuotaFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCPURtRuntimeFeildFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCPUSetCPUsFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerCPUSetMemsFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerShmSizeFieldFocus))

		createDialogApp.SetFocus(createDialog)
		createDialogApp.Draw()
		createDialog.setResourceSettingsPageNextFocus()
		Expect(createDialog.focusElement).To(Equal(createContainerShmSizeSystemdFieldFocus))
	})

	It("hide", func() {
		createDialog.Hide()
		Expect(createDialog.IsDisplay()).To(Equal(false))
	})

	AfterAll(func() {
		createDialogApp.Stop()
	})
})
