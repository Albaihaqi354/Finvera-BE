package blacklist

import (
	"sync"
	"time"
)

var (
	tokens = make(map[string]time.Time)
	mutex  sync.RWMutex
)

// Add adds a token to the blacklist with its expiration time
func Add(token string, exp time.Time) {
	mutex.Lock()
	defer mutex.Unlock()
	tokens[token] = exp
}

// IsBlacklisted checks if a token is in the blacklist and not expired
func IsBlacklisted(token string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	
	exp, exists := tokens[token]
	if !exists {
		return false
	}
	
	if time.Now().After(exp) {
		// Clean up expired token asynchronously (or let a cron do it)
		// For simplicity, we just return false
		return false
	}
	
	return true
}

// Cleanup removes expired tokens from the memory
func Cleanup() {
	mutex.Lock()
	defer mutex.Unlock()
	
	now := time.Now()
	for token, exp := range tokens {
		if now.After(exp) {
			delete(tokens, token)
		}
	}
}
