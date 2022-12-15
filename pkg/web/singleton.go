package web

import (
	"reflect"
	"sync"
)

var cache sync.Map

func GetSingleton[T any]() (t *T) {
	hash := reflect.TypeOf(t)
	instance, hasInstance := cache.Load(hash)
	if !hasInstance {
		instance = new(T)
		instance, _ = cache.LoadOrStore(hash, instance)
	}
	return instance.(*T)
}
