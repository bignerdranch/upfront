package main

import (
	"sync"
)

// MutexDB provides an in-memory key/value store
type MutexDB[T any] struct {
	mu   sync.RWMutex
	Data map[string]T
}

// Get unmarshals a key into your desired type, where the last type returned
// is if the key was present.
func (db *MutexDB[T]) Get(key string) (T, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	val, ok := db.Data[key]
	return val, ok
}

// Set stores the value, marshaled to json, into the database, otherwise returns
// any error that may have occurred.
//
// Returns if the key already existed
func (db *MutexDB[T]) Set(key string, value T) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.Data[key]
	db.Data[key] = value

	return ok
}
