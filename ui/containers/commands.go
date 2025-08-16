package containers

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	bcontainers "github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

func (cnt *Containers) runCommand(cmd string) { //nolint:cyclop
	switch cmd {
	case "attach":
		cnt.attach()
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
	case "healthcheck":
		cnt.preHealthcheck()
	case "inspect":
		cnt.inspect()
	case "kill":
		cnt.kill()
	case "logs":
		cnt.logs()
	case "pause":
		cnt.pause()
	case utils.PruneCommandLabel:
		cnt.cprune()
	case "rename":
		cnt.rename()
	case "restore":
		cnt.preRestore()
	case "port":
		cnt.port()
	case "rm":
		cnt.rm()
	case "run":
		cnt.runDialog.Display()
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

func (cnt *Containers) attach() {
	cntID, cntName := cnt.getSelectedItem()
	if cntID == "" {
		cnt.displayError("", errNoContainerAttach)

		return
	}

	cnt.progressDialog.SetTitle("container attach in progress")
	cnt.progressDialog.Display()

	attachReady := make(chan bool)
	stdin, stdout := cnt.terminalDialog.InitAttachChannels()
	detachKeys := cnt.terminalDialog.DetachKeys()

	attach := func() {
		err := containers.Attach(cntID, stdin, stdout, attachReady, detachKeys)
		if err != nil {
			attachReady <- false

			title := fmt.Sprintf("CONTAINER (%s) ATTACH ERROR", cntID)

			cnt.progressDialog.Hide()
			cnt.displayError(title, err)
			cnt.appFocusHandler()

			return
		}
	}

	waitForAttach := func() {
		isReady := <-attachReady
		if isReady {
			cnt.progressDialog.Hide()
			cnt.terminalDialog.SetContainerInfo(cntID, cntName)
			cnt.terminalDialog.Display()
			cnt.appFocusHandler()
		}
	}

	go waitForAttach()
	go attach()
}

func (cnt *Containers) preHealthcheck() {
	cntID, cntName := cnt.getSelectedItem()
	if cntID == "" {
		cnt.displayError("", errNoContainerHealthCheck)

		return
	}

	cnt.progressDialog.SetTitle("container healthcheck in progress")
	cnt.progressDialog.Display()

	cntHealthCheck := func() {
		report, err := containers.HealthCheck(cntID)

		cnt.progressDialog.Hide()

		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) HEALTHCHECK ERROR", cntID)

			cnt.displayError(title, err)
			cnt.appFocusHandler()

			return
		}

		headerLabel := fmt.Sprintf("%s (%s)", cntID, cntName)

		cnt.messageDialog.SetTitle("podman container healthcheck")
		cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, report)
		cnt.messageDialog.Display()
		cnt.appFocusHandler()
	}

	go cntHealthCheck()
}

func (cnt *Containers) preRestore() {
	var ( //nolint:prealloc
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
			cnt.appFocusHandler()

			return
		}

		headerLabel := fmt.Sprintf("%s (%s)", restoreOptions.ContainerID, restoreOptions.Name)

		cnt.progressDialog.Hide()
		cnt.messageDialog.SetTitle("podman container restore")
		cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, report)
		cnt.messageDialog.Display()
		cnt.appFocusHandler()
	}

	go restore()
}

func (cnt *Containers) preCheckpoint() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerCheckpoint)

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
			cnt.appFocusHandler()

			return
		}

		headerLabel := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)

		cnt.progressDialog.Hide()
		cnt.messageDialog.SetTitle("podman container checkpoint")
		cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, report)
		cnt.messageDialog.Display()
		cnt.appFocusHandler()
	}

	go checkpoint()
}

func (cnt *Containers) preCommit() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerCommit)

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
			cnt.appFocusHandler()

			return
		}

		headerLabel := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)

		cnt.progressDialog.Hide()
		cnt.messageDialog.SetTitle("podman container commit")
		cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, response)
		cnt.messageDialog.Display()
		cnt.appFocusHandler()
	}

	go cntCommit()
}

func (cnt *Containers) stats() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerStat)

		return
	}

	cntID, cntName := cnt.getSelectedItem()

	cntStatus, err := containers.Status(cntID)
	if err != nil {
		cnt.displayError("", err)

		return
	}

	if cntStatus != "running" {
		cnt.displayError("", fmt.Errorf("container (%s) status improper", cntID)) //nolint:err113

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
	cntID, cntName := cnt.getSelectedItem()

	if cntID == "" {
		cnt.displayError("", errNoContainerExec)

		return
	}

	cnt.execDialog.SetContainerID(cntID, cntName)
	cnt.execDialog.Display()
}

