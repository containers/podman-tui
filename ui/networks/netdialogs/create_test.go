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
	var netDialogApp *tview.Application
	var netDialogScreen tcell.SimulationScreen
	var netCreateDialog *NetworkCreateDialog
	var runApp func()

	BeforeAll(func() {
		netDialogApp = tview.NewApplication()
		netCreateDialog = NewNetworkCreateDialog()
		netDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := netDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := netDialogApp.SetScreen(netDialogScreen).SetRoot(netCreateDialog, true).Run(); err != nil {
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
		netDialogApp.SetFocus(netCreateDialog)
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
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netDialogApp.Draw()
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
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		netDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netDialogApp.Draw()
		Expect(createAction).To(Equal(createWants))
	})

	It("hide", func() {
		netCreateDialog.Hide()
		Expect(netCreateDialog.IsDisplay()).To(Equal(false))
	})

	It("basic info page next focus", func() {
		netCreateDialog.focusElement = networkNameFieldFocus
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkLabelFieldFocus))

		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkInternalCheckBoxFocus))

		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkDriverFieldFocus))

		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkDriverOptionsFieldFocus))

		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(formFocus))

	})

	It("ip settings page next focus", func() {
		netCreateDialog.setActiveCategory(ipSettingsPageIndex)
		netCreateDialog.focusElement = networkIPv6CheckBoxFocus
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkGatewatFieldFocus))

		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkIPRangeFieldFocus))

		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkSubnetFieldFocus))

		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(networkDisableDNSCheckBoxFocus))

		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setIPSettingsPageNextFocus()
		Expect(netCreateDialog.focusElement).To(Equal(formFocus))

	})

	It("next category", func() {
		netCreateDialog.Hide()
		netCreateDialog.Display()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.nextCategory()
		Expect(netCreateDialog.activePageIndex).To(Equal(ipSettingsPageIndex))
		netCreateDialog.nextCategory()
		Expect(netCreateDialog.activePageIndex).To(Equal(basicInfoPageIndex))
	})

	It("previous category", func() {
		netCreateDialog.Hide()
		netCreateDialog.Display()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
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
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netNameEvents := utils.StringToEventKey(netName)
		for i := 0; i < len(netNameEvents); i++ {
			netDialogApp.QueueEvent(netNameEvents[i])
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
		}

		// set network labels
		netCreateDialog.setBasicInfoPageNextFocus()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netLabelEvents := utils.StringToEventKey(netLabelStr)
		for i := 0; i < len(netLabelEvents); i++ {
			netDialogApp.QueueEvent(netLabelEvents[i])
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
		}

		// set internal
		netCreateDialog.setBasicInfoPageNextFocus()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		if netInternal {
			netDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
		}

		// set network options
		netCreateDialog.setBasicInfoPageNextFocus()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netCreateDialog.setBasicInfoPageNextFocus()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netOptionEvents := utils.StringToEventKey(netOptionStr)
		for i := 0; i < len(netOptionStr); i++ {
			netDialogApp.QueueEvent(netOptionEvents[i])
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
		}

		// switch to ip settings page
		netCreateDialog.nextCategory()
		netCreateDialog.focusElement = networkIPv6CheckBoxFocus
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()

		// set IPv6
		if ipv6 {
			netDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
		}

		// set gateway
		netCreateDialog.setIPSettingsPageNextFocus()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netGatwayEvents := utils.StringToEventKey(gateway)
		for i := 0; i < len(netGatwayEvents); i++ {
			netDialogApp.QueueEvent(netGatwayEvents[i])
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
		}

		// set ip range
		netCreateDialog.setIPSettingsPageNextFocus()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netIPRangeEvents := utils.StringToEventKey(iprange)
		for i := 0; i < len(netIPRangeEvents); i++ {
			netDialogApp.QueueEvent(netIPRangeEvents[i])
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
		}

		// set sebnet
		netCreateDialog.setIPSettingsPageNextFocus()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		netSubnetEvents := utils.StringToEventKey(subnet)
		for i := 0; i < len(netSubnetEvents); i++ {
			netDialogApp.QueueEvent(netSubnetEvents[i])
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
		}

		// set disable DNS
		netCreateDialog.setIPSettingsPageNextFocus()
		netDialogApp.SetFocus(netCreateDialog)
		netDialogApp.Draw()
		if disableDNS {
			netDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			netDialogApp.SetFocus(netCreateDialog)
			netDialogApp.Draw()
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
		netDialogApp.Stop()
	})

})
