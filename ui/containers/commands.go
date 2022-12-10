package containers

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	bcontainers "github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

func (cnt *Containers) runCommand(cmd string) {
	switch cmd {
	case "checkpoint":
		cnt.preCheckpoint()
	case "commit":
		cnt.preCommit()
	case "create":
		cnt.createDialog.Display()
	case "diff":
		cnt.diff()
	case "exec":
		cnt.cexec()
	case "inspect":
		cnt.inspect()
	case "kill":
		cnt.kill()
	case "logs":
		cnt.logs()
	case "pause":
		cnt.pause()
	case "prune":
		cnt.cprune()
	case "rename":
		cnt.rename()
	case "restore":
		cnt.preRestore()
	case "port":
		cnt.port()
	case "rm":
		cnt.rm()
	case "start":
		cnt.start()
	case "stats":
		cnt.stats()
	case "stop":
		cnt.stop()
	case "top":
		cnt.top()
	case "unpause":
		cnt.unpause()
	}
}

func (cnt *Containers) displayError(title string, err error) {
	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	cnt.errorDialog.SetTitle(title)
	cnt.errorDialog.SetText(fmt.Sprintf("%v", err))
	cnt.errorDialog.Display()
}

func (cnt *Containers) preRestore() {
	var (
		containersList [][]string
		podsList       [][]string
	)

	cnt.progressDialog.SetTitle("operation in progress")
	cnt.progressDialog.Display()

	// get current containers
	cntList, err := containers.List()
	if err != nil {
		cnt.progressDialog.Hide()
		cnt.displayError("CONTAINER RESTORE ERROR", err)
		return
	}

	for _, cnt := range cntList {
		containersList = append(containersList, []string{
			cnt.ID,
			cnt.Names[0],
		})
	}

	cnt.restoreDialog.SetContainers(containersList)

	// get current pods
	podList, err := pods.List()
	if err != nil {
		cnt.progressDialog.Hide()
		cnt.displayError("CONTAINER RESTORE ERROR", err)
		return
	}

	for _, pod := range podList {
		podsList = append(podsList, []string{
			pod.Id,
			pod.Name,
		})
	}

	cnt.restoreDialog.SetPods(podsList)

	cnt.progressDialog.Hide()
	cnt.restoreDialog.Display()
}

func (cnt *Containers) restore() {
	restoreOptions := cnt.restoreDialog.GetRestoreOptions()

	cnt.restoreDialog.Hide()
	cnt.progressDialog.SetTitle("container restore in progress")
	cnt.progressDialog.Display()

	restore := func() {
		report, err := containers.Restore(restoreOptions)
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) RESTORE ERROR", cnt.selectedID)

			cnt.progressDialog.Hide()
			cnt.displayError(title, err)

			return
		}

		headerLabel := fmt.Sprintf("%s (%s)", restoreOptions.ContainerID, restoreOptions.Name)
		cnt.progressDialog.Hide()
		cnt.messageDialog.SetTitle("podman container restore")
		cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, report)
		cnt.messageDialog.Display()
	}

	go restore()
}

func (cnt *Containers) preCheckpoint() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to perform checkpoint command"))
		return
	}
	cntID, cntName := cnt.getSelectedItem()
	cnt.checkpointDialog.SetContainerInfo(cntID, cntName)
	cnt.checkpointDialog.Display()
}

func (cnt *Containers) checkpoint() {
	checkpointOptions := cnt.checkpointDialog.GetCheckpointOptions()

	cnt.checkpointDialog.Hide()
	cnt.progressDialog.SetTitle("container checkpoint in progress")
	cnt.progressDialog.Display()

	checkpoint := func() {
		report, err := containers.Checkpoint(cnt.selectedID, checkpointOptions)
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) CHECKPOINT ERROR", cnt.selectedID)

			cnt.progressDialog.Hide()
			cnt.displayError(title, err)

			return
		}
		headerLabel := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)

		cnt.progressDialog.Hide()
		cnt.messageDialog.SetTitle("podman container checkpoint")
		cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, report)
		cnt.messageDialog.Display()
	}

	go checkpoint()
}

