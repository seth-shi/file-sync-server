package main

import (
	"flash-sync-server/config"
	"flash-sync-server/enums"
	"flash-sync-server/models"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"

	"github.com/iafan/go-l10n/loc"
	"github.com/iafan/go-l10n/locjson"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"gopkg.in/ini.v1"
)

var (
	appConfig = config.NewAppConfig()
	i18n      *loc.Context
	mw        *walk.MainWindow
	lb		  *walk.ListBox
	db        *leveldb.DB

	logModel      = models.NewLogModel()
	clientDevices = make(map[string]string)
)

func init() {

	// 加载配置文件
	iniPath := "conf.ini"
	err := ini.MapTo(appConfig, iniPath)
	if err != nil {
		panic(err)
	}

	// 不是开发环境才写入启动时间
	if !appConfig.Environment("dev") {

		appConfig.StartAt = time.Now().Format("2006-01-02 15:04:05")
		err = appConfig.SetSavePath(iniPath).Save()
		if err != nil {
			panic(err)
		}
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

	// 发送数据连接包
	ticker := time.NewTicker(5 * time.Second)
	go sendConnectPack(ticker)

	// 读取所有设备号
	iter := db.NewIterator(util.BytesPrefix([]byte("devices-")), nil)
	for iter.Next() {
		// Use key/value.
		clientDevices[string(iter.Key())] = string(iter.Value())
	}
	iter.Release()

	logTicker := time.NewTicker(3 * time.Second)
	go func() {

		for t := range logTicker.C {

			logModel.PushLog(t.String())

			err := lb.SetModel(logModel)
			log.Println("push items", err)
		}
	}()


	// 程序主窗口运行
	runMainWindow()
}

func sendConnectPack(ticker *time.Ticker) {

	udpPort, tcpPort := appConfig.Udp.Port, appConfig.Tcp.Port

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
			ListBox{
				AssignTo: &lb,
				Model: logModel,
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
