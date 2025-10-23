package mysync

import (
	"sync"
	"sync/atomic"
)

type TypedMap[Key comparable, Value any] struct {
	m   *sync.Map
	cnt atomic.Int64
}

func NewTypedMap[Key comparable, Value any]() *TypedMap[Key, Value] {
	return &TypedMap[Key, Value]{
		m: &sync.Map{},
	}
}

func (tm *TypedMap[Key, Value]) Get(key Key) (val Value, ok bool) {
	v, ok := tm.m.Load(key)
	if !ok {
		return val, false
	}

	val, ok = v.(Value)
	if !ok {
		panic("impossible type of value")
	}

	return val, true
}

func (tm *TypedMap[Key, Value]) Set(key Key, val Value) {
	tm.m.Store(key, val)
	tm.cnt.Add(1)
}

func (tm *TypedMap[Key, Value]) Del(key Key) {
	tm.m.Delete(key)
	tm.cnt.Add(-1)
}

func (tm *TypedMap[Key, Value]) Len() int64 {
	return tm.cnt.Load()
}

func (tm *TypedMap[Key, Value]) Range(f func(key Key, val Value) bool) {
	tm.m.Range(func(k, v any) bool {
		key, ok := k.(Key)
		if !ok {
			panic("impossible type of key")
		}

		val, ok := v.(Value)
		if !ok {
			panic("impossible type of value")
		}

		return f(key, val)
	})
}

func (tm *TypedMap[Key, Value]) Clear() {
	tm.m.Clear()
}

func (tm *TypedMap[Key, Value]) Map() map[Key]Value {
	m := make(map[Key]Value)
	tm.m.Range(func(k, v any) bool {
		key, ok := k.(Key)
		if !ok {
			panic("impossible type of key")
		}
		value, ok := v.(Value)
		if !ok {
			panic("impossible type of value")
		}
		m[key] = value
		return true
	})
	return m
}
