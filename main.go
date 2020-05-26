package main

import (
	"flash-sync-server/enums"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	app *App = NewApp()
)

func init() {

}

func main() {

	// 发送数据连接包
	ticker := time.NewTicker(5 * time.Second)
	go sendConnectPack(ticker)

	// 读取所有设备号
	iter := app.Db.NewIterator(util.BytesPrefix([]byte("devices-")), nil)
	for iter.Next() {
		// Use key/value.
		app.ClientDevices[string(iter.Key())] = string(iter.Value())
	}
	iter.Release()

	//logTicker := time.NewTicker(3 * time.Second)
	go func() {

		//for t := range logTicker.C {
		//
		//	App.MainWindow.Synchronize(func() {
		//		trackLatest := lb.ItemVisible(len(lb.Model())-1) && len(lb.SelectedIndexes()) <= 1
		//
		//		model.items = append(model.items, logEntry{time.Now(), "Some new stuff."})
		//		index := len(model.items) - 1
		//		model.PublishItemsInserted(index, index)
		//
		//		if trackLatest {
		//			lb.EnsureItemVisible(len(model.items) - 1)
		//		}
		//	})
		//	logModel.PushLog(t.String())
		//
		//	err := lb.SetModel(logModel)
		//	log.Println("push items", err)
		//}
	}()

	// 程序主窗口运行
	runMainWindow()
}

func sendConnectPack(ticker *time.Ticker) {

	udpPort, tcpPort := app.Config.Udp.Port, app.Config.Tcp.Port

	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: udpPort}

	broadcast, err := net.ListenUDP("udp", srcAddr)
	if err != nil {

		panic(err)
	}

	log.Printf("start udp broadcast, udp port: %d", udpPort)

	for _ = range ticker.C {

		// 广播自己的 tcp 端口
		msg := fmt.Sprintf("hello ! my tcp port=[%d]", tcpPort)
		_, err := broadcast.WriteToUDP([]byte(msg), dstAddr)
		log.Printf(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func runMainWindow() {

	var pathLabel *walk.Label

	_, err := MainWindow{
		AssignTo: &app.MainWindow,
		Title:    app.I18n.Tr("app_name"),
		Icon:     "assets/icons/app.png",
		Size:     Size{Width: 400, Height: 300},
		Layout:   VBox{},
		MenuItems: []MenuItem{
			Menu{
				Text: app.I18n.Tr("file"),
				Items: []MenuItem{
					Action{
						Text:     app.I18n.Tr("open"),
						Shortcut: Shortcut{walk.ModControl, walk.KeyO},
						OnTriggered: func() {

							fmt.Println("打开文件")
						},
					},
					Separator{},
					Action{
						Text:        app.I18n.Tr("exit"),
						OnTriggered: func() { app.MainWindow.Close() },
					},
				},
			},
			Menu{
				Text:  app.I18n.Tr("language"),
				Items: buildLangMenu(),
			},
			Action{
				Text: app.I18n.Tr("help"),
				OnTriggered: func() {

					err := exec.Command(`cmd`, `/c`, `start`, `https://github.com/seth-shi`).Start()

					if err != nil {

						// TODO
						dialog, err := walk.NewDialog(app.MainWindow)
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
				Text:     app.Config.Data.Path,
			},
			PushButton{
				Text: app.I18n.Tr("select sync path"),
				OnClicked: func() {

					dlg := new(walk.FileDialog)
					dlg.Title = app.I18n.Tr("select sync path")

					if ok, err := dlg.ShowBrowseFolder(app.MainWindow); err != nil {
						fmt.Println(err)
						return
					} else if !ok {
						fmt.Println("no ok")
						return
					}

					// 存储到环境目录
					app.Config.Data.Path = dlg.FilePath
					if nil == app.Config.Save() {
						_ = pathLabel.SetText(app.Config.Data.Path)
					}
				},
			},
			ListBox{
				AssignTo: &app.LogList,
				Model:    app.Logs,
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
			Checked:  app.Config.Language == enums.ZH,
			OnTriggered: func() {

				app.Config.Language = enums.ZH
				if nil == app.Config.Save() {
					_ = zhMenu.SetChecked(true)
					_ = enMenu.SetChecked(false)
					walk.MsgBox(app.MainWindow, app.I18n.Tr("switch success"), app.I18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
				}
			},
		},
		Action{
			AssignTo: &enMenu,
			Text:     enums.EN,
			Checked:  app.Config.Language == enums.EN,
			OnTriggered: func() {

				app.Config.Language = enums.EN
				if nil == app.Config.Save() {
					_ = zhMenu.SetChecked(false)
					_ = enMenu.SetChecked(true)
					walk.MsgBox(app.MainWindow, app.I18n.Tr("switch success"), app.I18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
				}
			},
		},
	}
}
