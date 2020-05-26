package main

import (
	"flash-sync-server/config"
	"flash-sync-server/enums"
	"flash-sync-server/models"
	"time"

	"github.com/iafan/go-l10n/loc"
	"github.com/iafan/go-l10n/locjson"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/ini.v1"

	"github.com/lxn/walk"
)

type App struct {

	/*******************************************
	*  所有的界面 UI 管理
	 */
	MainWindow *walk.MainWindow

	LogList *walk.ListBox

	/*******************************************
	*  所有的数据管理
	 */
	// 数据库操作对象
	Db *leveldb.DB
	// 数据的配置
	Config *config.AppConfig
	// 本地化
	I18n *loc.Context

	// 日志
	Logs []models.LogModel
	// 客户端的设备号
	ClientDevices map[string]string
}

func NewApp() *App {

	// 加载配置文件
	iniPath := "conf.ini"

	appConfig := config.NewAppConfig()
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
	i18n := lp.GetContext(lang)

	db, err := leveldb.OpenFile("data", nil)
	if err != nil {
		panic(err)
	}

	return &App{
		Config:        appConfig,
		I18n:          i18n,
		Db:            db,
		ClientDevices: make(map[string]string),
	}
}
