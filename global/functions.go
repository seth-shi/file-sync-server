package global

import "flash-sync-server/models"

func LogErrorHandle(err error)  {

	if err != nil {

		App.LogChan <- models.ErrorLog(err.Error())
	}
}

func LogError(msg string)  {

	App.LogChan <- models.ErrorLog(msg)
}

func LogWaring(msg string)  {

	App.LogChan <- models.ErrorLog(msg)
}

func LogInfo(msg string)  {

	App.LogChan <- models.InfoLog(msg)
}