func (cnt *Containers) preCommit() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to perform commit command"))
		return
	}
	cntID, cntName := cnt.getSelectedItem()
	cnt.commitDialog.SetContainerInfo(cntID, cntName)
	cnt.commitDialog.Display()
}

func (cnt *Containers) commit() {
	commitOpts := cnt.commitDialog.GetContainerCommitOptions()
	cnt.commitDialog.Hide()
	cnt.progressDialog.SetTitle("container commit in progress")
	cnt.progressDialog.Display()
	cntCommit := func() {
		response, err := containers.Commit(cnt.selectedID, commitOpts)
		if err != nil {
			cnt.progressDialog.Hide()
			title := fmt.Sprintf("CONTAINER (%s) COMMIT ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}

		headerLabel := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)

		cnt.progressDialog.Hide()
		cnt.messageDialog.SetTitle("podman container commit")
		cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, response)
		cnt.messageDialog.Display()
	}
	go cntCommit()
}

func (cnt *Containers) stats() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to perform stats command"))
		return
	}
	cntID, cntName := cnt.getSelectedItem()
	cntStatus, err := containers.Status(cntID)
	if err != nil {
		cnt.displayError("", fmt.Errorf("there is no container to perform stats command"))
		return
	}
	if cntStatus != "running" {
		cnt.displayError("", fmt.Errorf("container (%s) status improper", cntID))
		return
	}
	stream := true
	statOption := new(bcontainers.StatsOptions)
	statOption.Stream = &stream
	statsChan, err := containers.Stats(cntID, statOption)
	if err != nil {
		cnt.displayError("CONTAINER STATS ERROR", err)
		return
	}
	cnt.statsDialog.SetContainerInfo(cntID, cntName)
	cnt.statsDialog.SetStatsChannel(&statsChan)
	cnt.statsDialog.SetStatsStream(&stream)
	cnt.statsDialog.Display()
}

func (cnt *Containers) cexec() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to perform exec command"))
		return
	}
	cntID, cntName := cnt.getSelectedItem()
	cnt.execDialog.SetContainerID(cntID, cntName)
	cnt.execDialog.Display()
}

func (cnt *Containers) exec() {
	cnt.execDialog.Hide()
	_, _, width, height := cnt.table.GetInnerRect()
	// TODO better calculation
	width = width - (2 * dialogs.DialogPadding) - 6
	height = height - (2 * (dialogs.DialogPadding - 1)) - 2*dialogs.DialogFormHeight - 4

	execOpts := cnt.execDialog.ContainerExecOptions()
	execOpts.TtyWidth = width
	execOpts.TtyHeight = height

	err := cnt.execTerminalDialog.PrepareForExec(cnt.selectedID, cnt.selectedName, &execOpts)
	if err != nil {
		title := fmt.Sprintf("CONTAINER (%s) EXEC ERROR", cnt.selectedID)
		cnt.displayError(title, err)
		return
	}

	prepareAndExec := func() {
		execSessionID, err := containers.NewExecSession(cnt.selectedID, execOpts)
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) EXEC ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}
		cnt.execTerminalDialog.SetExecInfo(cnt.selectedID, cnt.selectedName, execSessionID)
		containers.Exec(execSessionID, execOpts)
	}
	go prepareAndExec()

	cnt.execTerminalDialog.Display()

}

func (cnt *Containers) create() {
	createOpts := cnt.createDialog.ContainerCreateOptions()
	if createOpts.Name == "" || createOpts.Image == "" {
		cnt.displayError("CONTAINER CREATE ERROR", fmt.Errorf("container name or image name is empty"))
		return
	}
	cnt.progressDialog.SetTitle("container create in progress")
	cnt.progressDialog.Display()
	create := func() {
		warnings, err := containers.Create(createOpts)
		cnt.progressDialog.Hide()
		if err != nil {
			cnt.displayError("CONTAINER CREATE ERROR", err)
			return
		}
		if len(warnings) > 0 {
			headerLabel := fmt.Sprintf("%s (%s)", "", createOpts.Name)
			cnt.messageDialog.SetTitle("CONTAINER CREATE WARNINGS")
			cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, strings.Join(warnings, "\n"))
			cnt.messageDialog.Display()
		}
	}
	go create()
}

