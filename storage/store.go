package storage

import "sync"

var (
	storage = make(map[string]string)
	mutex   sync.RWMutex
)

func Set(key, value string) {
	mutex.Lock()
	defer mutex.Unlock()
	storage[key] = value
}

func Get(key string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	value, exits := storage[key]
	return value, exits
}

func Delete(key string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, exists := storage[key]
	if exists {
		delete(storage, key)
	}
	return exists
}
