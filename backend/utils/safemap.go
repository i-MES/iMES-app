package utils

import "sync"

/*
线程安全的 Map
*/
type RWMap struct {
	Lock sync.RWMutex
	DMap map[string]string
}

func (r *RWMap) Get(key string) string {
	r.Lock.RLock()
	defer r.Lock.RUnlock()

	if v, ok := r.DMap[key]; ok {
		return v
	}
	return ""
}

func (r *RWMap) Put(k, v string) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	if r.DMap == nil {
		r.DMap = make(map[string]string)
	}
	r.DMap[k] = v
}
