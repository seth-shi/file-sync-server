package main

import (
	"flash-sync-server/config"
	"flash-sync-server/enums"
	"fmt"
	"github.com/iafan/go-l10n/loc"
	"github.com/iafan/go-l10n/locjson"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"gopkg.in/ini.v1"
	"time"
)

var appConfig = &config.AppConfig{}
var i18n *loc.Context
var mw *walk.MainWindow

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
}

func main() {

	var pathLabel *walk.Label

	MainWindow{
		AssignTo: &mw,
		Title:   i18n.Tr("app_name"),
		Icon: "assets/icons/app.png",
		Size:    Size{Width: 400, Height: 300},
		Layout: VBox{},
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
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text: "About",
						OnTriggered: func() {

							fmt.Println("帮助")
						},
					},
				},
			},
		},
		Children: []Widget{
			Label{
				AssignTo: &pathLabel,
				Text: appConfig.Data.Path,
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
