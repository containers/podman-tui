package netdialogs

import (
	"fmt"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("network connect", Ordered, func() {
	var netConnectDialogApp *tview.Application
	var netConnectDialogScreen tcell.SimulationScreen
	var netConnectDialog *NetworkConnectDialog
	var runApp func()

	BeforeAll(func() {
		netConnectDialogApp = tview.NewApplication()
		netConnectDialog = NewNetworkConnectDialog()
		netConnectDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := netConnectDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := netConnectDialogApp.SetScreen(netConnectDialogScreen).SetRoot(netConnectDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		netConnectDialog.Display()
		Expect(netConnectDialog.IsDisplay()).To(Equal(true))
		Expect(netConnectDialog.focusElement).To(Equal(netConnectContainerFocus))
	})

	It("set focus", func() {
		netConnectDialogApp.SetFocus(netConnectDialog)
		Expect(netConnectDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		netConnectDialog.SetCancelFunc(cancelFunc)
		netConnectDialog.focusElement = netConnectFormFocus
		netConnectDialogApp.SetFocus(netConnectDialog)
		netConnectDialogApp.Draw()
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netConnectDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("connect button selected", func() {
		connectWants := "connect selected"
		connectAction := "connect init"
		cancelFunc := func() {
			connectAction = connectWants
		}
		netConnectDialog.SetConnectFunc(cancelFunc)
		netConnectDialog.focusElement = netConnectFormFocus
		netConnectDialogApp.SetFocus(netConnectDialog)
		netConnectDialogApp.Draw()
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netConnectDialogApp.Draw()
		Expect(connectAction).To(Equal(connectWants))
	})

	It("set containers", func() {
		containerList := make([]entities.ListContainer, 0)
		containerList = append(containerList, entities.ListContainer{
			ID:    "f7db5ff00f23f7db5ff00f23",
			Names: []string{"container01"},
		})
		containerList = append(containerList, entities.ListContainer{
			ID:    "a92c29b48f32a92c29b48f32",
			Names: []string{"container02"},
		})

		netConnectDialog.SetContainers(containerList)
		_, cnt := netConnectDialog.container.GetCurrentOption()

		expectedContainer := fmt.Sprintf("%s (%s)", containerList[0].ID[0:12], containerList[0].Names[0])
		Expect(expectedContainer).To(Equal(cnt))
	})

	It("get connect options", func() {
		opts := struct {
			alias string
			ipv4  string
			ipv6  string
			mac   string
		}{
			alias: "a", // (256,97,0)
			ipv4:  "c", // (256,99,0)
			ipv6:  "d", // (256,100,0)
			mac:   "e", // (256,101,0)
		}
		containerList := make([]entities.ListContainer, 0)
		containerList = append(containerList, entities.ListContainer{
			ID:    "f7db5ff00f23f7db5ff00f23",
			Names: []string{"container01"},
		})
		containerList = append(containerList, entities.ListContainer{
			ID:    "a92c29b48f32a92c29b48f32",
			Names: []string{"container02"},
		})

		netConnectDialog.Hide()
		netConnectDialogApp.Draw()
		netConnectDialog.SetContainers(containerList)
		netConnectDialog.Display()
		netConnectDialogApp.SetFocus(netConnectDialog)
		netConnectDialogApp.Draw()

		// container
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netConnectDialogApp.Draw()

		// alias
		netConnectDialog.setFocusElement()
		netConnectDialogApp.SetFocus(netConnectDialog)
		netConnectDialogApp.Draw()
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(256, 97, tcell.ModNone))
		netConnectDialogApp.Draw()

		// ipv4
		netConnectDialog.setFocusElement()
		netConnectDialogApp.SetFocus(netConnectDialog)
		netConnectDialogApp.Draw()
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone))
		netConnectDialogApp.Draw()

		// ipv6
		netConnectDialog.setFocusElement()
		netConnectDialogApp.SetFocus(netConnectDialog)
		netConnectDialogApp.Draw()
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(256, 100, tcell.ModNone))
		netConnectDialogApp.Draw()

		// mac address
		netConnectDialog.setFocusElement()
		netConnectDialogApp.SetFocus(netConnectDialog)
		netConnectDialogApp.Draw()
		netConnectDialogApp.QueueEvent(tcell.NewEventKey(256, 101, tcell.ModNone))
		netConnectDialogApp.Draw()

		connectOptions := netConnectDialog.GetConnectOptions()
		Expect(connectOptions.Container).To(Equal(containerList[1].ID[0:12]))
		Expect(connectOptions.Aliases[0]).To(Equal(opts.alias))
		Expect(connectOptions.IPv4).To(Equal(opts.ipv4))
		Expect(connectOptions.IPv6).To(Equal(opts.ipv6))
		Expect(connectOptions.MacAddress).To(Equal(opts.mac))
	})

	It("hide", func() {
		netConnectDialog.Hide()
		Expect(netConnectDialog.IsDisplay()).To(Equal(false))
		Expect(netConnectDialog.aliases.GetText()).To(Equal(""))
		Expect(netConnectDialog.ipv4.GetText()).To(Equal(""))
		Expect(netConnectDialog.ipv6.GetText()).To(Equal(""))
		Expect(netConnectDialog.macAddr.GetText()).To(Equal(""))
	})

	AfterAll(func() {
		netConnectDialogApp.Stop()
	})
})