func (cnt *Containers) diff() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to display diff"))
		return
	}
	data, err := containers.Diff(cnt.selectedID)
	if err != nil {
		title := fmt.Sprintf("CONTAINER (%s) DIFF ERROR", cnt.selectedID)
		cnt.displayError(title, err)
		return
	}

	headerLabel := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)

	cnt.messageDialog.SetTitle("podman container diff")
	cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, strings.Join(data, "\n"))
	cnt.messageDialog.Display()
}

func (cnt *Containers) inspect() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to display inspect"))
		return
	}
	data, err := containers.Inspect(cnt.selectedID)
	if err != nil {
		title := fmt.Sprintf("CONTAINER (%s) INSPECT ERROR", cnt.selectedID)
		cnt.displayError(title, err)
		return
	}

	headerLabel := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)

	cnt.messageDialog.SetTitle("podman container inspect")
	cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, data)
	cnt.messageDialog.Display()
}

func (cnt *Containers) kill() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to kill"))
		return
	}
	cnt.progressDialog.SetTitle("container kill in progress")
	cnt.progressDialog.Display()
	kill := func(id string) {
		err := containers.Kill(id)
		cnt.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) KILL ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}
	}
	go kill(cnt.selectedID)
}

func (cnt *Containers) logs() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to display log"))
		return
	}
	logs, err := containers.Logs(cnt.selectedID)
	if err != nil {
		title := fmt.Sprintf("CONTAINER (%s) DISPLAY LOG ERROR", cnt.selectedID)
		cnt.displayError(title, err)
		return
	}

	headerLabel := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)

	cntLogs := strings.Join(logs, "\n")
	cntLogs = strings.ReplaceAll(cntLogs, "[", "")
	cntLogs = strings.ReplaceAll(cntLogs, "]", "")
	cnt.messageDialog.SetTitle("podman container logs")
	cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, cntLogs)
	cnt.messageDialog.TextScrollToEnd()
	cnt.messageDialog.Display()
}

func (cnt *Containers) pause() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to pause"))
		return
	}
	cnt.progressDialog.SetTitle("container pause in progress")
	cnt.progressDialog.Display()
	pause := func(id string) {
		err := containers.Pause(id)
		cnt.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) PAUSE ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}
	}
	go pause(cnt.selectedID)
}

func (cnt *Containers) port() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to display port"))
		return
	}
	data, err := containers.Port(cnt.selectedID)
	if err != nil {
		title := fmt.Sprintf("CONTAINER (%s) DISPLAY PORT ERROR", cnt.selectedID)
		cnt.displayError(title, err)
		return
	}

	headerLabel := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)

	cnt.messageDialog.SetTitle("podman container port")
	cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, strings.Join(data, "\n"))
	cnt.messageDialog.Display()
}

func (cnt *Containers) cprune() {
	cnt.confirmDialog.SetTitle("podman container prune")
	cnt.confirmData = "prune"
	cnt.confirmDialog.SetText("Are you sure you want to remove all unused containers ?")
	cnt.confirmDialog.Display()
}

func (cnt *Containers) prune() {
	cnt.progressDialog.SetTitle("container purne in progress")
	cnt.progressDialog.Display()
	prune := func() {
		errData, err := containers.Prune()
		cnt.progressDialog.Hide()
		if err != nil {
			cnt.displayError("CONTAINER PRUNE ERROR", err)
			return
		}
		if len(errData) > 0 {
			cnt.displayError("CONTAINER PRUNE ERROR", fmt.Errorf("%v", errData))
		}

	}
	go prune()
}

