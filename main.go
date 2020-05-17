package main

import (
	"flash-sync-server/config"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"gopkg.in/ini.v1"
	"time"
)

var appConfig = &config.AppConfig{}

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

	// TODO i18n
}

func main() {

	var mw *walk.MainWindow

	MainWindow{
		AssignTo: &mw,
		Title:   "SCREAMO",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						Text:        "&Open",
						Shortcut:    Shortcut{walk.ModControl, walk.KeyO},
						OnTriggered: func() {

							fmt.Println("打开文件")
						},
					},
					Separator{},
					Action{
						Text:        "E&xit",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&View",
				Items: []MenuItem{
					Action{
						Text:    "浏览",
						OnTriggered: func() {

							fmt.Println("浏览")
						},
					},
				},
			},
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text:        "About",
						OnTriggered: func() {

							fmt.Println("帮助")
						},
					},
				},
			},
		},
		Children: []Widget{
			PushButton{
				Text: "SCREAM",
				OnClicked: func() {

					dlg := new(walk.FileDialog)

					dlg.Title = "选择文件夹"

					if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
						fmt.Println(err)
						return
					} else if !ok {
						fmt.Println("no ok")
						return
					}

					// 存储到环境目录
					appConfig.Data.Path = dlg.FilePath
					_ = appConfig.Save()
					fmt.Println(dlg.FilePath)
				},
			},
		},
	}.Run()
}
