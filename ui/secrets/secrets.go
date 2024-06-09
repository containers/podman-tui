package secrets

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
)

const (
	viewSecretsIDColIndex = 0 + iota
	viewSecretsNameColIndex
	viewSecretsDriverColIndex
	viewSecretsCreatedColIndex
	viewSecretsUpdatedColIndex
)

var (
	errNoSecretRemove  = errors.New("there is no secret to remove")
	errNoSecretInspect = errors.New("there is no secret to display inspect")
)

// Secrets implements the secrets page primitive.
type Secrets struct {
	*tview.Box
	title          string
	headers        []string
	table          *tview.Table
	cmdDialog      *dialogs.CommandDialog
	messageDialog  *dialogs.MessageDialog
	errorDialog    *dialogs.ErrorDialog
	progressDialog *dialogs.ProgressDialog
	confirmDialog  *dialogs.ConfirmDialog
}

// NewSecrets returns secrets page view.
func NewSecrets() *Secrets {
	secrets := &Secrets{
		Box:            tview.NewBox(),
		title:          "secrets",
		headers:        []string{"id", "name", "driver", "created", "updated"},
		table:          tview.NewTable(),
		messageDialog:  dialogs.NewMessageDialog(""),
		errorDialog:    dialogs.NewErrorDialog(),
		progressDialog: dialogs.NewProgressDialog(),
		confirmDialog:  dialogs.NewConfirmDialog(),
	}

	secrets.cmdDialog = dialogs.NewCommandDialog([][]string{
		// {"create", "create a new secret"},
		{"inspect", "inspect a secret"},
		{"rm", "remove a secret"},
	})

	secrets.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(secrets.title)))
	secrets.table.SetBorderColor(style.BorderColor)
	secrets.table.SetBackgroundColor(style.BgColor)
	secrets.table.SetTitleColor(style.FgColor)
	secrets.table.SetBorder(true)

	secrets.table.SetFixed(1, 1)
	secrets.table.SetSelectable(true, false)

	// set command dialog functions
	secrets.cmdDialog.SetSelectedFunc(func() {
		secrets.cmdDialog.Hide()
		secrets.runCommand(secrets.cmdDialog.GetSelectedItem())
	})

	secrets.cmdDialog.SetCancelFunc(func() {
		secrets.cmdDialog.Hide()
	})

	// set message dialog function
	secrets.messageDialog.SetCancelFunc(func() {
		secrets.messageDialog.Hide()
	})

	// set confirm dialog functions
	secrets.confirmDialog.SetSelectedFunc(func() {
		secrets.confirmDialog.Hide()

		secrets.remove()
	})

	secrets.confirmDialog.SetCancelFunc(func() {
		secrets.confirmDialog.Hide()
	})

	return secrets
}

// GetTitle returns primitive title.
func (s *Secrets) GetTitle() string {
	return s.title
}

// HasFocus returns whether or not this primitive has focus.
func (s *Secrets) HasFocus() bool {
	if s.table.HasFocus() || s.Box.HasFocus() {
		return true
	}

	if s.cmdDialog.HasFocus() || s.errorDialog.HasFocus() {
		return true
	}

	if s.messageDialog.HasFocus() || s.progressDialog.HasFocus() {
		return true
	}

	if s.confirmDialog.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (s *Secrets) SubDialogHasFocus() bool {
	if s.cmdDialog.HasFocus() || s.errorDialog.HasFocus() {
		return true
	}

	if s.messageDialog.HasFocus() || s.progressDialog.HasFocus() {
		return true
	}

	if s.confirmDialog.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus.
func (s *Secrets) Focus(delegate func(p tview.Primitive)) {
	// error dialog
	if s.errorDialog.IsDisplay() {
		delegate(s.errorDialog)

		return
	}

	// message dialog
	if s.messageDialog.IsDisplay() {
		delegate(s.messageDialog)

		return
	}

	// confirmation dialog
	if s.confirmDialog.IsDisplay() {
		delegate(s.confirmDialog)

		return
	}

	// cmd dialog
	if s.cmdDialog.IsDisplay() {
		delegate(s.cmdDialog)

		return
	}

	delegate(s.table)
}

// HideAllDialogs hides all sub dialogs.
func (s *Secrets) HideAllDialogs() {
	if s.errorDialog.IsDisplay() {
		s.errorDialog.Hide()
	}

	if s.messageDialog.IsDisplay() {
		s.messageDialog.Hide()
	}

	if s.progressDialog.IsDisplay() {
		s.progressDialog.Hide()
	}

	if s.confirmDialog.IsDisplay() {
		s.confirmDialog.Hide()
	}

	if s.cmdDialog.IsDisplay() {
		s.cmdDialog.Hide()
	}
}

func (s *Secrets) getSelectedItem() (string, string) {
	var (
		secID   string
		secName string
	)

	if s.table.GetRowCount() <= 1 {
		return secID, secName
	}

	row, _ := s.table.GetSelection()
	secID = s.table.GetCell(row, 0).Text
	secName = s.table.GetCell(row, 1).Text

	return secID, secName
}
