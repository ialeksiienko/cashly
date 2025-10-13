package session

import (
	"sync"
)

type UserTextStateType int

const (
	StateNone UserTextStateType = iota
	StateWaitingFamilyName
	StateWaitingFamilyCode
	StateWaitingBankToken
	StateWaitingPassword
)

var (
	userTextState      = make(map[int64]UserTextStateType)
	userTextStateMutex sync.RWMutex
)

func SetTextState(userID int64, state UserTextStateType) {
	userTextStateMutex.Lock()
	defer userTextStateMutex.Unlock()

	userTextState[userID] = state
}

func GetTextState(userID int64) UserTextStateType {
	userTextStateMutex.RLock()
	defer userTextStateMutex.RUnlock()

	state, exists := userTextState[userID]
	if !exists {
		return StateNone
	}
	return state
}

func ClearTextState(userID int64) {
	userTextStateMutex.Lock()
	defer userTextStateMutex.Unlock()

	delete(userTextState, userID)
}
