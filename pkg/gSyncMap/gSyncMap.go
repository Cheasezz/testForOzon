package gSyncMap

//go:generate mockgen -package gSyncMap -source=gSyncMap.go -destination=mocks_test.go

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

type IGSyncMap[K comparable, V any] interface {
	Load(key K) (V, error)
	Delete(key K)
	Range(f func(key K, value V) bool)
	LoadOrStore(key K, value V) (V, bool)
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

func (gm *GSyncMap[K, V]) Range(f func(key K, value V) bool) {
	gm.m.Range(func(key, value interface{}) bool {
		k, ok1 := key.(K)
		v, ok2 := value.(V)
		if !ok1 || !ok2 {
			return false
		}
		return f(k, v)
	})
}

func (gm *GSyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := gm.m.LoadOrStore(key, value)

	return actual.(V), loaded
}
