package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/cxfksword/fnsync-desktop/client"
	"github.com/cxfksword/fnsync-desktop/config"
	"github.com/cxfksword/fnsync-desktop/msg"
	"github.com/cxfksword/fnsync-desktop/utils"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/mac"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

// App application struct
type App struct {
	runtime *wails.Runtime

	startsAtLoginMenu *menu.MenuItem
}

// NewApp creates a new Basic application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (app *App) startup(runtime *wails.Runtime) {
	// Perform your setup here
	app.runtime = runtime

	app.startsAtLoginMenu = &menu.MenuItem{
		Label:   "开机自启动",
		Type:    menu.CheckboxType,
		Checked: false,
		Click:   app.updateStartOnLogin,
	}
	startsAtLogin, err := mac.StartsAtLogin()
	if err != nil {
		app.startsAtLoginMenu.Label = "开机自启动(不支持)"
		app.startsAtLoginMenu.Disabled = true
	} else {
		app.startsAtLoginMenu.Checked = startsAtLogin
	}
	app.initMenus(runtime)

	go msg.UIMsgHandler.StartMsgLoop()
	go client.SleepNotifier.StartSubscribe()
	// start listen client connect
	go client.Listener.StartAccept()

}

func (app *App) initMenus(runtime *wails.Runtime) {
	var items []*menu.MenuItem

	items = append(items, &menu.MenuItem{
		Type:     menu.TextType,
		Label:    "已连接的设备:",
		Disabled: true,
	})
	if len(config.App.Devices) > 0 {
		for _, v := range config.App.Devices {
			items = append(items, &menu.MenuItem{
				Type:  menu.TextType,
				Label: v.Name,
			})
		}
	} else {
		items = append(items, &menu.MenuItem{
			Type:     menu.TextType,
			Label:    "(无)",
			Disabled: true,
		})
	}
	items = append(items, menu.Separator())
	items = append(items, &menu.MenuItem{
		Type:  menu.TextType,
		Label: "连接其他设备",
		Click: app.onConnectMenuClicked,
	})
	items = append(items, &menu.MenuItem{
		Label:   "同步剪贴板到手机",
		Type:    menu.CheckboxType,
		Checked: config.App.ClipboardSync,
		Click:   app.onClipboardSyncMenuClicked,
	})
	items = append(items, app.startsAtLoginMenu)
	// items = append(items, &menu.MenuItem{
	// 	Type:     menu.TextType,
	// 	Label:    "设置",
	// 	Disabled: true,
	// })
	items = append(items, menu.Separator())
	items = append(items, &menu.MenuItem{
		Type:        menu.TextType,
		Label:       "退出",
		Accelerator: keys.CmdOrCtrl("q"),
		Click:       app.onQuitMenuClicked,
	})

	m := &menu.Menu{Items: items}
	runtime.Menu.SetTrayMenu(&menu.TrayMenu{
		Image: "icon",
		Menu:  m,
	})
}

func (app *App) onQuitMenuClicked(_ *menu.CallbackData) {
	app.runtime.Quit()
}

func (app *App) onConnectMenuClicked(_ *menu.CallbackData) {
	app.runtime.Window.Show()
}

func (app *App) onClipboardSyncMenuClicked(_ *menu.CallbackData) {
	config.App.ClipboardSync = !config.App.ClipboardSync
	config.App.Save()
}

func (app *App) updateStartOnLogin(data *menu.CallbackData) {
	err := mac.StartAtLogin(data.MenuItem.Checked)
	if err != nil {
		app.startsAtLoginMenu.Label = "开机自启动(不支持)"
		app.startsAtLoginMenu.Disabled = true
	}
	// We need to refresh all as the menuitem is used in multiple places.
	// If we don't refresh, only the menuitem clicked will toggle in the UI.
	// app.refreshMenus()
}

// shutdown is called at application termination
func (app *App) shutdown() {
	client.Listener.Terminate()

	// TODO: wait for all connection terminate
	time.Sleep(200 * time.Millisecond)

	// Perform your teardown here
	log.Info().Msg("Shutdown app")
}

func (app *App) GenerateQRCode() string {
	client.Listener.RefreshCode()
	helloMsg := msg.Builder.MakeHello(client.Listener.GetCode())

	helloJson := utils.ToJSON(helloMsg)
	png, _ := qrcode.Encode(string(helloJson), qrcode.Medium, -10)
	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(png))
}
