package imgdialogs

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

var _ = Describe("image push", Ordered, func() {
	var imagePushDialogApp *tview.Application
	var imagePushDialogScreen tcell.SimulationScreen
	var imagePushDialog *ImagePushDialog
	var runApp func()

	BeforeAll(func() {
		imagePushDialogApp = tview.NewApplication()
		imagePushDialog = NewImagePushDialog()
		imagePushDialogScreen = tcell.NewSimulationScreen("UTF-8")
		err := imagePushDialogScreen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := imagePushDialogApp.SetScreen(imagePushDialogScreen).SetRoot(imagePushDialog, true).Run(); err != nil {
				panic(err)
			}
		}
		zerolog.SetGlobalLevel(zerolog.Disabled)
		go runApp()
	})

	It("display", func() {
		imagePushDialog.Display()
		Expect(imagePushDialog.IsDisplay()).To(Equal(true))
	})

	It("set focus", func() {
		imagePushDialogApp.SetFocus(imagePushDialog)
		Expect(imagePushDialog.HasFocus()).To(Equal(true))
	})

	It("set image info", func() {
		imageID := "imageID"
		imageName := "imageName"
		imageInfoWants := fmt.Sprintf("%12s (%s)", imageID, imageName)
		imagePushDialog.SetImageInfo(imageID, imageName)
		Expect(imagePushDialog.imageInfo.GetText()).To(Equal(imageInfoWants))
	})

	It("cancel button selected", func() {
		cancelFunc := func() {
			imagePushDialog.Hide()
		}
		imagePushDialog.Hide()
		imagePushDialogApp.Draw()
		imagePushDialog.SetCancelFunc(cancelFunc)
		imagePushDialog.Display()
		imagePushDialogApp.Draw()
		imagePushDialogApp.SetFocus(imagePushDialog.form)
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		imagePushDialogApp.Draw()
		Expect(imagePushDialog.IsDisplay()).To(Equal(false))
	})

	It("push button selected", func() {
		pushButton := "initial"
		pushButtonWants := "push selected"
		commitFunc := func() {
			pushButton = pushButtonWants
		}
		imagePushDialog.Hide()
		imagePushDialogApp.Draw()
		imagePushDialog.SetPushFunc(commitFunc)
		imagePushDialog.Display()
		imagePushDialogApp.Draw()
		imagePushDialogApp.SetFocus(imagePushDialog.form)
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		imagePushDialogApp.Draw()
		Expect(pushButton).To(Equal(pushButtonWants))
	})

	It("set push handler", func() {
		push := "initial"
		pushWants := "push"
		pushHandler := func() {
			push = pushWants
		}
		imagePushDialog.SetPushFunc(pushHandler)
		imagePushDialog.pushHandler()
		Expect(push).To(Equal(pushWants))
	})

	It("set push handler", func() {
		cancel := "initial"
		cancelWants := "cancel"
		cancelHandler := func() {
			cancel = cancelWants
		}
		imagePushDialog.SetCancelFunc(cancelHandler)
		imagePushDialog.cancelHandler()
		Expect(cancel).To(Equal(cancelWants))
	})

	It("get push options", func() {
		opts := struct {
			Description   string
			Compress      bool
			Format        string
			SkipTLSVerify bool
			Username      string
			Password      string
			Authfile      string
		}{
			Description:   "a", // (256,97,0)
			Format:        "v2v2",
			Compress:      true,
			SkipTLSVerify: true,
			Username:      "c", // (256,99,0)
			Password:      "d", // (256,100,0)
			Authfile:      "e", // (256,101,0)
		}

		imagePushDialog.Hide()
		imagePushDialogApp.Draw()
		imagePushDialog.Display()
		//  description field
		imagePushDialogApp.SetFocus(imagePushDialog)
		imagePushDialogApp.Draw()
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(256, 97, tcell.ModNone))
		imagePushDialogApp.Draw()
		// compress field
		imagePushDialog.setFocusElement()
		imagePushDialogApp.SetFocus(imagePushDialog)
		imagePushDialogApp.Draw()
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(256, 32, tcell.ModNone)) // space
		imagePushDialogApp.Draw()
		// format dropdown
		imagePushDialog.setFocusElement()
		imagePushDialogApp.SetFocus(imagePushDialog)
		imagePushDialogApp.Draw()
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		imagePushDialogApp.Draw()
		// skip TLS verify field
		imagePushDialog.setFocusElement()
		imagePushDialogApp.SetFocus(imagePushDialog)
		imagePushDialogApp.Draw()
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(256, 32, tcell.ModNone)) // space
		imagePushDialogApp.Draw()
		// username field
		imagePushDialog.setFocusElement()
		imagePushDialogApp.SetFocus(imagePushDialog)
		imagePushDialogApp.Draw()
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone))
		imagePushDialogApp.Draw()
		// password field
		imagePushDialog.setFocusElement()
		imagePushDialogApp.SetFocus(imagePushDialog)
		imagePushDialogApp.Draw()
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(256, 100, tcell.ModNone))
		imagePushDialogApp.Draw()
		// authfile field
		imagePushDialog.setFocusElement()
		imagePushDialogApp.SetFocus(imagePushDialog)
		imagePushDialogApp.Draw()
		imagePushDialogApp.QueueEvent(tcell.NewEventKey(256, 101, tcell.ModNone))
		imagePushDialogApp.Draw()

		// get and check push options
		pushOptions := imagePushDialog.GetImagePushOptions()
		Expect(pushOptions.Destination).To(Equal(opts.Description))
		Expect(pushOptions.Format).To(Equal(opts.Format))
		Expect(pushOptions.Compress).To(Equal(opts.Compress))
		Expect(pushOptions.SkipTLSVerify).To(Equal(opts.SkipTLSVerify))
		Expect(pushOptions.Username).To(Equal(opts.Username))
		Expect(pushOptions.Password).To(Equal(opts.Password))
		Expect(pushOptions.AuthFile).To(Equal(opts.Authfile))
	})

	AfterAll(func() {
		imagePushDialogApp.Stop()
	})
})
