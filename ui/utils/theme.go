package utils

import (
	"github.com/gdamore/tcell/v2"
)

type theme struct {
	InfoBar                     infoBar
	Menu                        menu
	PageTable                   pageTable
	Help                        help
	ConnectionProgressDialog    connectionProgressDialog
	CommandDialog               commandDialog
	ConfirmDialog               confirmDialog
	ImageSearchDialog           imageSearchDialog
	ImageHistoryDialog          imageHistoryDialog
	ImageBuildDialog            imageBuildDialog
	ImageBuildPrgDialog         imageBuildPrgDialog
	ImageSaveDialog             imageSaveDialog
	ContainerExecDialog         containerExecDialog
	ContainerExecTerminalDialog containerExecTerminalDialog
	ContainerStatsDialog        containerStatsDialog
	PodStatsDialog              podStatsDialog
	DropdownStyle               dropdownStyle
	EventsDialog                eventsDialog
	ConnectionAddDialog         connectionAddDialog
}

type connectionProgressDialog struct {
	BgColor tcell.Color
	FgColor tcell.Color
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

type help struct {
	BorderColor   tcell.Color
	BgColor       tcell.Color
	FgColor       tcell.Color
	HeaderFgColor tcell.Color
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

type imageBuildDialog struct {
	BgColor tcell.Color
	FgColor tcell.Color
}

type imageBuildPrgDialog struct {
	BgColor      tcell.Color
	FgColor      tcell.Color
	PrgCellColor tcell.Color
	Terminal     terminal
}

type imageSaveDialog struct {
	BgColor tcell.Color
	FgColor tcell.Color
}

type containerExecDialog struct {
	BgColor tcell.Color
	FgColor tcell.Color
}

type containerExecTerminalDialog struct {
	BgColor  tcell.Color
	FgColor  tcell.Color
	Terminal terminal
}

type containerStatsDialog struct {
	TableHeaderFgColor tcell.Color
	BgColor            tcell.Color
	FgColor            tcell.Color
}
type terminal struct {
	BgColor tcell.Color
	FgColor tcell.Color
}

type podStatsDialog struct {
	BgColor tcell.Color
	FgColor tcell.Color
}

type dropdownStyle struct {
	Unselected tcell.Style
	Selected   tcell.Style
}

type eventsDialog struct {
	BgColor     tcell.Color
	FgColor     tcell.Color
	EventViewer terminal
}

type connectionAddDialog struct {
	FgColor tcell.Color
	BgColor tcell.Color
}
