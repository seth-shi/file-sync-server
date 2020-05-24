package main

import (
	"flash-sync-server/config"
	"flash-sync-server/enums"
	"fmt"
	"os/exec"
	"time"

	"github.com/firstrow/tcp_server"
	"github.com/iafan/go-l10n/loc"
	"github.com/iafan/go-l10n/locjson"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"gopkg.in/ini.v1"
)

var appConfig = &config.AppConfig{}
var i18n *loc.Context
var mw *walk.MainWindow
var db *leveldb.DB

func init() {

	// 加载配置文件
	iniPath := "conf.ini"
	err := ini.MapTo(appConfig, iniPath)
	if err != nil {
		panic(err)
	}

	appConfig.StartAt = time.Now().Format("2006-01-02 15:04:05")
	err = appConfig.SetSavePath(iniPath).Save()
	if err != nil {
		panic(err)
	}

	// i18n
	lang := enums.ZH
	if appConfig.Language != enums.ZH {
		lang = enums.EN
	}
	lp := loc.NewPool(lang)
	lp.Resources[enums.ZH] = locjson.Load("translates/zh.json")
	lp.Resources[enums.EN] = locjson.Load("translates/en.json")
	i18n = lp.GetContext(lang)

	db, err = leveldb.OpenFile("data", nil)
	if err != nil {
		panic(err)
	}
}

func main() {

	var pathLabel *walk.Label
	// 我的设备
	iter := db.NewIterator(util.BytesPrefix([]byte("devices-")), nil)
	devices := make(map[string]string)
	for iter.Next() {
		// Use key/value.
		devices[string(iter.Key())] = string(iter.Value())
	}
	iter.Release()

	server := tcp_server.New("localhost:" + appConfig.Tcp.Port)
	server.OnNewClient(func(c *tcp_server.Client) {
		// new client connected
		// lets send some message
		c.Send("Hello")
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// new message received
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
	})

	go server.Listen()

	fmt.Println("server")
	MainWindow{
		AssignTo: &mw,
		Title:    i18n.Tr("app_name"),
		Icon:     "assets/icons/app.png",
		Size:     Size{Width: 400, Height: 300},
		Layout:   VBox{},
		MenuItems: []MenuItem{
			Menu{
				Text: i18n.Tr("file"),
				Items: []MenuItem{
					Action{
						Text:     i18n.Tr("open"),
						Shortcut: Shortcut{walk.ModControl, walk.KeyO},
						OnTriggered: func() {

							fmt.Println("打开文件")
						},
					},
					Separator{},
					Action{
						Text:        i18n.Tr("exit"),
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text:  i18n.Tr("language"),
				Items: buildLangMenu(),
			},
			Action{
				Text: i18n.Tr("help"),
				OnTriggered: func() {

					err := exec.Command(`cmd`, `/c`, `start`, `https://github.com/seth-shi`).Start()

					if err != nil {

						// TODO
						dialog, err := walk.NewDialog(mw)
						if err != nil {
							panic(err)
						}
						dialog.SetTitle(`https://github.com/seth-shi`)
						dialog.Show()
					}
				},
			},
		},
		Children: []Widget{
			Label{
				AssignTo: &pathLabel,
				Text:     appConfig.Data.Path,
			},
			PushButton{
				Text: i18n.Tr("select sync path"),
				OnClicked: func() {

					dlg := new(walk.FileDialog)
					dlg.Title = i18n.Tr("select sync path")

					if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
						fmt.Println(err)
						return
					} else if !ok {
						fmt.Println("no ok")
						return
					}

					// 存储到环境目录
					appConfig.Data.Path = dlg.FilePath
					if nil == appConfig.Save() {
						_ = pathLabel.SetText(appConfig.Data.Path)
					}
				},
			},
		},
	}.Run()
}

func buildLangMenu() []MenuItem {

	var zhMenu, enMenu *walk.Action

	return []MenuItem{
		Action{
			AssignTo: &zhMenu,
			Text:     enums.ZH,
			Checked:  appConfig.Language == enums.ZH,
			OnTriggered: func() {

				appConfig.Language = enums.ZH
				if nil == appConfig.Save() {
					_ = zhMenu.SetChecked(true)
					_ = enMenu.SetChecked(false)
					walk.MsgBox(mw, i18n.Tr("switch success"), i18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
				}
			},
		},
		Action{
			AssignTo: &enMenu,
			Text:     enums.EN,
			Checked:  appConfig.Language == enums.EN,
			OnTriggered: func() {

				appConfig.Language = enums.EN
				if nil == appConfig.Save() {
					_ = zhMenu.SetChecked(false)
					_ = enMenu.SetChecked(true)
					walk.MsgBox(mw, i18n.Tr("switch success"), i18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
				}
			},
		},
	}
}
