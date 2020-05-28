package events

import (
	. "flash-sync-server/global"
	"github.com/lxn/walk"
	"os/exec"
)


func ClearLogs() {

	App.MainWindow.Synchronize(func() {

		err := App.LogView.Children().Clear()
		LogErrorHandle(err)
	})
}

func ExitApp()  {

	err := App.MainWindow.Close()
	LogErrorHandle(err)
}

func SelectDataPath()  {

	dlg := new(walk.FileDialog)
	dlg.Title = App.I18n.Tr("select sync path")

	if ok, err := dlg.ShowBrowseFolder(App.MainWindow); err != nil {

		LogErrorHandle(err)
		return

	} else if !ok {
		// 不选择目录,取消的操作
		return
	}

	// 存储到环境目录
	App.Config.Data.Path = dlg.FilePath
	if nil == App.Config.Save() {

		LogErrorHandle(App.DataPathLabel.SetText(App.Config.Data.Path))
	}
}

func OpenHelp()  {

	err := exec.Command(`cmd`, `/c`, `start`, `https://github.com/seth-shi`).Start()

	if err != nil {

		dialog, err := walk.NewDialog(App.MainWindow)
		if err != nil {
			LogErrorHandle(err)
			return
		}

		err = dialog.SetTitle(`https://github.com/seth-shi`)
		if err != nil {
			LogErrorHandle(err)
			return
		}
		dialog.Show()
	}
}

func SwitchLang(lang string) func() {

	return func() {
		
		App.Config.Language = lang
		if err := App.Config.Save(); err != nil {

			LogError(App.I18n.Tr("switch fail"))
			return
		}


		for menuLang, menu := range App.Menus.GetLangMenus() {

			if err := menu.SetChecked(lang == menuLang); err != nil {

				LogErrorHandle(err)
			}
		}

		walk.MsgBox(App.MainWindow, App.I18n.Tr("switch success"), App.I18n.Tr("please reboot soft"), walk.MsgBoxIconInformation)
	}
}