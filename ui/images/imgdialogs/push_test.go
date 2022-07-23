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
	var app *tview.Application
	var screen tcell.SimulationScreen
	var imagePushDialog *ImagePushDialog
	var runApp func()

	BeforeAll(func() {
		app = tview.NewApplication()
		imagePushDialog = NewImagePushDialog()
		screen = tcell.NewSimulationScreen("UTF-8")
		err := screen.Init()
		if err != nil {
			panic(err)
		}
		runApp = func() {
			if err := app.SetScreen(screen).SetRoot(imagePushDialog, true).Run(); err != nil {
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
		app.SetFocus(imagePushDialog)
		Expect(imagePushDialog.HasFocus()).To(Equal(true))
	})

	It("set image info", func() {
		imageID := "imageID"
		imageName := "imageName"
		imageInfoWants := fmt.Sprintf("%-13s%s (%s)", "Image ID:", imageID, imageName)
		imagePushDialog.SetImageInfo(imageID, imageName)
		Expect(imagePushDialog.imageInfo.GetText(true)).To(Equal(imageInfoWants))
	})

	It("cancel button selected", func() {
		cancelFunc := func() {
			imagePushDialog.Hide()
		}
		imagePushDialog.Hide()
		app.Draw()
		imagePushDialog.SetCancelFunc(cancelFunc)
		imagePushDialog.Display()
		app.Draw()
		app.SetFocus(imagePushDialog.form)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		Expect(imagePushDialog.IsDisplay()).To(Equal(false))
	})

	It("push button selected", func() {
		pushButton := "initial"
		pushButtonWants := "push selected"
		commitFunc := func() {
			pushButton = pushButtonWants
		}
		imagePushDialog.Hide()
		app.Draw()
		imagePushDialog.SetPushFunc(commitFunc)
		imagePushDialog.Display()
		app.Draw()
		app.SetFocus(imagePushDialog.form)
		app.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
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
		app.Draw()
		imagePushDialog.Display()
		//  description field
		app.SetFocus(imagePushDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 97, tcell.ModNone))
		app.Draw()
		// compress field
		imagePushDialog.setFocusElement()
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 32, tcell.ModNone)) // space
		app.SetFocus(imagePushDialog)
		app.Draw()
		// format dropdown
		imagePushDialog.setFocusElement()
		app.SetFocus(imagePushDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		app.Draw()
		// skip TLS verify field
		imagePushDialog.setFocusElement()
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 32, tcell.ModNone)) // space
		app.SetFocus(imagePushDialog)
		app.Draw()
		// username field
		imagePushDialog.setFocusElement()
		app.SetFocus(imagePushDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 99, tcell.ModNone))
		app.Draw()
		// password field
		imagePushDialog.setFocusElement()
		app.SetFocus(imagePushDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 100, tcell.ModNone))
		app.Draw()
		// authfile field
		imagePushDialog.setFocusElement()
		app.SetFocus(imagePushDialog)
		app.Draw()
		app.QueueEvent(tcell.NewEventKey(256, 101, tcell.ModNone))
		app.Draw()

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
		app.Stop()
	})
})
