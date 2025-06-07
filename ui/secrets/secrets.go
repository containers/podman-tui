package secrets

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/secrets/secdialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
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
	errNoSecretRemove        = errors.New("there is no secret to remove")
	errNoSecretInspect       = errors.New("there is no secret to display inspect")
	errSecretFileAndText     = errors.New("cannot select secret file and secret text together")
	errEmptySecretFileOrText = errors.New("secret content not provided")
)

// Secrets implements the secrets page primitive.
type Secrets struct {
	*tview.Box
	title           string
	headers         []string
	table           *tview.Table
	cmdDialog       *dialogs.CommandDialog
	messageDialog   *dialogs.MessageDialog
	errorDialog     *dialogs.ErrorDialog
	progressDialog  *dialogs.ProgressDialog
	confirmDialog   *dialogs.ConfirmDialog
	createDialog    *secdialogs.SecretCreateDialog
	secretList      secretListReport
	appFocusHandler func()
}

type secretListReport struct {
	mu     sync.Mutex
	report []*types.SecretInfoReport
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
		createDialog:   secdialogs.NewSecretCreateDialog(),
	}

	secrets.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"create", "create a new secret"},
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

	// set create dialog function
	secrets.createDialog.SetCancelFunc(func() {
		secrets.createDialog.Hide()
	})

	secrets.createDialog.SetCreateFunc(func() {
		secrets.createDialog.Hide()
		secrets.create()
	})

	return secrets
}

// SetAppFocusHandler sets application focus handler.
func (s *Secrets) SetAppFocusHandler(handler func()) {
	s.appFocusHandler = handler
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

	if s.confirmDialog.HasFocus() || s.createDialog.HasFocus() {
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

	if s.confirmDialog.HasFocus() || s.createDialog.HasFocus() {
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

	// create dialog
	if s.createDialog.IsDisplay() {
		delegate(s.createDialog)

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

	if s.createDialog.IsDisplay() {
		s.createDialog.Hide()
	}
}

func (s *Secrets) getSelectedItem() (int, string, string) {
	var (
		rowIndex int
		secID    string
		secName  string
	)

	if s.table.GetRowCount() <= 1 {
		return rowIndex, secID, secName
	}

	rowIndex, _ = s.table.GetSelection()
	secID = s.table.GetCell(rowIndex, 0).Text
	secName = s.table.GetCell(rowIndex, 1).Text

	return rowIndex, secID, secName
}
