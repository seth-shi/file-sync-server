package config

import (
	"gopkg.in/ini.v1"
)

type AppConfig struct {
	Name     string
	Env      string
	Language string
	StartAt  string

	Data DataConfig

	Udp UdpConfig
	Tcp TcpConfig

	// 存储目录，和配置无关
	savePath string
}

func NewAppConfig() *AppConfig {

	return &AppConfig{
		Env: "dev",
	}
}

func (app *AppConfig) Environment(env string) bool {

	return env == app.Env
}

func (app *AppConfig) SetSavePath(path string) *AppConfig {
	app.savePath = path
	return app
}

func (app *AppConfig) Save() error {

	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, app)
	if err != nil {
		return err
	}

	return cfg.SaveTo(app.savePath)
}
