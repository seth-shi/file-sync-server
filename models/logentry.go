package models

import (
	"flash-sync-server/enums"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
	"time"
)

const (
	INFO = "info"
	WARING = "waring"
	ERROR = "error"
)

var labelBgColor = SolidColorBrush{Color: walk.RGB(240, 240, 240)}
var labelTextColorMap = map[string]walk.Color{
	INFO: walk.RGB(100, 100, 100),
	WARING: walk.RGB(255, 241, 0),
	ERROR: walk.RGB(153, 0, 51),
}


type LogEntry struct {

	createdAt time.Time

	messageType    string
	messageContent string
}

func NewLogEntry(msgType, content string) *LogEntry {

	return &LogEntry{time.Now(), msgType, content}
}

func InfoLog(content string) *LogEntry {

	return &LogEntry{time.Now(), INFO, content}
}

func WaringLog(content string) *LogEntry {

	return &LogEntry{time.Now(), WARING, content}
}

func ErrorLog(content string) *LogEntry {

	return &LogEntry{time.Now(), ERROR, content}
}

func (l *LogEntry) PushToView(parent walk.Container) error {

	logs := []string{
		l.createdAt.Format("2006-01-02 15:04:05"),
		//l.messageType,
		l.messageContent,
	}


	color, exists := labelTextColorMap[l.messageType]
	if ! exists {
		color, _ = labelTextColorMap[INFO]
	}

	label := Label{
		MinSize:Size{enums.APP_WIDTH - (9*6), 0},
		Alignment: AlignHNearVNear,
		TextColor: color,
		Background: labelBgColor,
		Text: strings.Join(logs, "  "),
	}

	return label.Create(NewBuilder(parent))
}
