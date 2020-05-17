package main

import (
	"flash-sync-server/config"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"gopkg.in/ini.v1"
	"strings"
	"time"
)

var appConfig = &config.AppConfig{}

func init()  {

	iniPath := "conf.ini"
	err := ini.MapTo(appConfig, iniPath)
	if err != nil {
		panic(err)
	}

	appConfig.StartAt = time.Now().Format("2006-01-02 15:04:05")
	err = appConfig.SaveToFile(iniPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	var inTE, outTE *walk.TextEdit

	MainWindow{
		Title:   "SCREAMO",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &inTE},
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			PushButton{
				Text: "SCREAM",
				OnClicked: func() {
					outTE.SetText(strings.ToUpper(inTE.Text()))
				},
			},
		},
	}.Run()
}