package global

import (
	"flash-sync-server/config"
	"flash-sync-server/enums"
	"flash-sync-server/models"
	"time"

	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/iafan/go-l10n/loc"
	"github.com/iafan/go-l10n/locjson"
	"github.com/lxn/walk"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/ini.v1"
)

var App *application

func init() {

	// 相对工程目录,而不是文件
	iniPath := "./conf.ini"
	dbPath := "./data"

	appConfig := loadIniConfig(iniPath)
	i18n := loadI18nConfig(appConfig.Language)
	db := openDatabase(dbPath)
	devices := loadDevices(db)

	App = &application{
		Config:        appConfig,
		I18n:          i18n,
		Db:            db,
		ClientDevices: devices,
		LogChan:       make(chan *models.LogEntry, 100),

		LinkCode: RandomString(enums.LINK_CODE_LENGTH),
		Menus:    &appMenus{},
	}
}

func loadIniConfig(path string) *config.AppConfig {
	appConfig := config.NewAppConfig()
	err := ini.MapTo(appConfig, path)
	if err != nil {
		panic(err)
	}

	appConfig.SetSavePath(path)

	// 不是开发环境才写入启动时间
	if !appConfig.Environment("dev") {

		appConfig.StartAt = time.Now().Format("2006-01-02 15:04:05")
		err = appConfig.Save()
		if err != nil {
			panic(err)
		}
	}

	return appConfig
}

func loadI18nConfig(defaultLang string) *loc.Context {

	// i18n
	lang := enums.ZH
	if defaultLang != enums.ZH {
		lang = enums.EN
	}
	lp := loc.NewPool(lang)
	lp.Resources[enums.ZH] = locjson.Load("resources/lang/zh.json")
	lp.Resources[enums.EN] = locjson.Load("resources/lang/en.json")
	i18n := lp.GetContext(lang)

	return i18n
}

func openDatabase(path string) *leveldb.DB {

	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}

	return db
}

func loadDevices(db *leveldb.DB) map[string]string {

	devices := make(map[string]string)
	// 读取所有设备号
	iter := db.NewIterator(util.BytesPrefix([]byte("devices-")), nil)
	for iter.Next() {
		// Use key/value.
		devices[string(iter.Key())] = string(iter.Value())
	}
	iter.Release()

	return devices
}

type application struct {

	/*******************************************
	*  所有的界面 UI 管理
	 */
	MainWindow *walk.MainWindow

	// 日志视图
	LogView *walk.ScrollView
	// 验证码, 数据路径
	LinkCodeLabel, DataPathLabel *walk.Label
	// 菜单
	Menus *appMenus
	// 链接验证码的label

	/*******************************************
	*  所有的数据管理
	 */
	// 验证码连接
	LinkCode string

	// 数据库操作对象
	Db *leveldb.DB
	// 数据的配置
	Config *config.AppConfig
	// 本地化
	I18n *loc.Context

	LogChan chan *models.LogEntry

	// 客户端的设备号
	ClientDevices map[string]string

	// 事件
	StartedHandle func()
}
