package state

import (
	"sync"
)

type UserTextType int

const (
	None UserTextType = iota
	WaitingFamilyName
	WaitingFamilyCode
	WaitingBankToken
)

var (
	ut  = make(map[int64]UserTextType)
	utm sync.RWMutex
)

func SetTextState(uid int64, s UserTextType) {
	utm.Lock()
	defer utm.Unlock()

	ut[uid] = s
}

func GetTextState(uid int64) UserTextType {
	utm.RLock()
	defer utm.RUnlock()

	s, ok := ut[uid]
	if !ok {
		return None
	}
	return s
}

func ClearTextState(uid int64) {
	utm.Lock()
	defer utm.Unlock()

	delete(ut, uid)
}
