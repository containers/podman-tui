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
		FgColor: tcell.ColorLightCyan,
		BgColor: tcell.ColorLightSkyBlue,
		HeaderRow: headerRow{
			FgColor: tcell.ColorBlack,
			BgColor: tcell.ColorSteelBlue,
		},
	},
	InfoBar: infoBar{
		ItemFgColor:  tcell.ColorLightSkyBlue,
		ValueFgColor: tcell.ColorWhite,
		ProgressBar: progressBar{
			FgColor:       tcell.ColorWhite,
			BarEmptyColor: tcell.ColorWhite,
			BarOKColor:    tcell.ColorGreen,
			BarWarnColor:  tcell.ColorOrange,
			BarCritColor:  tcell.ColorRed,
		},
	},
	Menu: menu{
		FgColor: tcell.ColorWhite,
		BgColor: tcell.ColorBlack,
		Item: menuItem{
			FgColor: tcell.ColorBlack,
			BgColor: tcell.ColorSteelBlue,
		},
	},
	Help: help{
		BorderColor:   tcell.ColorLightSkyBlue,
		BgColor:       tview.Styles.PrimitiveBackgroundColor,
		FgColor:       tcell.ColorWhite,
		HeaderFgColor: tcell.ColorSteelBlue,
	},
	ConnectionProgressDialog: connectionProgressDialog{
		BgColor:     tcell.ColorOrangeRed,
		FgColor:     tcell.ColorOrange,
		PrgBarColor: tcell.ColorOrange,
		BorderColor: tcell.ColorLightSkyBlue,
		TitleColor:  tcell.ColorLightCyan,
	},
	CommandDialog: commandDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorNavy,
		},
	},
	TopDialog: topDialog{
		BgColor:                tcell.ColorSteelBlue,
		FgColor:                tcell.ColorBlack,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorWhite,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorSteelBlue,
		},
	},
	ConfirmDialog: confirmDialog{
		BgColor: tcell.ColorOrange,
		FgColor: tcell.ColorBlack,
	},
	ImageSearchDialog: imageSearchDialog{
		BgColor:                tcell.ColorSteelBlue,
		FgColor:                tcell.ColorWhite,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorWhite,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorSteelBlue,
		},
	},
	ImageHistoryDialog: imageHistoryDialog{
		BgColor:                tcell.ColorSteelBlue,
		FgColor:                tcell.ColorBlack,
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorWhite,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorSteelBlue,
		},
	},
	ImageImportDialog: imageImportDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	ImageBuildDialog: imageBuildDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	ImageBuildPrgDialog: imageBuildPrgDialog{
		BgColor:      tcell.ColorSteelBlue,
		FgColor:      tcell.ColorWhite,
		PrgCellColor: tcell.ColorOrange,
		Terminal: terminal{
			BgColor: tview.Styles.PrimitiveBackgroundColor,
			FgColor: tcell.ColorWhite,
		},
	},
	ImageSaveDialog: imageSaveDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	VolumeCreateDialog: volumeCreateDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	NetworkCreateDialog: networkCreateDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	ContainerCreateDialog: containerCreateDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	ContainerExecDialog: containerExecDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	ContainerExecTerminalDialog: containerExecTerminalDialog{
		BgColor:       tcell.ColorSteelBlue,
		FgColor:       tcell.ColorWhite,
		HeaderBgColor: tcell.ColorNavy,
		Terminal: terminal{
			BgColor: tcell.NewRGBColor(0, 0, 0),
			FgColor: tcell.NewRGBColor(255, 255, 255),
		},
	},
	ContainerStatsDialog: containerStatsDialog{
		ResultTableBgColor:     tview.Styles.PrimitiveBackgroundColor,
		ResultTableBorderColor: tcell.ColorWhite,
		TableHeaderFgColor:     tcell.ColorLightSkyBlue,
		BgColor:                tcell.ColorSteelBlue,
		FgColor:                tcell.ColorWhite,
	},
	PodCreateDialog: podCreateDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	PodStatsDialog: podStatsDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	DropdownStyle: dropdownStyle{
		Unselected: tcell.StyleDefault.Background(tcell.ColorLightSkyBlue).Foreground(tcell.ColorBlack),
		Selected:   tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
	},
	EventsDialog: eventsDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
		EventViewer: terminal{
			BgColor: tview.Styles.PrimitiveBackgroundColor,
			FgColor: tview.Styles.PrimaryTextColor,
		},
	},
	ConnectionAddDialog: connectionAddDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	DiskUageDialog: diskUageDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorNavy,
		},
	},
	MessageDialog: messageDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
		Terminal: terminal{
			BgColor: tview.Styles.PrimitiveBackgroundColor,
			FgColor: tview.Styles.PrimaryTextColor,
		},
	},
	ButtonPrimitive: buttonPrimitive{
		BgColor: tcell.ColorBlue,
	},
	InputFieldPrimitive: inputFieldPrimitive{
		BgColor: tcell.ColorBlue,
	},
	ProgressDailog: progressDailog{
		PgBarColor: tcell.ColorOrange,
	},
	InputDialog: inputDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
	},
	ErrorDialog: errorDialog{
		HeaderFgColor: tcell.ColorDarkRed,
		BgColor:       tcell.ColorOrangeRed,
		FgColor:       tcell.ColorWhite,
	},
}
