package images

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rs/zerolog/log"
)

func (img *Images) runCommand(cmd string) { //nolint:cyclop
	switch cmd {
	case "build":
		img.buildDialog.Display()
	case "diff":
		img.diff()
	case "history":
		img.history()
	case "import":
		img.importDialog.Display()
	case "inspect":
		img.inspect()
	case "prune": //nolint:goconst
		img.cprune()
	case "push":
		img.cpush()
	case "rm":
		img.rm()
	case "save":
		img.csave()
	case "search/pull":
		img.searchDialog.Display()
	case "tag":
		img.ctag()
	case "tree":
		img.tree()
	case "untag":
		img.cuntag()
	}
}

func (img *Images) displayError(title string, err error) {
	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	img.errorDialog.SetTitle(title)
	img.errorDialog.SetText(fmt.Sprintf("%v", err))
	img.errorDialog.Display()
}

func (img *Images) build() {
	img.buildDialog.Hide()

	opts, err := img.buildDialog.ImageBuildOptions()
	if err != nil {
		img.buildPrgDialog.Hide()
		img.displayError("IMAGE BUILD ERROR", err)

		return
	}

	if opts.BuildOptions.ContextDirectory == "" && len(opts.ContainerFiles) == 0 {
		img.displayError("IMAGE BUILD ERROR", errNoBuildDirOrCntFile)

		return
	}

	img.buildPrgDialog.Display()
	writer := img.buildPrgDialog.LogWriter()
	opts.BuildOptions.Out = writer
	opts.BuildOptions.Err = writer

	buildFunc := func() {
		report, err := images.Build(opts)

		img.buildPrgDialog.Hide()

		if err != nil {
			img.displayError("IMAGE BUILD ERROR", err)

			return
		}

		img.messageDialog.SetTitle("podman image build")
		img.messageDialog.SetText(dialogs.MessageImageInfo, report, "")
		img.messageDialog.Display()
	}

	go buildFunc()
}

func (img *Images) diff() {
	imageID, imageName := img.getSelectedItem()

	if imageID == "" {
		img.displayError("", errNoImageToDiff)

		return
	}

	img.progressDialog.SetTitle("image diff in progress")
	img.progressDialog.Display()

	diff := func() {
		data, err := images.Diff(imageID)

		img.progressDialog.Hide()

		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) DIFF ERROR", imageID)
			img.displayError(title, err)

			return
		}

		headerLabel := fmt.Sprintf("%12s (%s)", imageID, imageName)

		img.messageDialog.SetTitle("podman image diff")
		img.messageDialog.SetText(dialogs.MessageImageInfo, headerLabel, strings.Join(data, "\n"))
		img.messageDialog.DisplayFullSize()
	}

	go diff()
}

func (img *Images) history() {
	if img.selectedID == "" {
		img.displayError("", errNoImageToHistory)

		return
	}

	result, err := images.History(img.selectedID)
	if err != nil {
		title := fmt.Sprintf("IMAGE (%s) HISTORY ERROR", img.selectedID)
		img.displayError(title, err)
	}

	img.historyDialog.SetImageInfo(img.selectedID, img.selectedName)
	img.historyDialog.UpdateResults(result)
	img.historyDialog.Display()
}

func (img *Images) imageImport() {
	importOpts, err := img.importDialog.ImageImportOptions()
	if err != nil {
		img.displayError("IMAGE IMPORT ERROR", err)

		return
	}

	img.importDialog.Hide()
	img.progressDialog.SetTitle("image import in progress")
	img.progressDialog.Display()

	importFunc := func() {
		newImageID, err := images.Import(importOpts)

		img.progressDialog.Hide()

		if err != nil {
			img.displayError("IMAGE IMPORT ERROR", err)

			return
		}

		img.messageDialog.SetTitle("podman image import")
		img.messageDialog.SetText(dialogs.MessageImageInfo, newImageID, "")
		img.messageDialog.Display()
	}

	go importFunc()
}