func (cnt *Containers) exec() {
	cnt.execDialog.Hide()

	cntID, cntName := cnt.getSelectedItem()
	_, _, width, height := cnt.table.GetInnerRect()

	width = width - (2 * dialogs.DialogPadding) - 6                                      //nolint:mnd
	height = height - (2 * (dialogs.DialogPadding - 1)) - 2*dialogs.DialogFormHeight - 4 //nolint:mnd

	execOpts := cnt.execDialog.ContainerExecOptions()
	execOpts.TtyWidth = width
	execOpts.TtyHeight = height

	execOpts.InputStream, execOpts.OutputStream = cnt.terminalDialog.InitExecChannels()
	execOpts.DetachKeys = cnt.terminalDialog.DetachKeys()

	cnt.terminalDialog.SetContainerInfo(cntID, cntName)

	execSessionID, err := containers.NewExecSession(cnt.selectedID, execOpts)
	if err != nil {
		title := fmt.Sprintf("CONTAINER (%s) EXEC ERROR", cnt.selectedID)

		cnt.displayError(title, err)

		return
	}

	prepareAndExec := func() {
		cnt.terminalDialog.SetSessionID(execSessionID)
		containers.Exec(execSessionID, execOpts)
		cnt.terminalDialog.SetAlreadyDetach(true)
	}

	go prepareAndExec()

	cnt.terminalDialog.Display()
}

func (cnt *Containers) run() {
	runOpts := cnt.runDialog.ContainerCreateOptions()
	if runOpts.Image == "" {
		cnt.displayError("CONTAINER RUN ERROR", errEmptyContainerImageName)

		return
	}

	cnt.progressDialog.SetTitle("container run in progress")
	cnt.progressDialog.Display()

	if runOpts.Detach {
		cnt.runDetach(runOpts)

		return
	}

	cnt.runAttach(runOpts)
}

func (cnt *Containers) runDetach(runOpts containers.CreateOptions) {
	go func() {
		warnings, cntID, err := containers.Create(runOpts, true)
		if err != nil {
			cnt.progressDialog.Hide()
			cnt.displayError("CONTAINER RUN ERROR", err)
			cnt.appFocusHandler()

			return
		}

		if len(warnings) > 0 {
			cnt.progressDialog.Hide()

			headerLabel := fmt.Sprintf("%s (%s)", "", runOpts.Name)

			cnt.messageDialog.SetTitle("CONTAINER RUN WARNINGS")
			cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, strings.Join(warnings, "\n"))
			cnt.messageDialog.Display()
			cnt.appFocusHandler()

			return
		}

		err = containers.Start(cntID)
		if err != nil {
			cnt.progressDialog.Hide()
			cnt.displayError("CONTAINER RUN ERROR", err)
			cnt.appFocusHandler()

			return
		}

		cnt.progressDialog.Hide()
	}()
}

func (cnt *Containers) runAttach(runOpts containers.CreateOptions) {
	runStatusChan := make(chan bool)
	attachReady := make(chan bool)
	runIDChan := make(chan string)
	stdin, stdout := cnt.terminalDialog.InitAttachChannels()
	detachKeys := cnt.terminalDialog.DetachKeys()

	run := func() {
		warnings, cntID, err := containers.Create(runOpts, true)
		if err != nil {
			runStatusChan <- false

			cnt.displayError("CONTAINER RUN ERROR", err)
			cnt.appFocusHandler()

			return
		}

		if len(warnings) > 0 {
			runStatusChan <- false

			headerLabel := fmt.Sprintf("%s (%s)", "", runOpts.Name)

			cnt.messageDialog.SetTitle("CONTAINER RUN WARNINGS")
			cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, strings.Join(warnings, "\n"))
			cnt.messageDialog.Display()
			cnt.appFocusHandler()

			return
		}

		runIDChan <- cntID

		err = containers.RunInitAttach(cntID, stdin, stdout, attachReady, detachKeys)
		if err != nil {
			attachReady <- false

			cnt.displayError("CONTAINER RUN ERROR", err)
			cnt.appFocusHandler()

			return
		}

		runStatusChan <- true
	}

	waitForAttach := func() {
		cntID := ""

		for {
			select {
			case <-runStatusChan:
				cnt.terminalDialog.SetAlreadyDetach(true)
				cnt.progressDialog.Hide()

				return

			case id := <-runIDChan:
				cntID = id
			case isReady := <-attachReady:
				cnt.progressDialog.Hide()

				if isReady {
					err := containers.Start(cntID)
					if err != nil {
						cnt.displayError("CONTAINER RUN ERROR", err)
						cnt.appFocusHandler()

						return
					}

					cnt.terminalDialog.SetContainerInfo(cntID, "")
					cnt.terminalDialog.Display()
					cnt.appFocusHandler()
				}
			}
		}
	}

	go waitForAttach()
	go run()
}

func (cnt *Containers) create() {
	createOpts := cnt.createDialog.ContainerCreateOptions()
	if createOpts.Image == "" {
		cnt.displayError("CONTAINER CREATE ERROR", errEmptyContainerImageName)

		return
	}

	cnt.progressDialog.SetTitle("container create in progress")
	cnt.progressDialog.Display()

	create := func() {
		warnings, _, err := containers.Create(createOpts, false)

		cnt.progressDialog.Hide()

		if err != nil {
			cnt.displayError("CONTAINER CREATE ERROR", err)
			cnt.appFocusHandler()

			return
		}

		if len(warnings) > 0 {
			headerLabel := fmt.Sprintf("%s (%s)", "", createOpts.Name)

			cnt.messageDialog.SetTitle("CONTAINER CREATE WARNINGS")
			cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, strings.Join(warnings, "\n"))
			cnt.messageDialog.Display()
			cnt.appFocusHandler()
		}
	}

	go create()
}

