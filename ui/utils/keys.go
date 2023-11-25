package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
)

// application key bindings names.

var (
	CommandMenuKey = uiKeyInfo{
		Key:      tcell.Key(256), //nolint:gomnd
		KeyRune:  rune('m'),
		KeyLabel: "m",
		KeyDesc:  "display command menu",
	}
	NextScreenKey = uiKeyInfo{
		Key:      tcell.Key(256), //nolint:gomnd
		KeyRune:  rune('l'),
		KeyLabel: "l",
		KeyDesc:  "switch to next screen",
	}
	PreviousScreenKey = uiKeyInfo{
		Key:      tcell.Key(256), //nolint:gomnd
		KeyRune:  rune('h'),
		KeyLabel: "h",
		KeyDesc:  "switch to previous screen",
	}
	MoveUpKey = uiKeyInfo{
		Key:      tcell.KeyUp,
		KeyRune:  rune('k'),
		KeyLabel: "k",
		KeyDesc:  "move up",
	}
	MoveDownKey = uiKeyInfo{
		Key:      tcell.KeyDown,
		KeyRune:  rune('j'),
		KeyLabel: "j",
		KeyDesc:  "move down",
	}
	CloseDialogKey = uiKeyInfo{
		Key:      tcell.KeyEsc,
		KeyLabel: "Esc",
		KeyDesc:  "close the active dialog",
	}
	SwitchFocusKey = uiKeyInfo{
		Key:      tcell.KeyTab,
		KeyLabel: "Tab",
		KeyDesc:  "switch between widgets",
	}
	DeleteKey = uiKeyInfo{
		Key:      tcell.KeyDelete,
		KeyLabel: "Delete",
		KeyDesc:  "delete the selected item",
	}
	ArrowUpKey = uiKeyInfo{
		Key:      tcell.KeyUp,
		KeyLabel: "Arrow Up",
		KeyDesc:  "move up",
	}
	ArrowDownKey = uiKeyInfo{
		Key:      tcell.KeyDown,
		KeyLabel: "Arrow Down",
		KeyDesc:  "move down",
	}
	ArrowLeftKey = uiKeyInfo{
		Key:      tcell.KeyLeft,
		KeyLabel: "Arrow Left",
		KeyDesc:  "previous screen",
	}
	ArrowRightKey = uiKeyInfo{
		Key:      tcell.KeyRight,
		KeyLabel: "Arrow Right",
		KeyDesc:  "next screen",
	}
	ScrollUpKey = uiKeyInfo{
		Key:      tcell.KeyPgUp,
		KeyLabel: "Page Up",
		KeyDesc:  "scroll up",
	}
	ScrollDownKey = uiKeyInfo{
		Key:      tcell.KeyPgDn,
		KeyLabel: "Page Down",
		KeyDesc:  "scroll down",
	}
	AppExitKey = uiKeyInfo{
		Key:      tcell.KeyCtrlC,
		KeyLabel: "Ctrl+c",
		KeyDesc:  "exit application",
	}
	HelpScreenKey = uiKeyInfo{
		Key:      tcell.KeyF1,
		KeyLabel: "F1",
		KeyDesc:  "display help screen",
	}
	SystemScreenKey = uiKeyInfo{
		Key:      tcell.KeyF2,
		KeyLabel: "F2",
		KeyDesc:  "display system screen",
	}
	PodsScreenKey = uiKeyInfo{
		Key:      tcell.KeyF3,
		KeyLabel: "F3",
		KeyDesc:  "display pods screen",
	}
	ContainersScreenKey = uiKeyInfo{
		Key:      tcell.KeyF4,
		KeyLabel: "F4",
		KeyDesc:  "display containers screen",
	}
	VolumesScreenKey = uiKeyInfo{
		Key:      tcell.KeyF5,
		KeyLabel: "F5",
		KeyDesc:  "display volumes screen",
	}
	ImagesScreenKey = uiKeyInfo{
		Key:      tcell.KeyF6,
		KeyLabel: "F6",
		KeyDesc:  "display images screen",
	}
	NetworksScreenKey = uiKeyInfo{
		Key:      tcell.KeyF7,
		KeyLabel: "F7",
		KeyDesc:  "display networks screen",
	}
)

// UIKeysBindings user interface key bindings.
var UIKeysBindings = []uiKeyInfo{
	CommandMenuKey,
	NextScreenKey,
	PreviousScreenKey,
	MoveUpKey,
	MoveDownKey,
	CloseDialogKey,
	SwitchFocusKey,
	DeleteKey,
	ArrowUpKey,
	ArrowDownKey,
	ArrowLeftKey,
	ArrowRightKey,
	ScrollUpKey,
	ScrollDownKey,
	AppExitKey,
	HelpScreenKey,
	SystemScreenKey,
	PodsScreenKey,
	ContainersScreenKey,
	VolumesScreenKey,
	ImagesScreenKey,
	NetworksScreenKey,
}

type uiKeyInfo struct {
	Key      tcell.Key
	KeyRune  rune
	KeyLabel string
	KeyDesc  string
}

func (key *uiKeyInfo) Label() string {
	return key.KeyLabel
}

func (key *uiKeyInfo) Rune() rune {
	return key.KeyRune
}

func (key *uiKeyInfo) EventKey() tcell.Key {
	return key.Key
}

func (key *uiKeyInfo) Description() string {
	return key.KeyDesc
}

// ParseKeyEventKey parsed and changes key events key and rune base on keyname.
func ParseKeyEventKey(event *tcell.EventKey) *tcell.EventKey {
	log.Debug().Msgf("utils: parse key event (%v) key=%v name=%v", event, event.Key(), event.Name())

	switch event.Rune() {
	case MoveUpKey.KeyRune:
		return tcell.NewEventKey(MoveUpKey.Key, MoveUpKey.KeyRune, tcell.ModNone)
	case MoveDownKey.KeyRune:
		return tcell.NewEventKey(MoveDownKey.Key, MoveDownKey.KeyRune, tcell.ModNone)
	}

	switch event.Key() { //nolint:exhaustive
	case ArrowLeftKey.Key:
		return tcell.NewEventKey(PreviousScreenKey.Key, PreviousScreenKey.KeyRune, tcell.ModNone)
	case ArrowRightKey.Key:
		return tcell.NewEventKey(NextScreenKey.Key, NextScreenKey.KeyRune, tcell.ModNone)
	}

	return event
}
