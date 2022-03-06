package main

import (
	"sync"
)

type LockingMap[K string | uint | int, V any] struct {
	Lock       *sync.RWMutex
	Underlying map[K]V
}

type Entry[K string | uint | int, V any] struct {
	Key   K
	Value V
}

// NewMap constructor for creating new locking maps
func NewMap[K string | uint | int, V any]() LockingMap[K, V] {
	m := make(map[K]V)
	out := LockingMap[K, V]{
		Lock:       &sync.RWMutex{},
		Underlying: m,
	}
	return out
}

// NewMapOf constructor for creating new locking maps with the provided
// default key value pairs
func NewMapOf[K string | uint | int, V any](entries []Entry[K, V]) LockingMap[K, V] {
	m := make(map[K]V)
	out := LockingMap[K, V]{
		Lock:       &sync.RWMutex{},
		Underlying: m,
	}
	for _, entry := range entries {
		m[entry.Key] = entry.Value
	}
	return out
}

// NewMapOfArrays constructor for creating new locking maps from the two arrays provided
// the first array will be used as keys and the second will be used as values. The
// arrays must both have a matching length
func NewMapOfArrays[K string | uint | int, V any](keys []K, values []V) LockingMap[K, V] {
	if len(keys) != len(values) {
		panic("Expected keys and values to have same length")
	}
	m := make(map[K]V)
	out := LockingMap[K, V]{
		Lock:       &sync.RWMutex{},
		Underlying: m,
	}
	for i, key := range keys {
		value := values[i]
		m[key] = value
	}
	return out
}

// ForEach iterates over all the key values in the underlying map
// and runs the action function for each of them. Note: DO NOT MODIFY
// THE UNDERLYING MAP WITHIN THIS FUNCTION USE ForEachSafe if you are
// going to access any write methods
func (m *LockingMap[K, V]) ForEach(action func(key K, value V)) {
	m.Lock.RLock()
	for k, v := range m.Underlying {
		action(k, v)
	}
	m.Lock.RUnlock()
}

// ForEachSafe iterates over all the key values in the underlying map
// and runs the action function for each of them. This function is
// safe for map write operations because it copies the entries before
// iterating
func (m *LockingMap[K, V]) ForEachSafe(action func(key K, value V)) {
	entries := m.GetEntries()
	for _, entry := range entries {
		action(entry.Key, entry.Value)
	}
}

// ForEachUntil iterates over all the key values in the underlying map
// and runs the action function for each of them until the action function
// returns true
func (m *LockingMap[K, V]) ForEachUntil(action func(key K, value V) bool) {
	m.Lock.RLock()
	for k, v := range m.Underlying {
		if action(k, v) {
			break
		}
	}
	m.Lock.RUnlock()
}

func (m *LockingMap[K, V]) Size() int {
	m.Lock.RLock()
	//goland:noinspection GoVarAndConstTypeMayBeOmitted
	var u map[K]V = m.Underlying // IDE gets confused without this
	l := len(u)
	m.Lock.RUnlock()
	return l
}

// PutAll inserts all the provided entries into the map
func (m *LockingMap[K, V]) PutAll(entries []Entry[K, V]) {
	m.Lock.Lock()
	for _, entry := range entries {
		m.Underlying[entry.Key] = entry.Value
	}
	m.Lock.Unlock()
}

// Put inserts the key value pair into the map.
func (m *LockingMap[K, V]) Put(key K, value V) {
	m.Lock.Lock()
	m.Underlying[key] = value
	m.Lock.Unlock()
}

// Contains Returns whether the map contains the provided key
func (m *LockingMap[K, V]) Contains(key K) bool {
	m.Lock.RLock()
	_, exists := m.Underlying[key]
	m.Lock.RUnlock()
	return exists
}

// Get retrieves the value from the map using its key. The second
// return value indicates whether a value is present
func (m *LockingMap[K, V]) Get(key K) (V, bool) {
	m.Lock.RLock()
	value, exists := m.Underlying[key]
	m.Lock.RUnlock()
	return value, exists
}

// GetPointer retrieves the value from the map using its key. Will
// return a pointer to the value or nil if it doesn't exist
func (m *LockingMap[K, V]) GetPointer(key K) *V {
	m.Lock.RLock()
	value, exists := m.Underlying[key]
	m.Lock.RUnlock()
	if exists {
		return &value
	} else {
		return nil
	}
}

