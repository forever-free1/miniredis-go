package server

import (
	"sync"
	"time"
)

// Data types for Redis values
type Value struct {
	Data      string
	ExpireAt  *time.Time // nil means no expiration
	Type      string     // "string", "list", "hash", "set"
}

// ListValue represents a Redis list
type ListValue struct {
	Data     []string
	ExpireAt *time.Time
}

// HashValue represents a Redis hash
type HashValue struct {
	Data     map[string]string
	ExpireAt *time.Time
}

// SetValue represents a Redis set
type SetValue struct {
	Data     map[string]struct{}
	ExpireAt *time.Time
}

// Store is the in-memory data store
var (
	store    = make(map[string]*Value)
	listStore    = make(map[string]*ListValue)
	hashStore    = make(map[string]*HashValue)
	setStore    = make(map[string]*SetValue)
	mu       sync.RWMutex
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

// ==================== List Commands ====================

// LPush pushes values to the left of a list
func LPush(key string, values ...string) int {
	mu.Lock()
	defer mu.Unlock()

	list, ok := listStore[key]
	if !ok {
		list = &ListValue{Data: make([]string, 0)}
		listStore[key] = list
	}

	// Insert at the beginning
	newData := make([]string, 0, len(values)+len(list.Data))
	newData = append(newData, values...)
	newData = append(newData, list.Data...)
	list.Data = newData

	return len(list.Data)
}

// RPush pushes values to the right of a list
func RPush(key string, values ...string) int {
	mu.Lock()
	defer mu.Unlock()

	list, ok := listStore[key]
	if !ok {
		list = &ListValue{Data: make([]string, 0)}
		listStore[key] = list
	}

	list.Data = append(list.Data, values...)
	return len(list.Data)
}

// LRange returns elements from start to stop
func LRange(key string, start, stop int) ([]string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	list, ok := listStore[key]
	if !ok {
		return nil, false
	}

	data := list.Data
	n := len(data)

	// Handle negative indices
	if start < 0 {
		start = n + start
	}
	if stop < 0 {
		stop = n + stop
	}

	// Handle out of bounds
	if start < 0 {
		start = 0
	}
	if stop >= n {
		stop = n - 1
	}
	if start > stop || start >= n {
		return []string{}, true
	}

	return data[start : stop+1], true
}

// LLen returns the length of a list
func LLen(key string) int {
	mu.RLock()
	defer mu.RUnlock()

	list, ok := listStore[key]
	if !ok {
		return 0
	}

	return len(list.Data)
}

// LIndex returns element at index
func LIndex(key string, index int) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	list, ok := listStore[key]
	if !ok {
		return "", false
	}

	n := len(list.Data)
	if index < 0 {
		index = n + index
	}
	if index < 0 || index >= n {
		return "", false
	}

	return list.Data[index], true
}

// ==================== Hash Commands ====================

// HSet sets field in hash
func HSet(key, field, value string) int {
	mu.Lock()
	defer mu.Unlock()

	hash, ok := hashStore[key]
	if !ok {
		hash = &HashValue{Data: make(map[string]string)}
		hashStore[key] = hash
	}

	_, isNew := hash.Data[field]
	hash.Data[field] = value
	if isNew {
		return 0
	}
	return 1
}

// HGet gets field from hash
func HGet(key, field string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	hash, ok := hashStore[key]
	if !ok {
		return "", false
	}

	value, ok := hash.Data[field]
	return value, ok
}

// HGetAll gets all fields and values from hash
func HGetAll(key string) (map[string]string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	hash, ok := hashStore[key]
	if !ok {
		return nil, false
	}

	result := make(map[string]string)
	for k, v := range hash.Data {
		result[k] = v
	}
	return result, true
}

// HDel deletes fields from hash
func HDel(key string, fields ...string) int {
	mu.Lock()
	defer mu.Unlock()

	hash, ok := hashStore[key]
	if !ok {
		return 0
	}

	count := 0
	for _, field := range fields {
		if _, ok := hash.Data[field]; ok {
			delete(hash.Data, field)
			count++
		}
	}
	return count
}

// HExists checks if field exists in hash
func HExists(key, field string) bool {
	mu.RLock()
	defer mu.RUnlock()

	hash, ok := hashStore[key]
	if !ok {
		return false
	}

	_, ok = hash.Data[field]
	return ok
}

// HLen returns number of fields in hash
func HLen(key string) int {
	mu.RLock()
	defer mu.RUnlock()

	hash, ok := hashStore[key]
	if !ok {
		return 0
	}

	return len(hash.Data)
}

// ==================== Set Commands ====================

// SAdd adds members to a set
func SAdd(key string, members ...string) int {
	mu.Lock()
	defer mu.Unlock()

	set, ok := setStore[key]
	if !ok {
		set = &SetValue{Data: make(map[string]struct{})}
		setStore[key] = set
	}

	count := 0
	for _, member := range members {
		if _, exists := set.Data[member]; !exists {
			set.Data[member] = struct{}{}
			count++
		}
	}
	return count
}

