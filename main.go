package main

import (
	"flash-sync-server/enums"
	"flash-sync-server/events"
	. "flash-sync-server/global"
	"fmt"
	"github.com/lxn/win"
	"log"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	APP_WIDTH = 400
)

func init() {

	// 每隔 5s 发送一次 udp 数据包
	go events.SendConnectUdpPack(time.NewTicker(5 * time.Second))


	App.StartedHandle = func() {

		// 冲刷日志到视图
		go events.FlushLogToView()
	}
}

func main() {

	// 程序主窗口运行
	runMainWindow()
}

func runMainWindow() {

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
				Items: []MenuItem{
					Action{
						AssignTo: &App.Menus.SwitchToZH,
						Text:     enums.ZH,
						Checked:  App.Config.Language == enums.ZH,
						OnTriggered: events.SwitchLang(enums.ZH),
					},
					Action{
						AssignTo: &App.Menus.SwitchToEn,
						Text:     enums.EN,
						Checked:  App.Config.Language == enums.EN,
						OnTriggered: events.SwitchLang(enums.EN),
					},
				},
			},
			Action{
				Text: App.I18n.Tr("help"),
				OnTriggered: events.OpenHelp,
			},
		},
		Children: []Widget{
			Label{
				AssignTo: &App.DataPathLabel,
				Text:     App.Config.Data.Path,
			},
			PushButton{
				Text: App.I18n.Tr("select sync path"),
				OnClicked: events.SelectDataPath,
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
