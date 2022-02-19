package containers

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/ui/dialogs"
	bcontainers "github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

func (cnt *Containers) runCommand(cmd string) {
	switch cmd {
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
	var message string
	if title != "" {
		message = fmt.Sprintf("%s: %v", title, err)
	} else {
		message = fmt.Sprintf("%v", err)
	}

	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	cnt.errorDialog.SetText(message)
	cnt.errorDialog.Display()
}

func (cnt *Containers) stats() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to perform exec command"))
		return
	}
	cntID, cntName := cnt.getSelectedItem()
	cntStatus, err := containers.Status(cntID)
	if err != nil {
		cnt.displayError("", fmt.Errorf("there is no container to perform exec command"))
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

	sessionID, err := cnt.execTerminalDialog.PrepareForExec(cnt.selectedID, cnt.selectedName, &execOpts)
	if err != nil {
		title := fmt.Sprintf("CONTAINER (%s) EXEC ERROR", cnt.selectedID)
		cnt.displayError(title, err)
		return
	}
	go containers.Exec(sessionID, execOpts)
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
			cnt.messageDialog.SetTitle("CONTAINER CREATE WARNINGS")
			cnt.messageDialog.SetText(strings.Join(warnings, "\n"))
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
	cnt.messageDialog.SetTitle("podman container diff")
	cnt.messageDialog.SetText(strings.Join(data, "\n"))
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
	cnt.messageDialog.SetTitle("podman container inspect")
	cnt.messageDialog.SetText(data)
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
	cntLogs := strings.Join(logs, "\n")
	cntLogs = strings.ReplaceAll(cntLogs, "[", "")
	cntLogs = strings.ReplaceAll(cntLogs, "]", "")
	cnt.messageDialog.SetTitle("podman container logs")
	cnt.messageDialog.SetText(cntLogs)
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
	cnt.messageDialog.SetTitle("podman container port")
	cnt.messageDialog.SetText(strings.Join(data, "\n"))
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
	description := fmt.Sprintf("[white::]container name : [black::]%s[white::]\ncontainer ID   : [black::]%s", cnt.selectedName, cnt.selectedID)
	cnt.cmdInputDialog.SetDescription(description)
	cnt.cmdInputDialog.SetSelectButtonLabel("rename")
	cnt.cmdInputDialog.SetLabel("target name")
	cnt.cmdInputDialog.SetSelectedFunc(func() {
		newName := cnt.cmdInputDialog.GetInputText()
		cnt.cmdInputDialog.Hide()
		err := containers.Rename(cnt.selectedID, newName)
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) RENAME ERROR", cnt.selectedID)
			cnt.displayError(title, err)
		}
	})
	cnt.cmdInputDialog.Display()
}

func (cnt *Containers) rm() {
	if cnt.selectedID == "" {
		cnt.displayError("", fmt.Errorf("there is no container to remove"))
		return
	}
	cnt.confirmDialog.SetTitle("podman container remove")
	cnt.confirmData = "rm"
	description := fmt.Sprintf("Are you sure you want to remove following container ? \n\nCONTAINER ID : %s", cnt.selectedID)
	cnt.confirmDialog.SetText(description)
	cnt.confirmDialog.Display()
}

func (cnt *Containers) remove() {
	cnt.progressDialog.SetTitle("container remove in progress")
	cnt.progressDialog.Display()
	remove := func(id string) {
		err := containers.Remove(id)
		cnt.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("CONTAINER (%s) REMOVE ERROR", cnt.selectedID)
			cnt.displayError(title, err)
			return
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
	cnt.topDialog.UpdateResults(data)
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
