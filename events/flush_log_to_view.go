package events

import . "flash-sync-server/global"

func FlushLogToView() {

	for l := range App.LogChan {

		lg := l

		App.MainWindow.Synchronize(func() {

			if err := lg.PushToView(App.LogView); err != nil {

				println(err)
			}
		})
	}
}
