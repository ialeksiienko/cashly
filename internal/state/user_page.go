package state

import (
	"cashly/internal/entity"
	"sync"
)

type UserPage struct {
	Page     int
	Families []entity.Family
}

var (
	upage = make(map[int64]*UserPage)
	upm   sync.RWMutex
)

func SetUserPageState(uid int64, s *UserPage) {
	upm.Lock()
	defer upm.Unlock()

	upage[uid] = s
}

func GetUserPageState(uid int64) (*UserPage, bool) {
	upm.RLock()
	defer upm.RUnlock()

	s, ok := upage[uid]
	return s, ok
}
