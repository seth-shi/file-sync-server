package config

import (
	"gopkg.in/ini.v1"
)

type AppConfig struct {

	Name string
	StartAt string

	Data Data

	Tcp TcpConfig
}

func (app *AppConfig) SaveToFile(filePath string) error  {

	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, app)
	if err != nil {
		return err
	}

	return cfg.SaveTo(filePath)
}