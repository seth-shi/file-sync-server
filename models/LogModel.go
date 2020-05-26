package models

import (
	"github.com/lxn/walk"
)

type LogModel struct {
	walk.ListModelBase
	items []string
}

func (l *LogModel) PushLog(item string) *LogModel {

	l.items = append(l.items, item)
	return l
}

func NewLogModel() *LogModel {

	return &LogModel{items: make([]string, 100)}
}

func (m *LogModel) ItemCount() int {
	return len(m.items)
}

func (m *LogModel) Value(index int) interface{} {
	return m.items[index]
}