func (cnt *Containers) rename() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to rename"))
		return
	}
	cnt.cmdInputDialog.SetTitle("podman container rename")
	fgColor := style.GetColorHex(style.DialogFgColor)
	bgColor := fmt.Sprintf("#%x", style.DialogBorderColor.Hex())
	containerInfo := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)
	description := fmt.Sprintf("[%s:%s:b]CONTAINER ID:[:-:-] %s",
		fgColor, bgColor, containerInfo)
	cnt.cmdInputDialog.SetDescription(description)
	cnt.cmdInputDialog.SetSelectButtonLabel("rename")
	cnt.cmdInputDialog.SetLabel("target name ")

	cnt.cmdInputDialog.SetSelectedFunc(func() {
		newName := cnt.cmdInputDialog.GetInputText()
		cnt.cmdInputDialog.Hide()
		cnt.renameContainer(cnt.selectedID, newName)
	})
	cnt.cmdInputDialog.Display()
}

func (cnt *Containers) renameContainer(id string, newName string) {
	cnt.progressDialog.SetTitle("container rename in progress")
	cnt.progressDialog.Display()
	renameFunc := func() {
		err := containers.Rename(id, newName)
		cnt.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) RENAME ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}
		cnt.UpdateData()
	}
	go renameFunc()
}

func (cnt *Containers) rm() {
	cntID, cntName := cnt.getSelectedItem()
	if cntID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to remove"))
		return
	}
	cnt.confirmDialog.SetTitle("podman container remove")
	cnt.confirmData = "rm"
	bgColor := style.GetColorHex(style.DialogBorderColor)
	fgColor := style.GetColorHex(style.DialogFgColor)
	containerItem := fmt.Sprintf("[%s:%s:b]CONTAINER ID:[:-:-] %s(%s)", fgColor, bgColor, cntID, cntName)

	description := fmt.Sprintf("%s\n\nAre you sure you want to remove the selected container ?", containerItem)
	cnt.confirmDialog.SetText(description)
	cnt.confirmDialog.Display()
}

func (cnt *Containers) remove() {
	cnt.progressDialog.SetTitle("container remove in progress")
	cnt.progressDialog.Display()
	remove := func(id string) {
		errData, err := containers.Remove(id)
		cnt.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) REMOVE ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}
		if len(errData) > 0 {
			title := fmt.Sprintf("CONTAINER (%s) REMOVE ERROR", cnt.selectedID)
			cnt.displayError(title, fmt.Errorf("%v", errData))
		}
	}
	go remove(cnt.selectedID)
}

func (cnt *Containers) start() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to start"))
		return
	}
	cnt.progressDialog.SetTitle("container start in progress")
	cnt.progressDialog.Display()
	start := func(id string) {
		err := containers.Start(id)
		cnt.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) START ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}
	}
	go start(cnt.selectedID)
}

func (cnt *Containers) stop() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to stop"))
		return
	}
	cnt.progressDialog.SetTitle("container stop in progress")
	cnt.progressDialog.Display()
	stop := func(id string) {
		err := containers.Stop(id)
		cnt.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) STOP ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}
	}
	go stop(cnt.selectedID)
}

func (cnt *Containers) top() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to display top"))
		return
	}
	data, err := containers.Top(cnt.selectedID)
	if err != nil {
		title := fmt.Sprintf("CONTAINER (%s) TOP ERROR", cnt.selectedID)
		cnt.displayError(title, err)
		return
	}

	cntID, cntName := cnt.getSelectedItem()
	cnt.topDialog.UpdateResults(dialogs.TopContainerInfo, cntID, cntName, data)
	cnt.topDialog.Display()
}

func (cnt *Containers) unpause() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to unpause"))
		return
	}
	cnt.progressDialog.SetTitle("container unpause in progress")
	cnt.progressDialog.Display()
	unpause := func(id string) {
		err := containers.Unpause(id)
		cnt.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) UNPAUSE ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
		}
	}
	go unpause(cnt.selectedID)
}