func (cnt *Containers) diff() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerDiff)

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
	cnt.messageDialog.DisplayFullSize()
}

func (cnt *Containers) inspect() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerInspect)

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
	cnt.messageDialog.DisplayFullSize()
}

func (cnt *Containers) kill() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerKill)

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
			cnt.appFocusHandler()

			return
		}
	}

	go kill(cnt.selectedID)
}

func (cnt *Containers) logs() {
	cntID, cntName := cnt.getSelectedItem()
	if cntID == "" {
		cnt.displayError("", errNoContainerLogs)

		return
	}

	cnt.progressDialog.SetTitle("container logs in progress")
	cnt.progressDialog.Display()

	getLogs := func() {
		logData, err := containers.Logs(cntID)

		cnt.progressDialog.Hide()

		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) DISPLAY LOG ERROR", cntID)

			cnt.displayError(title, err)
			cnt.appFocusHandler()

			return
		}

		headerLabel := fmt.Sprintf("%s (%s)", cntID, cntName)

		cntLogs := strings.Join(logData, "\n")
		cntLogs = strings.ReplaceAll(cntLogs, "[", "")
		cntLogs = strings.ReplaceAll(cntLogs, "]", "")

		cnt.messageDialog.SetTitle("podman container logs")
		cnt.messageDialog.SetText(dialogs.MessageContainerInfo, headerLabel, cntLogs)
		cnt.messageDialog.TextScrollToEnd()
		cnt.messageDialog.DisplayFullSize()
		cnt.appFocusHandler()
	}

	go getLogs()
}

func (cnt *Containers) pause() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerPause)

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
			cnt.appFocusHandler()

			return
		}
	}

	go pause(cnt.selectedID)
}

func (cnt *Containers) port() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerPorts)

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
	cnt.progressDialog.SetTitle("container prune in progress")
	cnt.progressDialog.Display()

	prune := func() {
		errData, err := containers.Prune()

		cnt.progressDialog.Hide()

		if err != nil {
			cnt.displayError("CONTAINER PRUNE ERROR", err)
			cnt.appFocusHandler()

			return
		}

		if len(errData) > 0 {
			cnt.displayError("CONTAINER PRUNE ERROR", fmt.Errorf("%v", errData)) //nolint:err113
		}
	}

	go prune()
}

func (cnt *Containers) rename() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerRename)

		return
	}

	cnt.cmdInputDialog.SetTitle("podman container rename")

	fgColor := style.GetColorHex(style.DialogFgColor)
	bgColor := fmt.Sprintf("#%x", style.DialogBorderColor.Hex())
	containerInfo := fmt.Sprintf("%s (%s)", cnt.selectedID, cnt.selectedName)
	description := fmt.Sprintf("[%s:%s:b]%s[:-:-] %s",
		fgColor, bgColor, utils.ContainerIDLabel, containerInfo)

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
			cnt.appFocusHandler()

			return
		}

		cnt.UpdateData()
	}

	go renameFunc()
}

func (cnt *Containers) rm() {
	cntID, cntName := cnt.getSelectedItem()
	if cntID == "" {
		cnt.displayError("", errNoContainerRemove)

		return
	}

	cnt.confirmDialog.SetTitle("podman container remove")
	cnt.confirmData = "rm"
	bgColor := style.GetColorHex(style.DialogBorderColor)
	fgColor := style.GetColorHex(style.DialogFgColor)
	containerItem := fmt.Sprintf("[%s:%s:b]%s[:-:-] %s(%s)", fgColor, bgColor, utils.ContainerIDLabel, cntID, cntName)
	description := fmt.Sprintf("%s\n\nAre you sure you want to remove the selected container ?", //nolint:perfsprint
		containerItem)

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
			cnt.appFocusHandler()

			return
		}

		if len(errData) > 0 {
			title := fmt.Sprintf("CONTAINER (%s) REMOVE ERROR", cnt.selectedID)

			cnt.displayError(title, fmt.Errorf("%v", errData)) //nolint:err113
			cnt.appFocusHandler()
		}
	}

	go remove(cnt.selectedID)
}

func (cnt *Containers) start() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerStart)

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
			cnt.appFocusHandler()

			return
		}
	}
	go start(cnt.selectedID)
}

func (cnt *Containers) stop() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerStop)

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
			cnt.appFocusHandler()

			return
		}
	}

	go stop(cnt.selectedID)
}

func (cnt *Containers) top() {
	if cnt.selectedID == "" {
		cnt.displayError("", errNoContainerTop)

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
		cnt.displayError("", errNoContainerUnpause)

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
			cnt.appFocusHandler()

			return
		}
	}

	go unpause(cnt.selectedID)
}
