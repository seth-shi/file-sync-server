package main

import (
	"flash-sync-server/enums"
	"flash-sync-server/serveices"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func init() {

}

func main() {

	// 每隔 5s 发送一次 udp 数据包
	go serveices.SendConnectUdpPack(time.NewTicker(5 * time.Second))

	logTicker := time.NewTicker(3 * time.Second)
	go func() {

		for t := range logTicker.C {

			App.MainWindow.Synchronize(func() {

				index := App.Log.ItemCount() - 1

				trackLatest := App.LogList.ItemVisible(index) &&
					len(App.LogList.SelectedIndexes()) <= 1

				items := append(App.Log.Logs(), serveices.InfoLog("写入日志"))
				App.Log.PublishItemsInserted(index, index)

				if trackLatest {
					App.LogList.EnsureItemVisible(len(items) - 1)
				}
			})

			log.Println("push items", t)
		}
	}()

	// 程序主窗口运行
	runMainWindow()
}

func runMainWindow() {

	var pathLabel *walk.Label

	_, err := MainWindow{
		AssignTo: &App.MainWindow,
		Title:    App.I18n.Tr("app_name"),
		Icon:     "assets/icons/app.png",
		Size:     Size{Width: 400, Height: 300},
		Layout:   VBox{},
		MenuItems: []MenuItem{
			Menu{
				Text: App.I18n.Tr("file"),
				Items: []MenuItem{
					Action{
						Text:     App.I18n.Tr("open"),
						Shortcut: Shortcut{walk.ModControl, walk.KeyO},
						OnTriggered: func() {

							fmt.Println("打开文件")
						},
					},
					Separator{},
					Action{
						Text:        App.I18n.Tr("exit"),
						OnTriggered: func() { App.MainWindow.Close() },
					},
				},
			},
			Menu{
				Text:  App.I18n.Tr("language"),
				Items: buildLangMenu(),
			},
			Action{
				Text: App.I18n.Tr("help"),
				OnTriggered: func() {

					err := exec.Command(`cmd`, `/c`, `start`, `https://github.com/seth-shi`).Start()

					if err != nil {

						// TODO
						dialog, err := walk.NewDialog(App.MainWindow)
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
				Text:     App.Config.Data.Path,
			},
			PushButton{
				Text: App.I18n.Tr("select sync path"),
				OnClicked: func() {

					dlg := new(walk.FileDialog)
					dlg.Title = App.I18n.Tr("select sync path")

					if ok, err := dlg.ShowBrowseFolder(App.MainWindow); err != nil {
						fmt.Println(err)
						return
					} else if !ok {
						fmt.Println("no ok")
						return
					}

					// 存储到环境目录
					App.Config.Data.Path = dlg.FilePath
					if nil == App.Config.Save() {
						_ = pathLabel.SetText(App.Config.Data.Path)
					}
				},
			},
			ListBox{
				AssignTo: &App.LogList,
				Model:    App.Logs,
			},
		},
	}.Run()

	if err != nil {

		log.Println("exit err", err)
	}
}

func buildLangMenu() []MenuItem {

	var zhMenu, enMenu *walk.Action

	return []MenuItem{
		Action{
			AssignTo: &zhMenu,
			Text:     enums.ZH,
			Checked:  App.Config.Language == enums.ZH,
			OnTriggered: func() {

				App.Config.Language = enums.ZH
				if nil == App.Config.Save() {
					_ = zhMenu.SetChecked(true)
					_ = enMenu.SetChecked(false)
					walk.MsgBox(App.MainWindow, App.I18n.Tr("switch success"), App.I18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
				}
			},
		},
		Action{
			AssignTo: &enMenu,
			Text:     enums.EN,
			Checked:  App.Config.Language == enums.EN,
			OnTriggered: func() {

				App.Config.Language = enums.EN
				if nil == App.Config.Save() {
					_ = zhMenu.SetChecked(false)
					_ = enMenu.SetChecked(true)
					walk.MsgBox(App.MainWindow, App.I18n.Tr("switch success"), App.I18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
				}
			},
		},
	}
}