func (img *Images) inspect() {
	imageID, imageName := img.getSelectedItem()
	if imageID == "" {
		img.displayError("", errNoImageToInspect)

		return
	}

	data, err := images.Inspect(imageID)
	if err != nil {
		title := fmt.Sprintf("IMAGE (%s) INSPECT ERROR", imageID)
		img.displayError(title, err)

		return
	}

	headerLabel := fmt.Sprintf("%12s (%s)", imageID, imageName)

	img.messageDialog.SetTitle("podman image inspect")
	img.messageDialog.SetText(dialogs.MessageImageInfo, headerLabel, data)
	img.messageDialog.DisplayFullSize()
}

func (img *Images) cprune() {
	img.confirmDialog.SetTitle("podman image prune")
	img.confirmData = "prune"
	img.confirmDialog.SetText("Are you sure you want to remove all unused images ?")
	img.confirmDialog.Display()
}

func (img *Images) prune() {
	img.progressDialog.SetTitle("image prune in progress")
	img.progressDialog.Display()

	prune := func() {
		err := images.Prune()

		img.progressDialog.Hide()

		if err != nil {
			img.displayError("IMAGE PRUNE ERROR", err)

			return
		}
	}

	go prune()
}

func (img *Images) cpush() {
	id, name := img.getSelectedItem()
	if id == "" {
		img.displayError("", errNoImageToPush)

		return
	}

	img.pushDialog.SetImageInfo(id, name)
	img.pushDialog.Display()
}

func (img *Images) push() {
	pushOptions := img.pushDialog.GetImagePushOptions()
	img.pushDialog.Hide()
	img.progressDialog.SetTitle("image push in progress")
	img.progressDialog.Display()

	push := func() {
		if err := images.Push(img.selectedID, pushOptions); err != nil {
			img.progressDialog.Hide()
			title := fmt.Sprintf("IMAGE (%s) PUSH ERROR", img.selectedID)
			img.displayError(title, err)

			return
		}

		img.progressDialog.Hide()
	}

	go push()
}

func (img *Images) rm() {
	imageID, imageName := img.getSelectedItem()
	if imageID == "" {
		img.displayError("", errNoImageToRemove)

		return
	}

	img.confirmDialog.SetTitle("podman image remove")
	img.confirmData = "rm"
	bgColor := style.GetColorHex(style.DialogBorderColor)
	fgColor := style.GetColorHex(style.DialogFgColor)
	imageItem := fmt.Sprintf("[%s:%s:b]IMAGE ID:[:-:-] %s (%s)", fgColor, bgColor, imageID, imageName)
	description := fmt.Sprintf("%s\n\nAre you sure you want to remove the selected image?", imageItem) //nolint:perfsprint

	img.confirmDialog.SetText(description)
	img.confirmDialog.Display()
}

func (img *Images) remove() {
	imageID, imageName := img.getSelectedItem()
	img.progressDialog.SetTitle("image remove in progress")
	img.progressDialog.Display()

	remove := func() {
		data, err := images.Remove(imageID)

		img.progressDialog.Hide()

		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) REMOVE ERROR", imageID)
			img.displayError(title, err)
		} else {
			headerLabel := fmt.Sprintf("%12s (%s)", imageID, imageName)

			img.messageDialog.SetTitle("podman image remove")
			img.messageDialog.SetText(dialogs.MessageImageInfo, headerLabel, strings.Join(data, "\n"))
			img.messageDialog.Display()
		}
	}

	go remove()
}

func (img *Images) csave() {
	id, name := img.getSelectedItem()
	if id == "" {
		img.displayError("", errNoImageToSave)

		return
	}

	img.saveDialog.SetImageInfo(id, name)
	img.saveDialog.Display()
}

