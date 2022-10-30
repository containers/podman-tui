//go:build windows
// +build windows

package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Styles represent default application style
var Styles = theme{
	PageTable: pageTable{
		FgColor: tcell.ColorPink,
		BgColor: tcell.ColorPink,
		HeaderRow: headerRow{
			FgColor: tcell.ColorBlack,
			BgColor: tcell.ColorPink,
		},
	},
	InfoBar: infoBar{
		ItemFgColor:  tcell.ColorPink,
		ValueFgColor: tcell.ColorWhite,
		ProgressBar: progressBar{
			FgColor:       tcell.ColorWhite,
			BarEmptyColor: tcell.ColorWhite,
			BarOKColor:    tcell.ColorLime,
			BarWarnColor:  tcell.ColorYellow,
			BarCritColor:  tcell.ColorRed,
		},
	},
	Menu: menu{
		FgColor: tcell.ColorWhite,
		BgColor: tcell.ColorBlack,
		Item: menuItem{
			FgColor: tcell.ColorBlack,
			BgColor: tcell.ColorPink,
		},
	},
	Help: help{
		BorderColor:   tcell.ColorPink,
		BgColor:       tview.Styles.PrimitiveBackgroundColor,
		FgColor:       tcell.ColorWhite,
		HeaderFgColor: tcell.ColorPink,
	},
	ConnectionProgressDialog: connectionProgressDialog{
		BgColor:     tcell.ColorPink,
		FgColor:     tcell.ColorPurple,
		PrgBarColor: tcell.ColorFuchsia,
		BorderColor: tcell.ColorWhite,
		TitleColor:  tcell.ColorWhite,
	},
	CommandDialog: commandDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorPurple,
		},
	},
	TopDialog: topDialog{
		BgColor:                tcell.ColorPink,
		FgColor:                tcell.ColorWhite,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tview.Styles.PrimaryTextColor,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorPurple,
		},
	},
	ConfirmDialog: confirmDialog{
		BgColor:     tcell.ColorPink,
		FgColor:     tcell.ColorWhite,
		ButtonColor: tcell.ColorPurple,
	},
	ImageSearchDialog: imageSearchDialog{
		BgColor:                tcell.ColorPink,
		FgColor:                tcell.ColorWhite,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tview.Styles.PrimaryTextColor,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorPurple,
		},
	},
	ImageHistoryDialog: imageHistoryDialog{
		BgColor:                tcell.ColorPink,
		FgColor:                tcell.ColorWhite,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tview.Styles.PrimaryTextColor,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorPurple,
		},
	},
	ImageImportDialog: imageImportDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ImageBuildDialog: imageBuildDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ImageBuildPrgDialog: imageBuildPrgDialog{
		BgColor:      tcell.ColorPink,
		FgColor:      tcell.ColorWhite,
		PrgCellColor: tcell.ColorOrange,
		Terminal: terminal{
			BorderColor: tview.Styles.PrimaryTextColor,
			BgColor:     tview.Styles.PrimitiveBackgroundColor,
			FgColor:     tview.Styles.PrimaryTextColor,
		},
	},
	ImageSaveDialog: imageSaveDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ImagePushDialog: imagePushDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	VolumeCreateDialog: volumeCreateDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	NetworkCreateDialog: networkCreateDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	NetworkConnectDialog: networkConnectDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ContainerCreateDialog: containerCreateDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ContainerExecDialog: containerExecDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ContainerExecTerminalDialog: containerExecTerminalDialog{
		BgColor:       tcell.ColorPink,
		FgColor:       tcell.ColorWhite,
		HeaderBgColor: tcell.ColorPurple,
		Terminal: terminal{
			BorderColor: tview.Styles.PrimaryTextColor,
			BgColor:     tcell.NewRGBColor(0, 0, 0),
			FgColor:     tcell.NewRGBColor(255, 255, 255),
		},
	},
	ContainerStatsDialog: containerStatsDialog{
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorWhite,
		TableHeaderFgColor:     tcell.ColorPink,
		BgColor:                tcell.ColorPink,
		FgColor:                tcell.ColorWhite,
	},
	ContainerCommitDialog: containerCommitDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ContainerCheckpointDialog: containerCheckpointDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ContainerRestoreDialog: containerRestoreDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	PodCreateDialog: podCreateDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	PodStatsDialog: podStatsDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorPink,
		},
	},
	DropdownStyle: dropdownStyle{
		Unselected: tcell.StyleDefault.Background(tcell.ColorFuchsia).Foreground(tcell.ColorBlack),
		Selected:   tcell.StyleDefault.Background(tcell.ColorPurple).Foreground(tcell.ColorWhite),
	},
	EventsDialog: eventsDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
		EventViewer: terminal{
			BorderColor: tview.Styles.PrimaryTextColor,
			BgColor:     tview.Styles.PrimitiveBackgroundColor,
			FgColor:     tview.Styles.PrimaryTextColor,
		},
	},
	ConnectionAddDialog: connectionAddDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	DiskUageDialog: diskUageDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorPurple,
		},
	},
	MessageDialog: messageDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
		Terminal: terminal{
			BorderColor: tview.Styles.PrimaryTextColor,
			BgColor:     tview.Styles.PrimitiveBackgroundColor,
			FgColor:     tview.Styles.PrimaryTextColor,
		},
	},
	ButtonPrimitive: buttonPrimitive{
		BgColor: tcell.ColorPurple,
	},
	InputFieldPrimitive: inputFieldPrimitive{
		BgColor: tcell.ColorPurple,
	},
	ProgressDailog: progressDailog{
		PgBarColor: tcell.ColorPurple,
	},
	InputDialog: inputDialog{
		BgColor: tcell.ColorPink,
		FgColor: tcell.ColorWhite,
	},
	ErrorDialog: errorDialog{
		HeaderFgColor: tcell.ColorOrangeRed,
		BgColor:       tcell.ColorRed,
		FgColor:       tcell.ColorWhite,
		ButtonColor:   tcell.ColorPurple,
	},
}
