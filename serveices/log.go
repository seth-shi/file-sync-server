package serveices

import (
	"time"

	"github.com/lxn/walk"
)

type LogModel struct {
	walk.ReflectListModelBase
	items []*logEntry
}

func (m *LogModel) Logs() []*logEntry {
	return m.items
}

func (m *LogModel) Items() interface{} {
	return m.items
}

type logEntry struct {
	timestamp time.Time

	messageType    string
	messageContent string
}

func InfoLog(content string) *logEntry {

	return &logEntry{time.Now(), "info", content}
}

func (m *LogModel) ItemCount() int {
	return len(m.items)
}

func (m *LogModel) Value(index int) interface{} {
	return m.items[index]
}
