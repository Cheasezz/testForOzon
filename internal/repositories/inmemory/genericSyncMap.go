package inmemory

import (
	"fmt"
	"sync"
)

var (
	errLoad = fmt.Errorf("error GenericMap.load")
)

type GenericMap[K comparable, V any] struct {
	m sync.Map
}

func (gm *GenericMap[K, V]) Store(key K, value V) {

	gm.m.Store(key, value)
}

func (gm *GenericMap[K, V]) Load(key K) (*V, error) {
	value, ok := gm.m.Load(key)
	if !ok {
		return nil, errLoad
	}
	return value.(*V), nil
}

func (gm *GenericMap[K, V]) Delete(key K) {
	gm.m.Delete(key)
}
