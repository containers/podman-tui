package help

import (
	"fmt"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpScreen is a help primitive dialog
type HelpScreen struct {
	*tview.Box
	title  string
	layout *tview.Flex
}

// NewHelpScreen returns a help screen primitive.
func NewHelpScreen(appName string, appVersion string) *HelpScreen {
	// returns the help primitive
	help := &HelpScreen{
		Box:   tview.NewBox(),
		title: "help",
	}

	// colors
	headerColor := utils.Styles.HelpScreen.HeaderFgColor
	fgColor := utils.Styles.HelpScreen.FgColor
	bgColor := utils.Styles.HelpScreen.BgColor
	borderColor := utils.Styles.HelpScreen.BorderColor

	// application keys descriotion table
	keyinfo := tview.NewTable()
	keyinfo.SetBackgroundColor(bgColor)
	keyinfo.SetFixed(1, 1)
	keyinfo.SetSelectable(false, false)

	// application description and version text view
	appinfo := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	appinfo.SetBackgroundColor(bgColor)

	licenseInfo := "Released under the Apache License 2.0."
	appInfoText := fmt.Sprintf("%s %s - (C) 2022 podman-tui dev team.\n%s", appName, appVersion, licenseInfo)

	appinfo.SetText(appInfoText)
	appinfo.SetTextColor(headerColor)

	// help table items
	// the items will be devided into two separate tables
	rowIndex := 0
	colIndex := 0
	needInit := true
	maxRowIndex := int(len(utils.UIKeysBindings)/2) + 1
	for i := 0; i < len(utils.UIKeysBindings); i++ {
		if i >= maxRowIndex {
			if needInit {
				colIndex = 2
				rowIndex = 0
				needInit = false
			}
		}
		keyinfo.SetCell(rowIndex, colIndex,
			tview.NewTableCell(fmt.Sprintf("%s:", utils.UIKeysBindings[i].KeyLabel)).
				SetAlign(tview.AlignRight).
				SetBackgroundColor(bgColor).
				SetSelectable(true).SetTextColor(headerColor))
		keyinfo.SetCell(rowIndex, colIndex+1,
			tview.NewTableCell(utils.UIKeysBindings[i].KeyDesc).
				SetAlign(tview.AlignLeft).
				SetBackgroundColor(bgColor).
				SetSelectable(true).SetTextColor(fgColor))

		rowIndex = rowIndex + 1
	}

	// appinfo and appkeys layout
	mlayout := tview.NewFlex().SetDirection(tview.FlexRow)
	mlayout.AddItem(appinfo, 2, 0, false)
	mlayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mlayout.AddItem(keyinfo, 0, 1, false)
	mlayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	// layout
	help.layout = tview.NewFlex().SetDirection(tview.FlexColumn)
	help.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	help.layout.AddItem(mlayout, 0, 1, false)
	help.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	help.layout.SetBorder(true)
	help.layout.SetBackgroundColor(bgColor)
	help.layout.SetBorderColor(borderColor)

	return help
}

// GetTitle returns primitive title
func (help *HelpScreen) GetTitle() string {
	return help.title
}

// HasFocus returns whether or not this primitive has focus
func (help *HelpScreen) HasFocus() bool {
	return help.Box.HasFocus() || help.layout.HasFocus()
}

// Focus is called when this primitive receives focus
func (help *HelpScreen) Focus(delegate func(p tview.Primitive)) {
	delegate(help.layout)
}

// Draw draws this primitive onto the screen.
func (help *HelpScreen) Draw(screen tcell.Screen) {

	x, y, width, height := help.Box.GetInnerRect()
	if height <= 3 {
		return
	}
	help.Box.DrawForSubclass(screen, help)
	help.layout.SetRect(x, y, width, height)
	help.layout.Draw(screen)
}
