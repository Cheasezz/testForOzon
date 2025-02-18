package gSyncMap

import (
	"errors"
	"sync"
)

var (
	errLoad        = errors.New("error GSyncMap.load")
	errLoadOrStore = errors.New("error GenerGSyncMapicMap.LoadOrStore")
	errTypeAssert  = errors.New("invalid type assertion")
)

type GSyncMap[K comparable, V any] struct {
	m sync.Map
}

func NewGenericSyncMap[K comparable, V any]() *GSyncMap[K, V] {
	return new(GSyncMap[K, V])
}

func (gm *GSyncMap[K, V]) Store(key K, value V) {

	gm.m.Store(key, value)
}

func (gm *GSyncMap[K, V]) Load(key K) (V, error) {
	value, ok := gm.m.Load(key)
	if !ok {
		var zeroValue V
		return zeroValue, errLoad
	}

	return value.(V), nil
}

func (gm *GSyncMap[K, V]) Delete(key K) {
	gm.m.Delete(key)
}

func (gm *GSyncMap[K, V]) Range(f func(key any, value any) bool) {
	gm.m.Range(f)
}

func (gm *GSyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := gm.m.LoadOrStore(key, value)

	return actual.(V), loaded
}
