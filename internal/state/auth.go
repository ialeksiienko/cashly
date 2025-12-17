package state

import (
	"sync"
	"time"
)

var (
	lastAuthTime = make(map[int64]time.Time)
	authTimeout  = 5 * time.Minute
	mu           sync.RWMutex
)

func SetAuthorized(uid int64) {
	mu.Lock()
	defer mu.Unlock()

	lastAuthTime[uid] = time.Now()
}

func GetAuthorized(uid int64) (time.Time, bool) {
	mu.RLock()
	defer mu.RUnlock()

	t, ok := lastAuthTime[uid]
	return t, ok
}

func IsAuthorized(uid int64) bool {
	mu.RLock()
	defer mu.RUnlock()

	t, ok := lastAuthTime[uid]

	return ok && time.Since(t) < authTimeout
}