func (img *Images) save() {
	saveOpts, err := img.saveDialog.ImageSaveOptions()
	if err != nil {
		title := fmt.Sprintf("IMAGE (%s) SAVE ERROR", img.selectedID)
		img.displayError(title, err)

		return
	}

	img.saveDialog.Hide()
	img.progressDialog.SetTitle("image save in progress")
	img.progressDialog.Display()

	saveFunc := func() {
		err := images.Save(img.selectedID, saveOpts)

		img.progressDialog.Hide()

		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) SAVE ERROR", img.selectedID)
			img.displayError(title, err)

			return
		}
	}

	go saveFunc()
}

func (img *Images) search(term string) {
	img.searchDialog.ClearResults()
	img.progressDialog.SetTitle("image search in progress")
	img.progressDialog.Display()

	search := func(term string) {
		result, err := images.Search(term)
		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) SEARCH ERROR", img.selectedID)
			img.displayError(title, err)
		}

		img.searchDialog.UpdateResults(result)
		img.progressDialog.Hide()
	}

	go search(term)
}

func (img *Images) ctag() {
	if img.selectedID == "" {
		img.displayError("", errNoImageToTag)

		return
	}

	img.cmdInputDialog.SetTitle("podman image tag")

	fgColor := style.GetColorHex(style.DialogFgColor)
	bgColor := style.GetColorHex(style.DialogBorderColor)

	description := fmt.Sprintf("[%s:%s:b]IMAGE ID:[:-:-] %s (%s)",
		fgColor, bgColor, img.selectedID, img.selectedName)

	img.cmdInputDialog.SetDescription(description)
	img.cmdInputDialog.SetSelectButtonLabel("tag")
	img.cmdInputDialog.SetLabel("target name")
	img.cmdInputDialog.SetSelectedFunc(func() {
		img.tag(img.cmdInputDialog.GetInputText())
		img.cmdInputDialog.Hide()
	})

	img.cmdInputDialog.Display()
}

func (img *Images) tag(tag string) {
	if err := images.Tag(img.selectedID, tag); err != nil {
		title := fmt.Sprintf("IMAGE (%s) TAG ERROR", img.selectedID)
		img.displayError(title, err)
	}
}

func (img *Images) cuntag() {
	if img.selectedID == "" {
		img.displayError("", errNoImageToUntag)

		return
	}

	img.cmdInputDialog.SetTitle("podman image untag")
	img.cmdInputDialog.SetDescription("")
	img.cmdInputDialog.SetSelectButtonLabel("untag")
	img.cmdInputDialog.SetLabel("image")
	img.cmdInputDialog.SetInputText(img.selectedName)
	img.cmdInputDialog.SetSelectedFunc(func() {
		img.untag(img.cmdInputDialog.GetInputText())
		img.cmdInputDialog.Hide()
	})

	img.cmdInputDialog.Display()
}

func (img *Images) tree() {
	imageID, imageName := img.getSelectedItem()
	if imageID == "" {
		img.displayError("", errNoImageToTree)

		return
	}

	retTree := func() {
		tree, err := images.Tree(imageID)
		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) TREE ERROR", imageID)
			img.displayError(title, err)

			return
		}

		headerLabel := fmt.Sprintf("%12s (%s)", imageID, imageName)

		img.progressDialog.Hide()
		img.messageDialog.SetTitle("podman image tree")
		img.messageDialog.SetText(dialogs.MessageImageInfo, headerLabel, tree)
		img.messageDialog.Display()
	}

	img.progressDialog.SetTitle("image tree in progress")
	img.progressDialog.Display()

	go retTree()
}

func (img *Images) untag(id string) {
	if err := images.Untag(id); err != nil {
		title := fmt.Sprintf("IMAGE (%s) UNTAG ERROR", img.selectedID)

		img.displayError(title, err)
	}
}

func (img *Images) pull(image string) {
	img.progressDialog.SetTitle("image pull in progress")

	img.progressDialog.Display()

	pull := func(name string) {
		err := images.Pull(name)
		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) PULL ERROR", img.selectedID)
			img.displayError(title, err)
		}

		img.progressDialog.Hide()
	}

	go pull(image)
}
