package netdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("network create", Ordered, func() {
	var netCreateDialogApp *tview.Application
	var netCreateDialogScreen tcell.SimulationScreen
	var netCreateDialog *NetworkCreateDialog
	var runApp func()

	BeforeAll(func() {
		netCreateDialogApp = tview.NewApplication()
		netCreateDialog = NewNetworkCreateDialog()
		netCreateDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := netCreateDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := netCreateDialogApp.SetScreen(netCreateDialogScreen).SetRoot(netCreateDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		netCreateDialog.Display()
		Expect(netCreateDialog.IsDisplay()).To(Equal(true))
		Expect(netCreateDialog.focusElement).To(Equal(categoryPagesFocus))
	})

	It("init data", func() {
		netCreateDialog.networkNameField.SetText("sample")
		netCreateDialog.networkLabelsField.SetText("sample")
		netCreateDialog.networkInternalCheckBox.SetChecked(true)
		netCreateDialog.networkDriverField.SetText("sample")
		netCreateDialog.networkDriverOptionsField.SetText("sample")
		netCreateDialog.networkIpv6CheckBox.SetChecked(true)
		netCreateDialog.networkGatewayField.SetText("sample")
		netCreateDialog.networkIPRangeField.SetText("sample")
		netCreateDialog.networkSubnetField.SetText("sample")
		netCreateDialog.networkDisableDNSCheckBox.SetChecked(true)

		netCreateDialog.initData()

		Expect(netCreateDialog.networkNameField.GetText()).To(Equal(""))
		Expect(netCreateDialog.networkLabelsField.GetText()).To(Equal(""))
		Expect(netCreateDialog.networkInternalCheckBox.IsChecked()).To(Equal(false))
		Expect(netCreateDialog.networkDriverField.GetText()).To(Equal(networks.DefaultNetworkDriver()))
		Expect(netCreateDialog.networkDriverOptionsField.GetText()).To(Equal(""))
		Expect(netCreateDialog.networkIpv6CheckBox.IsChecked()).To(Equal(false))
		Expect(netCreateDialog.networkGatewayField.GetText()).To(Equal(""))
		Expect(netCreateDialog.networkIPRangeField.GetText()).To(Equal(""))
		Expect(netCreateDialog.networkSubnetField.GetText()).To(Equal(""))
		Expect(netCreateDialog.networkDisableDNSCheckBox.IsChecked()).To(Equal(false))
	})

	It("set focus", func() {
		netCreateDialogApp.SetFocus(netCreateDialog)
		Expect(netCreateDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		netCreateDialog.SetCancelFunc(cancelFunc)
		netCreateDialog.focusElement = formFocus
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netCreateDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("create button selected", func() {
		createWants := "create selected"
		createAction := "create init"
		cancelFunc := func() {
			createAction = createWants
		}
		netCreateDialog.SetCreateFunc(cancelFunc)
		netCreateDialog.focusElement = formFocus
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		netCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netCreateDialogApp.Draw()
		Expect(createAction).To(Equal(createWants))
	})

	It("hide", func() {
		netCreateDialog.Hide()
		Expect(netCreateDialog.IsDisplay()).To(Equal(false))
	})

	It("basic info page next focus", func() {
		netCreateDialog.focusElement = networkNameFieldFocus
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkLabelFieldFocus))

		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkInternalCheckBoxFocus))

		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkDriverFieldFocus))

		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkDriverOptionsFieldFocus))

		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(formFocus))

	})

	It("ip settings page next focus", func() {
		netCreateDialog.setActiveCategory(ipSettingsPageIndex)
		netCreateDialog.focusElement = networkIPv6CheckBoxFocus
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkGatewatFieldFocus))

		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkIPRangeFieldFocus))

		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkSubnetFieldFocus))

		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkDisableDNSCheckBoxFocus))

		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(formFocus))

	})

	It("next category", func() {
		netCreateDialog.Hide()
		netCreateDialog.Display()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.nextCategory()
		Expect(netCreateDialog.activePageIndex).To(Equal(ipSettingsPageIndex))
		netCreateDialog.nextCategory()
		Expect(netCreateDialog.activePageIndex).To(Equal(basicInfoPageIndex))
	})

	It("previous category", func() {
		netCreateDialog.Hide()
		netCreateDialog.Display()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.previousCategory()
		Expect(netCreateDialog.activePageIndex).To(Equal(ipSettingsPageIndex))
		netCreateDialog.previousCategory()
		Expect(netCreateDialog.activePageIndex).To(Equal(basicInfoPageIndex))
	})

	It("create options", func() {
		netName := "testnet"
		netLabel := struct {
			key   string
			value string
		}{key: "labelkey", value: "labelvalue"}
		netLabelStr := fmt.Sprintf("%s=%s", netLabel.key, netLabel.value)

		netInternal := true
		netDriver := networks.DefaultNetworkDriver()
		netOption := struct {
			key   string
			value string
		}{key: "optionkey", value: "optionvalue"}
		netOptionStr := fmt.Sprintf("%s=%s", netOption.key, netOption.value)

		ipv6 := true
		gateway := "192.168.1.254"
		iprange := "192.168.1.10 192.168.1.20"
		subnet := "255.255.255.0"
		disableDNS := true

		// set network name
		netCreateDialog.Hide()
		netCreateDialog.Display()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netNameEvents := utils.StringToEventKey(netName)
		for i := 0; i < len(netNameEvents); i++ {
			netCreateDialogApp.QueueEvent(netNameEvents[i])
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		// set network labels
		netCreateDialog.setBasicInfoPageNextFocus()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netLabelEvents := utils.StringToEventKey(netLabelStr)
		for i := 0; i < len(netLabelEvents); i++ {
			netCreateDialogApp.QueueEvent(netLabelEvents[i])
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		// set internal
		netCreateDialog.setBasicInfoPageNextFocus()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		if netInternal {
			netCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		// set network options
		netCreateDialog.setBasicInfoPageNextFocus()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netOptionEvents := utils.StringToEventKey(netOptionStr)
		for i := 0; i < len(netOptionStr); i++ {
			netCreateDialogApp.QueueEvent(netOptionEvents[i])
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		// switch to ip settings page
		netCreateDialog.nextCategory()
		netCreateDialog.focusElement = networkIPv6CheckBoxFocus
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()

		// set IPv6
		if ipv6 {
			netCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		// set gateway
		netCreateDialog.setIPSettingsPageNextFocus()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netGatwayEvents := utils.StringToEventKey(gateway)
		for i := 0; i < len(netGatwayEvents); i++ {
			netCreateDialogApp.QueueEvent(netGatwayEvents[i])
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		// set ip range
		netCreateDialog.setIPSettingsPageNextFocus()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netIPRangeEvents := utils.StringToEventKey(iprange)
		for i := 0; i < len(netIPRangeEvents); i++ {
			netCreateDialogApp.QueueEvent(netIPRangeEvents[i])
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		// set sebnet
		netCreateDialog.setIPSettingsPageNextFocus()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		netSubnetEvents := utils.StringToEventKey(subnet)
		for i := 0; i < len(netSubnetEvents); i++ {
			netCreateDialogApp.QueueEvent(netSubnetEvents[i])
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		// set disable DNS
		netCreateDialog.setIPSettingsPageNextFocus()
		netCreateDialogApp.SetFocus(netCreateDialog)
		netCreateDialogApp.Draw()
		if disableDNS {
			netCreateDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			netCreateDialogApp.SetFocus(netCreateDialog)
			netCreateDialogApp.Draw()
		}

		networkCreateOptions := netCreateDialog.NetworkCreateOptions()
		Expect(networkCreateOptions.Name).To(Equal(netName))
		netLabelValue := networkCreateOptions.Labels[netLabel.key]
		Expect(netLabelValue).To(Equal(netLabel.value))
		Expect(networkCreateOptions.Internal).To(Equal(netInternal))
		Expect(networkCreateOptions.Drivers).To(Equal(netDriver))
		netOptionValue := networkCreateOptions.DriversOptions[netOption.key]
		Expect(netOptionValue).To(Equal(netOption.value))
		Expect(networkCreateOptions.IPv6).To(Equal(ipv6))
		Expect(networkCreateOptions.Gateways[0]).To(Equal(gateway))
		iprangeSplited := strings.Split(iprange, " ")
		for i := 0; i < len(networkCreateOptions.IPRanges); i++ {
			Expect(networkCreateOptions.IPRanges[i]).To(Equal(iprangeSplited[i]))
		}
		Expect(networkCreateOptions.Subnets[0]).To(Equal(subnet))
		Expect(networkCreateOptions.DisableDNS).To(Equal(disableDNS))
	})

	AfterAll(func() {
		netCreateDialogApp.Stop()
	})

})
