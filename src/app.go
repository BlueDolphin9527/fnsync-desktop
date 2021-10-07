package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/cxfksword/fnsync-desktop/client"
	"github.com/cxfksword/fnsync-desktop/clipboard"
	"github.com/cxfksword/fnsync-desktop/config"
	"github.com/cxfksword/fnsync-desktop/entity"
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

	startsAtLoginMenu     *menu.MenuItem
	trayMenu              *menu.TrayMenu
	currentSelectedDevice *entity.Device
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

	go app.startUIMsgLoop()
	go client.SleepNotifier.StartSubscribe()
	// start listen client connect
	go client.Listener.StartAccept()

}

func (app *App) initMenus(runtime *wails.Runtime) {
	var items []*menu.MenuItem
	trayImage := "icon"

	items = append(items, &menu.MenuItem{
		Type:     menu.TextType,
		Label:    "连接过的设备:",
		Disabled: true,
	})
	devices := client.Listener.GetDevices()
	if len(devices) > 0 {
		for _, v := range devices {
			if v.IsAlive {
				items = append(items, &menu.MenuItem{
					Type:  menu.TextType,
					Label: v.Name,
					Click: func(data *menu.CallbackData) {
						app.currentSelectedDevice = &v
						app.onDeviceMenuClicked(data)
					},
				})
				trayImage = "icon1"
			} else {
				items = append(items, &menu.MenuItem{
					Type:  menu.TextType,
					Label: fmt.Sprintf("%s (未连接)", v.Name),
					Click: func(data *menu.CallbackData) {
						app.currentSelectedDevice = &v
						app.onTryConnectDeviceMenuClicked(data)
					},
				})
			}

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
	items = append(items, menu.Separator())
	items = append(items, &menu.MenuItem{
		Type:     menu.TextType,
		Label:    "设置:",
		Disabled: true,
	})
	items = append(items, &menu.MenuItem{
		Label:   "同步剪贴板到手机",
		Type:    menu.CheckboxType,
		Checked: config.App.ClipboardSync,
		Click:   app.onClipboardSyncMenuClicked,
	})
	items = append(items, app.startsAtLoginMenu)

	// items = append(items, &menu.MenuItem{
	// 	Type:  menu.TextType,
	// 	Label: "设置",
	// 	Click: app.onSettingMenuClicked,
	// })
	items = append(items, menu.Separator())
	items = append(items, &menu.MenuItem{
		Type:        menu.TextType,
		Label:       "退出",
		Accelerator: keys.CmdOrCtrl("q"),
		Click:       app.onQuitMenuClicked,
	})

	m := &menu.Menu{Items: items}
	if app.trayMenu == nil {
		app.trayMenu = &menu.TrayMenu{
			Image: trayImage,
			Menu:  m,
		}
	} else {
		app.trayMenu.Image = trayImage
		app.trayMenu.Menu = m
	}
	runtime.Menu.SetTrayMenu(app.trayMenu)
}

// refreshMenus refreshes all tray menus to ensure they
// are in sync with the data, EG: checkbox sync
func (app *App) refreshMenus() {
	app.initMenus(app.runtime)
}

func (app *App) onQuitMenuClicked(_ *menu.CallbackData) {
	app.runtime.Quit()
}

func (app *App) onConnectMenuClicked(_ *menu.CallbackData) {
	app.runtime.Window.SetTitle(AppName)
	app.runtime.Events.Emit("tray.click.qrcode")
	app.runtime.Events.Emit("app.qrcode.window.show")
	app.runtime.Window.SetSize(250, 325)
	app.runtime.Window.Show()
}

func (app *App) onClipboardSyncMenuClicked(_ *menu.CallbackData) {
	config.App.ClipboardSync = !config.App.ClipboardSync
	config.App.Save()
}

func (app *App) onDeviceMenuClicked(data *menu.CallbackData) {
	app.runtime.Window.SetTitle("发送文本")
	app.runtime.Events.Emit("tray.click.device")
	app.runtime.Events.Emit("app.device.window.show")
	app.runtime.Window.SetSize(400, 300)
	app.runtime.Window.Show()
}

func (app *App) onTryConnectDeviceMenuClicked(data *menu.CallbackData) {
	client.StartHandshakeDevice(*app.currentSelectedDevice)
}

func (app *App) onSettingMenuClicked(_ *menu.CallbackData) {
	app.runtime.Window.SetTitle("设置")
	app.runtime.Events.Emit("tray.click.settings")
	app.runtime.Events.Emit("app.settings.window.show")
	app.runtime.Window.SetSize(400, 300)
	app.runtime.Window.Show()
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

func (app *App) startUIMsgLoop() {
	app.runtime.Events.On("app.window.close", func(optionalData ...interface{}) {
		app.runtime.Window.Hide()
		app.runtime.Events.Emit("tray.click.qrcode")
	})

	msgCh := msg.UIMsgHandler.Start()
	for {
		msg := <-msgCh

		switch v := msg.(type) {
		case entity.UIUpdateStatusMsg:
			app.refreshMenus()
		case entity.UIScanConnectedMsg:
			app.runtime.Events.Emit("app.scan.connected", v.Code)
		default:
			v.Execute()
		}

	}
}

func (app *App) GenerateQRCode() entity.QRCode {
	client.Listener.RefreshCode()
	helloMsg := msg.Builder.MakeHello(client.Listener.GetCode())

	helloJson := utils.ToJSON(helloMsg)
	log.Debug().Msgf("Hello qrcode json: %s", string(helloJson))
	png, _ := qrcode.Encode(string(helloJson), qrcode.Medium, -10)
	return entity.QRCode{
		Base64Data: fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(png)),
		Data:       helloMsg,
	}
}

func (app *App) GetClipboardText() (string, error) {
	text, err := clipboard.Get()
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	return text, err
}

func (app *App) SendText(txt string) (bool, error) {
	err := client.Listener.SendMeg(app.currentSelectedDevice.Code, txt)
	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	} else {
		msg.UIMsgHandler.Send(entity.UINotifyMsg{
			Title:   AppName,
			Message: "已发送成功",
		})
		return true, nil
	}
}

func (app *App) LoadConfig() (config.AppConfig, error) {
	return *config.App, nil
}

func (app *App) GetStartAtLoginState() (bool, error) {
	startsAtLogin, err := mac.StartsAtLogin()
	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	} else {
		return startsAtLogin, nil
	}
}

func (app *App) ToggleStartAtLogin(checked bool) (bool, error) {
	err := mac.StartAtLogin(checked)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (app *App) ToggleClipboardSync(checked bool) (bool, error) {
	config.App.ClipboardSync = checked
	config.App.Save()
	return true, nil
}
