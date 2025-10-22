package session

import (
	"cashly/internal/entity"
	"sync"
)

type UserPageState struct {
	Page     int
	Families []entity.Family
}

var (
	userPageState      = make(map[int64]*UserPageState)
	userPageStateMutex sync.RWMutex
)

func SetUserPageState(userID int64, state *UserPageState) {
	userPageStateMutex.Lock()
	defer userPageStateMutex.Unlock()

	userPageState[userID] = state
}

func GetUserPageState(userID int64) (*UserPageState, bool) {
	userPageStateMutex.RLock()
	defer userPageStateMutex.RUnlock()

	state, exists := userPageState[userID]
	if !exists {
		return nil, false
	}
	return state, true
}
