package genericsync

import "sync"

// Map is a utility type wrapping a sync.Map using generics
type Map[K comparable, V any] struct {
	internal sync.Map
}

// Load tries to retrieve a value from m
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	untypedValue, ok := m.internal.Load(key)
	if !ok {
		return value, false
	}
	value = untypedValue.(V)
	return value, true
}

// Store stores a value in m
func (m *Map[K, V]) Store(key K, value V) {
	m.internal.Store(key, value)
}

// Range calls f sequentially for each key and value present in m
func (m *Map[K, V]) Range(f func(K, V)) {
	m.internal.Range(func(key, value any) bool {
		f(key.(K), value.(V))
		return true
	})
}
