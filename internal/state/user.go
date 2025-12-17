package state

import (
	"cashly/internal/entity"
	"sync"
)

type UserState struct {
	Family *entity.Family
}

var (
	us  = make(map[int64]*UserState)
	usm sync.RWMutex
)

func SetUserState(uid int64, s *UserState) {
	usm.Lock()
	defer usm.Unlock()

	us[uid] = s
}

func GetUserState(uid int64) (*UserState, bool) {
	usm.RLock()
	defer usm.RUnlock()

	s, ok := us[uid]
	return s, ok
}

func DeleteUserState(uid int64) {
	usm.Lock()
	defer usm.Unlock()

	delete(us, uid)
}
