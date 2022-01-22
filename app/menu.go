package app

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/utils"

	"github.com/rivo/tview"
)

func newMenu(menuItems [][]string) *tview.TextView {

	menu := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter)
	menu.SetBackgroundColor(utils.Styles.Menu.BgColor)
	var menuList []string
	for i := 0; i < len(menuItems); i++ {
		key, item := genMenuItem(menuItems[i])
		if i == len(menuItems)-1 {
			item = item + " "
		}
		menuList = append(menuList, key+item)
	}
	fmt.Fprintf(menu, "%s", strings.Join(menuList, " "))
	return menu
}

func genMenuItem(items []string) (string, string) {

	key := fmt.Sprintf("[%s:%s:b] <%s>", utils.GetColorName(utils.Styles.Menu.FgColor), utils.GetColorName(utils.Styles.Menu.BgColor), items[0])
	desc := fmt.Sprintf("[%s:%s:b] %s", utils.GetColorName(utils.Styles.Menu.Item.FgColor), utils.GetColorName(utils.Styles.Menu.Item.BgColor), strings.ToUpper(items[1]))

	return key, desc
}
