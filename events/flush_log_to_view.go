package events

import . "flash-sync-server/global"

func FlushLogToView() {

	for l := range App.LogChan {

		App.MainWindow.Synchronize(func() {

			if err := l.PushToView(App.LogView); err != nil {

				println(err)
			}
		})
	}
}
