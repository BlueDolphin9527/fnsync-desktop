package main

import (
	"log"

	// auto config log output format
	mylog "github.com/cxfksword/fnsync-desktop/logger"

	"github.com/cxfksword/fnsync-desktop/app"
	"github.com/cxfksword/fnsync-desktop/config"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

var (
	// App name
	AppName = "FnSync"
	// Version the version of app (Update by ci build).
	Version = "unknown"
	// Commit the git commit hash of this version (Update by ci build).
	Commit = "unknown"
	// BuildDate the date on which this binary was build (Update by ci build).
	BuildDate = "unknown"
	// Mode the build mode.
	Mode = app.DevMode
)

func main() {
	app.Init(AppName, Version, Commit, BuildDate, Mode)

	runApp := NewApp()

	opts := &options.App{
		Title:             AppName,
		Width:             250,
		Height:            300,
		MinWidth:          250,
		MinHeight:         300,
		HideWindowOnClose: true,
		Mac: &mac.Options{
			WebviewIsTransparent:          true,
			WindowBackgroundIsTranslucent: true,
			TitleBar:                      mac.TitleBarDefault(),
			Menu:                          menu.DefaultMacMenu(),
		},
		DevTools:      false,
		DisableResize: true,
		Startup:       runApp.startup,
		StartHidden:   len(config.App.Devices) > 0,
		Shutdown:      runApp.shutdown,
		Bind: []interface{}{
			runApp,
		},
	}
	if app.IsDebugMode {
		opts.LogLevel = logger.DEBUG
		mylog.SetDebugLevel()
	}
	err := wails.Run(opts)
	if err != nil {
		log.Fatal(err)
	}
}
