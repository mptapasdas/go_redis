package storage

import "sync"

var store sync.Map

func Set(key, value string) {
	store.Store(key, value)
}

func Get(key string) (string, bool) {
	value, exists := store.Load(key)
	if !exists {
		return "", false
	}
	return value.(string), true
}

func Delete(key string) bool {
	_, exists := store.Load(key)
	if exists {
		store.Delete(key)
	}
	return exists
}
