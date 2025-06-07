package pods

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (pods *Pods) Draw(screen tcell.Screen) {
	pods.Box.DrawForSubclass(screen, pods)
	pods.Box.SetBorder(false)

	podViewX, podViewY, podViewW, podViewH := pods.GetInnerRect()

	pods.refresh(podViewW)
	pods.table.SetRect(podViewX, podViewY, podViewW, podViewH)
	pods.table.SetBorder(true)

	pods.table.Draw(screen)

	x, y, width, height := pods.table.GetInnerRect()

	// error dialog
	if pods.errorDialog.IsDisplay() {
		pods.errorDialog.SetRect(x, y, width, height)
		pods.errorDialog.Draw(screen)

		return
	}

	// command dialog
	if pods.cmdDialog.IsDisplay() {
		pods.cmdDialog.SetRect(x, y, width, height)
		pods.cmdDialog.Draw(screen)

		return
	}

	// create dialog
	if pods.createDialog.IsDisplay() {
		pods.createDialog.SetRect(x, y, width, height)
		pods.createDialog.Draw(screen)

		return
	}

	// confirm dialog
	if pods.confirmDialog.IsDisplay() {
		pods.confirmDialog.SetRect(x, y, width, height)
		pods.confirmDialog.Draw(screen)

		return
	}

	// message dialog
	if pods.messageDialog.IsDisplay() {
		if pods.messageDialog.IsDisplayFullSize() {
			pods.messageDialog.SetRect(podViewX, podViewY, podViewW, podViewH)
		} else {
			pods.messageDialog.SetRect(x, y, width, height+1)
		}

		pods.messageDialog.Draw(screen)

		return
	}

	// progress dialog
	if pods.progressDialog.IsDisplay() {
		pods.progressDialog.SetRect(x, y, width, height)
		pods.progressDialog.Draw(screen)
	}

	// top dialog
	if pods.topDialog.IsDisplay() {
		pods.topDialog.SetRect(x, y, width, height)
		pods.topDialog.Draw(screen)

		return
	}

	// stats dialogs
	if pods.statsDialog.IsDisplay() {
		pods.statsDialog.SetRect(podViewX, podViewY, podViewW, podViewH)
		pods.statsDialog.Draw(screen)

		return
	}
}
