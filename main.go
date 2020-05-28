package main

import (
	"flash-sync-server/enums"
	. "flash-sync-server/global"
	"flash-sync-server/serveices"
	"fmt"
	"github.com/lxn/win"
	"log"
	"os/exec"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	APP_WIDTH = 400
)

func init() {

	// 每隔 5s 发送一次 udp 数据包
	go serveices.SendConnectUdpPack(time.NewTicker(5 * time.Second))

	App.StartedHandle = func() {

		go func() {

			for l := range App.LogChan {

				App.MainWindow.Synchronize(func() {

					if err := l.PushToView(App.LogView); err != nil {

						println(err)
					}
				})
			}
		}()
	}
}

func main() {

	// 程序主窗口运行
	runMainWindow()
}

func runMainWindow() {

	var pathLabel *walk.Label

	err := MainWindow{
		AssignTo: &App.MainWindow,
		Title:    App.I18n.Tr("app_name"),
		Icon:     "assets/icons/app.png",
		Size:     Size{Width: APP_WIDTH, Height: 300},
		Layout:   VBox{},
		MenuItems: []MenuItem{
			Menu{
				Text: App.I18n.Tr("file"),
				Items: []MenuItem{
					Action{
						Text:     App.I18n.Tr("clear logs"),
						Shortcut: Shortcut{walk.ModControl, walk.KeyC},
						OnTriggered: func() {

							App.MainWindow.Synchronize(func() {

								err := App.LogView.Children().Clear()
								fmt.Println(err)
							})
						},
					},
					Separator{},
					Action{
						Text: App.I18n.Tr("exit"),
						OnTriggered: func() {

							App.MainWindow.Close()
						},
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
			ScrollView{
				AssignTo:        &App.LogView,
				HorizontalFixed: true,
				Alignment:       AlignHNearVNear,
				Layout:          VBox{MarginsZero: true},
				Children:        []Widget{},
			},
		},
	}.Create()

	if err != nil {

		log.Println("exit err", err)
	}

	flag := win.GetWindowLong(App.MainWindow.Handle(), win.GWL_STYLE)
	// fixed size
	flag &= ^win.WS_THICKFRAME
	win.SetWindowLong(App.MainWindow.Handle(), win.GWL_STYLE, flag)

	if App.StartedHandle != nil {
		App.StartedHandle()
	}

	exitCode := App.MainWindow.Run()

	log.Println("exit:", exitCode)
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
				if err := App.Config.Save(); err != nil {

					fmt.Println(err)
					walk.MsgBox(App.MainWindow, App.I18n.Tr("switch fail"), App.I18n.Tr("please reboot soft"), walk.MsgBoxIconError)
					return
				}

				_ = zhMenu.SetChecked(true)
				_ = enMenu.SetChecked(false)
				walk.MsgBox(App.MainWindow, App.I18n.Tr("switch success"), App.I18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
			},
		},
		Action{
			AssignTo: &enMenu,
			Text:     enums.EN,
			Checked:  App.Config.Language == enums.EN,
			OnTriggered: func() {

				App.Config.Language = enums.EN
				if err := App.Config.Save(); err != nil {
					fmt.Println(err)
					walk.MsgBox(App.MainWindow, App.I18n.Tr("switch fail"), App.I18n.Tr("please reboot soft"), walk.MsgBoxIconError)
					return
				}

				_ = zhMenu.SetChecked(false)
				_ = enMenu.SetChecked(true)
				walk.MsgBox(App.MainWindow, App.I18n.Tr("switch success"), App.I18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
			},
		},
	}
}
