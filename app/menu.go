package app

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
)

func newMenu(menuItems [][]string) *tview.TextView {
	menu := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter)

	menu.SetBackgroundColor(style.BgColor)

	var menuList []string

	for i := 0; i < len(menuItems); i++ {
		key, item := genMenuItem(menuItems[i])
		if i == len(menuItems)-1 {
			item += " "
		}

		menuList = append(menuList, key+item)
	}

	fmt.Fprintf(menu, "%s", strings.Join(menuList, " "))

	return menu
}

func genMenuItem(items []string) (string, string) {
	key := fmt.Sprintf("[%s::b] <%s>[-:-:-]", style.GetColorHex(style.PageHeaderFgColor), items[0])
	desc := fmt.Sprintf("[%s:%s:b] %s [-:-:-]",
		style.GetColorHex(style.PageHeaderFgColor),
		style.GetColorHex(style.MenuBgColor),
		strings.ToUpper(items[1]))

	return key, desc
}
