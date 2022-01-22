package utils

import (
	"github.com/gdamore/tcell/v2"
)

// Styles represent default application style
var Styles = theme{
	PageTable: pageTable{
		FgColor: tcell.ColorLightCyan,
		BgColor: tcell.ColorSteelBlue,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
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
	CommandDialog: commandDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorWhite,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorNavy,
		},
	},
	ConfirmDialog: confirmDialog{
		BgColor: tcell.ColorOrange,
		FgColor: tcell.ColorBlack,
	},
	ImageSearchDialog: imageSearchDialog{
		BgColor:                tcell.ColorSteelBlue,
		FgColor:                tcell.ColorWhite,
		ResultTableBgColor:     tcell.ColorSteelBlue,
		ResultTableBorderColor: tcell.ColorNavy,
		ResultHeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorNavy,
		},
	},
	ImageHistoryDialog: imageHistoryDialog{
		BgColor: tcell.ColorSteelBlue,
		FgColor: tcell.ColorBlack,
		HeaderRow: headerRow{
			FgColor: tcell.ColorWhite,
			BgColor: tcell.ColorNavy,
		},
	},
}
