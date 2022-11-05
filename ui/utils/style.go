//go:build !windows
// +build !windows

package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Styles represent default application style
var Styles = theme{
	PageTable: pageTable{
		FgColor: tcell.ColorWhiteSmoke,
		BgColor: tcell.ColorMediumPurple,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhiteSmoke,
			BgColor: tcell.ColorRebeccaPurple,
		},
	},
	InfoBar: infoBar{
		ItemFgColor:  tcell.ColorSilver,
		ValueFgColor: tcell.ColorWhiteSmoke,
		ProgressBar: progressBar{
			FgColor:       tcell.ColorWhite,
			BarEmptyColor: tcell.ColorWhite,
			BarOKColor:    tcell.ColorGreen,
			BarWarnColor:  tcell.ColorOrange,
			BarCritColor:  tcell.ColorRed,
		},
	},
	Menu: menu{
		FgColor: tcell.ColorWhiteSmoke,
		BgColor: tcell.ColorBlack,
		Item: menuItem{
			FgColor: tcell.ColorWhiteSmoke,
			BgColor: tcell.ColorRebeccaPurple,
		},
	},
	Help: help{
		BorderColor:   tcell.ColorMediumPurple,
		BgColor:       tview.Styles.PrimitiveBackgroundColor,
		FgColor:       tcell.ColorWhiteSmoke,
		HeaderFgColor: tcell.ColorSilver,
	},
	ConnectionProgressDialog: connectionProgressDialog{
		BgColor:     tcell.ColorMediumPurple,
		FgColor:     tcell.ColorRebeccaPurple,
		PrgBarColor: tcell.ColorDarkOrange,
		BorderColor: tcell.ColorWhiteSmoke,
		TitleColor:  tcell.ColorWhiteSmoke,
	},
	CommandDialog: commandDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhiteSmoke,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhiteSmoke,
			BgColor: tcell.ColorRebeccaPurple,
		},
	},
	TopDialog: topDialog{
		BgColor:                tcell.ColorMediumPurple,
		FgColor:                tcell.ColorBlack,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorMediumPurple,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorRebeccaPurple,
		},
	},
	ConfirmDialog: confirmDialog{
		BgColor:     tcell.ColorMediumPurple,
		FgColor:     tcell.ColorWhiteSmoke,
		ButtonColor: tcell.ColorRebeccaPurple,
	},
	ImageSearchDialog: imageSearchDialog{
		BgColor:                tcell.ColorMediumPurple,
		FgColor:                tcell.ColorWhite,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorMediumPurple,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorRebeccaPurple,
		},
	},
	ImageHistoryDialog: imageHistoryDialog{
		BgColor:                tcell.ColorMediumPurple,
		FgColor:                tcell.ColorBlack,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorMediumPurple,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorRebeccaPurple,
		},
	},
	ImageImportDialog: imageImportDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ImageBuildDialog: imageBuildDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ImageBuildPrgDialog: imageBuildPrgDialog{
		BgColor:      tcell.ColorMediumPurple,
		FgColor:      tcell.ColorWhite,
		PrgCellColor: tcell.ColorDarkOrange,
		Terminal: terminal{
			BorderColor: tcell.ColorMediumPurple,
			BgColor:     tview.Styles.PrimitiveBackgroundColor,
			FgColor:     tcell.ColorWhite,
		},
	},
	ImageSaveDialog: imageSaveDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ImagePushDialog: imagePushDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	VolumeCreateDialog: volumeCreateDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	NetworkCreateDialog: networkCreateDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	NetworkConnectDialog: networkConnectDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ContainerCreateDialog: containerCreateDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ContainerExecDialog: containerExecDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ContainerExecTerminalDialog: containerExecTerminalDialog{
		BgColor:       tcell.ColorMediumPurple,
		FgColor:       tcell.ColorWhite,
		HeaderBgColor: tcell.ColorDarkOrchid,
		Terminal: terminal{
			BorderColor: tcell.ColorMediumPurple,
			BgColor:     tcell.NewRGBColor(0, 0, 0),
			FgColor:     tcell.NewRGBColor(255, 255, 255),
		},
	},
	ContainerStatsDialog: containerStatsDialog{
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorMediumPurple,
		TableHeaderFgColor:     tcell.ColorLightGray,
		BgColor:                tcell.ColorMediumPurple,
		FgColor:                tcell.ColorWhite,
	},
	ContainerCommitDialog: containerCommitDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ContainerCheckpointDialog: containerCheckpointDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ContainerRestoreDialog: containerRestoreDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	PodCreateDialog: podCreateDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	PodStatsDialog: podStatsDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorRebeccaPurple,
		},
	},
	DropdownStyle: dropdownStyle{
		Unselected: tcell.StyleDefault.Background(tcell.ColorMediumOrchid).Foreground(tcell.ColorBlack),
		Selected:   tcell.StyleDefault.Background(tcell.ColorPurple).Foreground(tcell.ColorWhite),
	},
	EventsDialog: eventsDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
		EventViewer: terminal{
			BorderColor: tcell.ColorMediumPurple,
			BgColor:     tview.Styles.PrimitiveBackgroundColor,
			FgColor:     tview.Styles.PrimaryTextColor,
		},
	},
	ConnectionAddDialog: connectionAddDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	DiskUageDialog: diskUageDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorDimGray,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorRebeccaPurple,
		},
	},
	MessageDialog: messageDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
		Terminal: terminal{
			BorderColor: tcell.ColorMediumPurple,
			BgColor:     tview.Styles.PrimitiveBackgroundColor,
			FgColor:     tview.Styles.PrimaryTextColor,
		},
	},
	ButtonPrimitive: buttonPrimitive{
		BgColor: tcell.ColorRebeccaPurple,
	},
	InputFieldPrimitive: inputFieldPrimitive{
		BgColor: tcell.ColorDarkOrchid,
	},
	ProgressDailog: progressDailog{
		PgBarColor: tcell.ColorDarkOrange,
	},
	InputDialog: inputDialog{
		BgColor: tcell.ColorMediumPurple,
		FgColor: tcell.ColorWhite,
	},
	ErrorDialog: errorDialog{
		HeaderFgColor: tcell.ColorWhiteSmoke,
		BgColor:       tcell.ColorIndianRed,
		FgColor:       tcell.ColorWhite,
		ButtonColor:   tcell.ColorDarkRed,
	},
}
