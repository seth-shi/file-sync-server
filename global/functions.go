package global

import (
	"flash-sync-server/models"
	"math/rand"
)

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


var defaultLetters = []byte("abcdefghijklmnopqrstuvwxyz0123456789")

// RandomString returns a random string with a fixed length
func RandomString(n int) string {
	var letters []byte

	letters = defaultLetters

	b := make([]byte, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}