// SMembers returns all members of a set
func SMembers(key string) ([]string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	set, ok := setStore[key]
	if !ok {
		return nil, false
	}

	members := make([]string, 0, len(set.Data))
	for member := range set.Data {
		members = append(members, member)
	}
	return members, true
}

// SIsMember checks if member exists in set
func SIsMember(key, member string) bool {
	mu.RLock()
	defer mu.RUnlock()

	set, ok := setStore[key]
	if !ok {
		return false
	}

	_, ok = set.Data[member]
	return ok
}

// SCard returns cardinality of set
func SCard(key string) int {
	mu.RLock()
	defer mu.RUnlock()

	set, ok := setStore[key]
	if !ok {
		return 0
	}

	return len(set.Data)
}

// SRem removes members from a set
func SRem(key string, members ...string) int {
	mu.Lock()
	defer mu.Unlock()

	set, ok := setStore[key]
	if !ok {
		return 0
	}

	count := 0
	for _, member := range members {
		if _, exists := set.Data[member]; exists {
			delete(set.Data, member)
			count++
		}
	}
	return count
}

// ==================== Pub/Sub ====================

// Subscriber represents a client subscription
type Subscriber struct {
	Channels map[string]bool
	Patterns map[string]bool
}

// PubSub manages publish/subscribe state
type PubSub struct {
	mu         sync.RWMutex
	subscribers map[string]map[*Subscriber]bool
	patterns    map[string]map[*Subscriber]bool
}

var pubsub *PubSub

func init() {
	pubsub = &PubSub{
		subscribers: make(map[string]map[*Subscriber]bool),
		patterns:    make(map[string]map[*Subscriber]bool),
	}
}

// Publish publishes a message to a channel
func Publish(channel, message string) int {
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()

	count := 0

	// Send to exact channel subscribers
	if subs, ok := pubsub.subscribers[channel]; ok {
		for range subs {
			count++
		}
	}

	// Send to pattern subscribers
	for pattern, subs := range pubsub.patterns {
		if matchPattern(pattern, channel) {
			for range subs {
				count++
			}
		}
	}

	return count
}

// Subscribe subscribes a client to channels
func Subscribe(sub *Subscriber, channels ...string) {
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()

	for _, ch := range channels {
		if sub.Channels == nil {
			sub.Channels = make(map[string]bool)
		}
		sub.Channels[ch] = true

		if pubsub.subscribers[ch] == nil {
			pubsub.subscribers[ch] = make(map[*Subscriber]bool)
		}
		pubsub.subscribers[ch][sub] = true
	}
}

// Unsubscribe unsubscribes a client from channels
func Unsubscribe(sub *Subscriber, channels ...string) {
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()

	if len(channels) == 0 {
		// Unsubscribe from all
		for ch := range sub.Channels {
			if pubsub.subscribers[ch] != nil {
				delete(pubsub.subscribers[ch], sub)
			}
		}
		sub.Channels = make(map[string]bool)
		for pat := range sub.Patterns {
			if pubsub.patterns[pat] != nil {
				delete(pubsub.patterns[pat], sub)
			}
		}
		sub.Patterns = make(map[string]bool)
		return
	}

	for _, ch := range channels {
		if sub.Channels != nil {
			delete(sub.Channels, ch)
		}
		if pubsub.subscribers[ch] != nil {
			delete(pubsub.subscribers[ch], sub)
		}
	}
}

// PSubscribe subscribes a client to patterns
func PSubscribe(sub *Subscriber, patterns ...string) {
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()

	for _, pat := range patterns {
		if sub.Patterns == nil {
			sub.Patterns = make(map[string]bool)
		}
		sub.Patterns[pat] = true

		if pubsub.patterns[pat] == nil {
			pubsub.patterns[pat] = make(map[*Subscriber]bool)
		}
		pubsub.patterns[pat][sub] = true
	}
}

// PUnsubscribe unsubscribes a client from patterns
func PUnsubscribe(sub *Subscriber, patterns ...string) {
	pubsub.mu.Lock()
	defer pubsub.mu.Unlock()

	if len(patterns) == 0 {
		// Unsubscribe from all patterns
		for pat := range sub.Patterns {
			if pubsub.patterns[pat] != nil {
				delete(pubsub.patterns[pat], sub)
			}
		}
		sub.Patterns = make(map[string]bool)
		return
	}

	for _, pat := range patterns {
		if sub.Patterns != nil {
			delete(sub.Patterns, pat)
		}
		if pubsub.patterns[pat] != nil {
			delete(pubsub.patterns[pat], sub)
		}
	}
}

// matchPattern checks if a channel matches a pattern
func matchPattern(pattern, channel string) bool {
	// Simple pattern matching: * matches any characters
	i := 0
	j := 0
	for i < len(pattern) && j < len(channel) {
		if pattern[i] == '*' {
			// Star matches any sequence
			i++
			if i >= len(pattern) {
				return true
			}
			// Find next character after star
			for j < len(channel) && channel[j] != pattern[i] {
				j++
			}
		} else if pattern[i] == '?' {
			// Question matches any single character
			i++
			j++
		} else if pattern[i] == channel[j] {
			i++
			j++
		} else {
			return false
		}
	}
	return i >= len(pattern) && j >= len(channel)
}
