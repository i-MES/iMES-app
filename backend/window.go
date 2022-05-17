package imes

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

func getMenu() *menu.Menu {
	return menu.NewMenuFromItems(
		menu.SubMenu("File", menu.NewMenuFromItems(
			menu.Text("&Open", keys.CmdOrCtrl("o"), func(data *menu.CallbackData) {
				fmt.Println(data)
			}),
			menu.Separator(),
		)),
	)
	// runtime.MenuSetApplicationMenu(imesMenu)
}