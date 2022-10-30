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

var _ = Describe("network disconnect", Ordered, func() {
	var netDisconnectDialogApp *tview.Application
	var netDialogScreen tcell.SimulationScreen
	var netDisconnectDialog *NetworkDisconnectDialog
	var runApp func()

	BeforeAll(func() {
		netDisconnectDialogApp = tview.NewApplication()
		netDisconnectDialog = NewNetworkDisconnectDialog()
		netDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := netDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := netDisconnectDialogApp.SetScreen(netDialogScreen).SetRoot(netDisconnectDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		netDisconnectDialog.Display()
		Expect(netDisconnectDialog.IsDisplay()).To(Equal(true))
		Expect(netDisconnectDialog.focusElement).To(Equal(netDisconnectContainerFocus))
	})

	It("set focus", func() {
		netDisconnectDialogApp.SetFocus(netDisconnectDialog)
		Expect(netDisconnectDialog.HasFocus()).To(Equal(true))
	})

	It("cancel button selected", func() {
		cancelWants := "cancel selected"
		cancelAction := "cancel init"
		cancelFunc := func() {
			cancelAction = cancelWants
		}
		netDisconnectDialog.SetCancelFunc(cancelFunc)
		netDisconnectDialog.focusElement = netConnectFormFocus
		netDisconnectDialogApp.SetFocus(netDisconnectDialog)
		netDisconnectDialogApp.Draw()
		netDisconnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netDisconnectDialogApp.Draw()
		Expect(cancelAction).To(Equal(cancelWants))
	})

	It("disconnect button selected", func() {
		disconnectWants := "disconnect selected"
		disconnectAction := "disconnect init"
		cancelFunc := func() {
			disconnectAction = disconnectWants
		}
		netDisconnectDialog.SetDisconnectFunc(cancelFunc)
		netDisconnectDialog.focusElement = netConnectFormFocus
		netDisconnectDialogApp.SetFocus(netDisconnectDialog)
		netDisconnectDialogApp.Draw()
		netDisconnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		netDisconnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netDisconnectDialogApp.Draw()
		Expect(disconnectWants).To(Equal(disconnectAction))
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

		netDisconnectDialog.SetContainers(containerList)
		_, cnt := netDisconnectDialog.container.GetCurrentOption()

		expectedContainer := fmt.Sprintf("%s (%s)", containerList[0].ID[0:12], containerList[0].Names[0])
		Expect(expectedContainer).To(Equal(cnt))
	})

	It("get disconnect options", func() {

		network := "network01"
		containerList := make([]entities.ListContainer, 0)
		containerList = append(containerList, entities.ListContainer{
			ID:    "f7db5ff00f23f7db5ff00f23",
			Names: []string{"container01"},
		})
		containerList = append(containerList, entities.ListContainer{
			ID:    "a92c29b48f32a92c29b48f32",
			Names: []string{"container02"},
		})

		netDisconnectDialog.Hide()
		netDisconnectDialogApp.Draw()
		netDisconnectDialog.SetContainers(containerList)
		netDisconnectDialog.SetNetworkInfo(network)
		netDisconnectDialog.Display()
		netDisconnectDialogApp.SetFocus(netDisconnectDialog)
		netDisconnectDialogApp.Draw()

		// container
		netDisconnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netDisconnectDialogApp.Draw()
		netDisconnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		netDisconnectDialogApp.Draw()
		netDisconnectDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		netDisconnectDialogApp.Draw()

		networkName, containerID := netDisconnectDialog.GetDisconnectOptions()
		Expect(containerID).To(Equal(containerList[1].ID[0:12]))
		Expect(networkName).To(Equal(network))
	})

	It("hide", func() {
		networkInfo := fmt.Sprintf("%-11s", "Network:")

		netDisconnectDialog.Hide()
		Expect(netDisconnectDialog.IsDisplay()).To(Equal(false))
		Expect(netDisconnectDialog.network.GetText(true)).To(Equal(networkInfo))
	})

	AfterAll(func() {
		netDisconnectDialogApp.Stop()
	})
})
