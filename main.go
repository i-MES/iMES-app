package main

import (
	"embed"
	"runtime"

	imes "github.com/i-mes/imes-app/backend"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	api := &imes.Api{}
	app := &App{}
	app.api = api

	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Load Config Data Only", keys.CmdOrCtrl("o"), app.loadConfigFileCallback)
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		wailsruntime.Quit(nil)
	})
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
	}
	ViewMenu := AppMenu.AddSubmenu("View")
	ViewMenu.AddText("...", keys.CmdOrCtrl(","), app.msgDialog)
	HelpMenu := AppMenu.AddSubmenu("Help")
	HelpMenu.AddText("...", keys.CmdOrCtrl("."), app.msgDialog)

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "iMES-app",
		Width:             1024,
		Height:            768,
		MinWidth:          900,
		MinHeight:         600,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		Assets:            assets,
		Menu:              AppMenu,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnBeforeClose:     app.beforeClose,
		OnShutdown:        app.shutdown,
		WindowStartState:  options.Normal,
		Bind: []interface{}{
			app,
			api, // bind a struct(not a file), and the struct's method( func (* xxx)XXX )
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarDefault(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "intelligent Manufacturing Execution System",
				Message: "专为中小制造企业服务的智能 MES 系统",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err)
	}
}