// GetOrDefault retrieves the value from the map using its key. Returns
// the value provided as d if the key doesn't exist
func (m *LockingMap[K, V]) GetOrDefault(key K, d V) V {
	m.Lock.RLock()
	value, exists := m.Underlying[key]
	m.Lock.RUnlock()
	if !exists {
		return d
	} else {
		return value
	}
}

// GetOrCompute retrieves the value from the map using its key. If
// the provided key doesn't exist then the compute function will be
// called and that will be inserted into the map
func (m *LockingMap[K, V]) GetOrCompute(key K, compute func() V) V {
	m.Lock.RLock()
	value, exists := m.Underlying[key]
	m.Lock.RUnlock()
	if !exists {
		computed := compute()
		m.Put(key, computed)
		return computed
	} else {
		return value
	}
}

// Remove safely removes the key from the underlying map.
func (m *LockingMap[K, V]) Remove(key K) {
	m.Lock.Lock()
	delete(m.Underlying, key)
	m.Lock.Unlock()
}

// RemoveAndGet safely removes the key from the underlying map. And
// returns the value that existed or nil
func (m *LockingMap[K, V]) RemoveAndGet(key K) V {
	value, _ := m.Get(key)
	m.Lock.Lock()
	delete(m.Underlying, key)
	m.Lock.Unlock()
	return value
}

// RemoveIf Runs the provided action on all the entries in the map
// any calls that return true will be deleted from the underlying map
func (m *LockingMap[K, V]) RemoveIf(action func(key K, value V) bool) {
	values := m.GetEntries()
	m.Lock.Lock()
	for _, entry := range values {
		if action(entry.Key, entry.Value) {
			delete(m.Underlying, entry.Key)
		}
	}
	m.Lock.Unlock()
}

// RemoveUnless Runs the provided action on all the entries in the map
// any calls that return false will be deleted from the underlying map
func (m *LockingMap[K, V]) RemoveUnless(action func(key K, value V) bool) {
	values := m.GetEntries()
	m.Lock.Lock()
	for _, entry := range values {
		if !action(entry.Key, entry.Value) {
			delete(m.Underlying, entry.Key)
		}
	}
	m.Lock.Unlock()
}

// Clear removes all keys and values from the underlying map.
func (m *LockingMap[K, V]) Clear() {
	keys := m.GetKeys()
	m.Lock.Lock()
	for _, k := range keys {
		delete(m.Underlying, k)
	}
	m.Lock.Unlock()
}

// GetValuePointers creates an array with pointers to all the values stored
// inside the locking map.
func (m *LockingMap[K, V]) GetValuePointers() []*V {
	m.Lock.RLock()
	out := make([]*V, m.Size())
	i := 0
	for _, v := range m.Underlying {
		out[i] = &v
		i++
	}
	m.Lock.RUnlock()
	return out
}

// GetValues creates an array with all the values stored inside the
// locking map.
func (m *LockingMap[K, V]) GetValues() []V {
	m.Lock.RLock()
	out := make([]V, m.Size())
	i := 0
	for _, v := range m.Underlying {
		out[i] = v
		i++
	}
	m.Lock.RUnlock()
	return out
}

// GetKeys creates an array with all the keys stored inside the
// locking map.
func (m *LockingMap[K, V]) GetKeys() []K {
	m.Lock.RLock()
	out := make([]K, m.Size())
	i := 0
	for k, _ := range m.Underlying {
		out[i] = k
		i++
	}
	m.Lock.RUnlock()
	return out
}

// GetEntries creates an array with all the key and values stored inside the
// locking map.
func (m *LockingMap[K, V]) GetEntries() []Entry[K, V] {
	m.Lock.RLock()
	out := make([]Entry[K, V], m.Size())
	i := 0
	for k, v := range m.Underlying {
		out[i] = Entry[K, V]{Key: k, Value: v}
		i++
	}
	m.Lock.RUnlock()
	return out
}

// GetEntryPointers creates an array with all the key and value pointers stored inside the
// locking map.
func (m *LockingMap[K, V]) GetEntryPointers() []Entry[K, *V] {
	m.Lock.RLock()
	out := make([]Entry[K, *V], m.Size())
	i := 0
	for k, v := range m.Underlying {
		out[i] = Entry[K, *V]{Key: k, Value: &v}
		i++
	}
	m.Lock.RUnlock()
	return out
}
