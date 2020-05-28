package global

import (
	"flash-sync-server/enums"
	"github.com/lxn/walk"
)

type appMenus struct {

	SwitchToZH, SwitchToEn *walk.Action
}

func (m *appMenus) GetLangMenus() map[string]*walk.Action {

	return map[string]*walk.Action{
		enums.ZH: m.SwitchToZH,
		enums.EN: m.SwitchToEn,
	}
}