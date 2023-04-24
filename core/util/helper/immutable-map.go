package helper

import "github.com/huandu/go-clone"

type immutableMap[K comparable, V any] struct {
	m map[K]V
}

func NewImmutableMap[K comparable, V any](m map[K]V) immutableMap[K, V] {
	return immutableMap[K, V]{
		m: clone.Slowly(m).(map[K]V),
	}
}

func (i immutableMap[K, V]) Get(key K) (V, bool) {
	v, ok := i.m[key]
	if !ok {
		return v, ok
	}
	return clone.Slowly(v).(V), ok
}

func (i immutableMap[K, V]) BindGet(key K) V {
	return clone.Slowly(i.m[key]).(V)
}

func (i immutableMap[K, V]) GetMap() map[K]V {
	return clone.Slowly(i.m).(map[K]V)
}
