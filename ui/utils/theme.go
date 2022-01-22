package utils

import (
	"github.com/gdamore/tcell/v2"
)

type theme struct {
	InfoBar            infoBar
	Menu               menu
	PageTable          pageTable
	CommandDialog      commandDialog
	ConfirmDialog      confirmDialog
	ImageSearchDialog  imageSearchDialog
	ImageHistoryDialog imageHistoryDialog
}

type infoBar struct {
	ItemFgColor  tcell.Color
	ValueFgColor tcell.Color
	ProgressBar  progressBar
}
type menu struct {
	FgColor tcell.Color
	BgColor tcell.Color
	Item    menuItem
}
type menuItem struct {
	FgColor tcell.Color
	BgColor tcell.Color
}

type progressBar struct {
	FgColor       tcell.Color
	BarOKColor    tcell.Color
	BarWarnColor  tcell.Color
	BarCritColor  tcell.Color
	BarEmptyColor tcell.Color
}

type pageTable struct {
	FgColor   tcell.Color
	BgColor   tcell.Color
	HeaderRow headerRow
}

type headerRow struct {
	FgColor tcell.Color
	BgColor tcell.Color
}

type commandDialog struct {
	BgColor   tcell.Color
	FgColor   tcell.Color
	HeaderRow headerRow
}

type confirmDialog struct {
	BgColor tcell.Color
	FgColor tcell.Color
}

type imageSearchDialog struct {
	BgColor                tcell.Color
	FgColor                tcell.Color
	ResultHeaderRow        headerRow
	ResultTableBgColor     tcell.Color
	ResultTableBorderColor tcell.Color
}

type imageHistoryDialog struct {
	BgColor   tcell.Color
	FgColor   tcell.Color
	HeaderRow headerRow
}
