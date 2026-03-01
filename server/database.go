package server

import (
	"sync"
	"time"
)

// Data types for Redis values
type Value struct {
	Data      string
	ExpireAt  *time.Time // nil means no expiration
}

// Store is the in-memory data store
var (
	store = make(map[string]*Value)
	mu    sync.RWMutex
)

// Set sets a key-value pair
func Set(key, value string) {
	mu.Lock()
	defer mu.Unlock()
	store[key] = &Value{Data: value}
}

// Get gets a value by key
func Get(key string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := store[key]
	if !ok {
		return "", false
	}

	// Check expiration
	if v.ExpireAt != nil && time.Now().After(*v.ExpireAt) {
		return "", false
	}

	return v.Data, true
}

// Delete deletes a key
func Delete(key string) bool {
	mu.Lock()
	defer mu.Unlock()

	_, ok := store[key]
	if ok {
		delete(store, key)
		return true
	}
	return false
}

// Exists checks if a key exists
func Exists(key string) bool {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := store[key]
	if !ok {
		return false
	}

	// Check expiration
	if v.ExpireAt != nil && time.Now().After(*v.ExpireAt) {
		return false
	}

	return true
}

// SetWithExpire sets a key-value pair with expiration
func SetWithExpire(key, value string, expire time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	expireAt := time.Now().Add(expire)
	store[key] = &Value{
		Data:     value,
		ExpireAt: &expireAt,
	}
}

// Expire sets expiration for a key
func Expire(key string, seconds int) bool {
	mu.Lock()
	defer mu.Unlock()

	v, ok := store[key]
	if !ok {
		return false
	}

	// Key exists but already has expiration in the past
	if v.ExpireAt != nil && time.Now().After(*v.ExpireAt) {
		return false
	}

	expireAt := time.Now().Add(time.Duration(seconds) * time.Second)
	v.ExpireAt = &expireAt
	return true
}

// TTL returns the remaining time to live for a key
func TTL(key string) (int, bool) {
	mu.RLock()
	defer mu.RUnlock()

	v, ok := store[key]
	if !ok {
		return -2, false
	}

	if v.ExpireAt == nil {
		return -1, true
	}

	remaining := time.Until(*v.ExpireAt)
	if remaining <= 0 {
		return -2, false
	}

	return int(remaining.Seconds()), true
}

// Persist removes expiration from a key
func Persist(key string) bool {
	mu.Lock()
	defer mu.Unlock()

	v, ok := store[key]
	if !ok {
		return false
	}

	v.ExpireAt = nil
	return true
}
