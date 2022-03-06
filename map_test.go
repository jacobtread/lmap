package main

import (
	"fmt"
	"reflect"
	"testing"
)

func PutJunkData(m *LockingMap[string, int], count int) {
	for i := 0; i < count; i++ {
		m.Put(fmt.Sprintf("Test%d", i), i)
	}
}

func ContentEquals[T string | int](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, aV := range a {
		found := false
		for _, bV := range b {
			if aV == bV {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func ContentEqualsEntry(a []Entry[string, int], b []Entry[string, int]) bool {
	if len(a) != len(b) {
		return false
	}
	for _, aV := range a {
		found := false
		for _, bV := range b {
			if aV.Key == bV.Key || aV.Value == bV.Value {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestLockingMap_Put(t *testing.T) {
	m := NewMap[string, int]()
	m.Put("Test", 1)
	if !m.Contains("Test") {
		t.Logf("Expected map to contain key 'Test'")
		t.FailNow()
	}
}

func TestLockingMap_Contains(t *testing.T) {
	m := NewMap[string, int]()
	m.Put("Test", 1)
	if !m.Contains("Test") {
		t.Logf("Expected map to contain key 'Test'")
		t.FailNow()
	}
}

func TestLockingMap_ForEach(t *testing.T) {
	m := NewMap[string, int]()
	expectedCount := 20
	PutJunkData(&m, expectedCount)
	i := 0
	m.ForEach(func(key string, value int) {
		i++
	})
	if i != expectedCount {
		t.Logf("Expected iteration of %d elements but only iterated %d times", expectedCount, i)
		t.FailNow()
	}
}

func TestLockingMap_ForEachSafe(t *testing.T) {
	m := NewMap[string, int]()
	expectedCount := 20
	PutJunkData(&m, expectedCount)
	i := 0
	m.ForEachSafe(func(key string, value int) {
		i++
	})
	if i != expectedCount {
		t.Logf("Expected iteration of %d elements but only iterated %d times", expectedCount, i)
		t.FailNow()
	}
}

func TestLockingMap_ForEachUntil(t *testing.T) {
	m := NewMap[string, int]()
	elementCount := 20
	expectedCount := elementCount / 2
	PutJunkData(&m, elementCount)
	i := 0
	m.ForEachUntil(func(key string, value int) bool {
		i++
		return i >= expectedCount
	})
	if i != expectedCount {
		t.Logf("Expected iteration of %d elements but only iterated %d times", expectedCount, i)
		t.FailNow()
	}
}

func TestLockingMap_GetKeys(t *testing.T) {
	expectedKeys := []string{"Test1", "Test2", "Test3"}
	expectedValues := []int{1, 2, 3}
	m := NewMapOfArrays(expectedKeys, expectedValues)
	keys := m.GetKeys()
	if !ContentEquals(expectedKeys, keys) {
		t.Log("Keys did not match", expectedKeys, keys)
		t.FailNow()
	}
}

func TestLockingMap_GetValues(t *testing.T) {
	expectedKeys := []string{"Test1", "Test2", "Test3"}
	expectedValues := []int{1, 2, 3}
	m := NewMapOfArrays(expectedKeys, expectedValues)
	values := m.GetValues()
	if !ContentEquals(expectedValues, values) {
		t.Log("Values did not match", expectedValues, values)
		t.FailNow()
	}
}

func TestLockingMap_Size(t *testing.T) {
	m := NewMap[string, int]()
	PutJunkData(&m, 10)
	size := m.Size()
	if size != 10 {
		t.Logf("Expected map size to be 10 but got %d instead", size)
		t.FailNow()
	}
}

func TestLockingMap_GetOrCompute(t *testing.T) {
	m := NewMap[string, int]()
	m.Put("Test", 1)
	value := m.GetOrCompute("Test", func() int {
		return 2
	})
	if value != 1 {
		t.Logf("Expected the value of 'Test' to be 1 but got %d instead", value)
		t.FailNow()
	}
	value = m.GetOrCompute("Test2", func() int {
		return 2
	})
	if value != 2 {
		t.Logf("Expected the value of 'Test2' to be 2 but got %d instead", value)
		t.FailNow()
	}
}
func TestNewMapOfArrays(t *testing.T) {
	expectedKeys := []string{"Test1", "Test2", "Test3"}
	expectedValues := []int{1, 2, 3}
	m := NewMapOfArrays(expectedKeys, expectedValues)
	keys := m.GetKeys()
	if !ContentEquals(expectedKeys, keys) {
		t.Log("Keys did not match", expectedKeys, keys)
		t.FailNow()
	}
	values := m.GetValues()
	if !ContentEquals(expectedValues, values) {
		t.Log("Values did not match", expectedValues, values)
		t.FailNow()
	}
}

func TestLockingMap_GetEntries(t *testing.T) {
	expectEntries := []Entry[string, int]{
		{Key: "Test1", Value: 1},
		{Key: "Test2", Value: 5},
		{Key: "Test3", Value: 9},
	}
	m := NewMapOf(expectEntries)
	entries := m.GetEntries()
	if !ContentEqualsEntry(expectEntries, entries) {
		t.Log("Entries did not match", expectEntries, entries)
		t.FailNow()
	}
}

func TestNewMapOf(t *testing.T) {
	TestLockingMap_GetEntries(t)
}

func TestLockingMap_RemoveAndGet(t *testing.T) {
	m := NewMap[string, int]()
	m.Put("Test1", 1)

	value := m.RemoveAndGet("Test1")
	if value != 1 {
		t.Logf("Expected the value of 'Test1' to be 1 got %d instead", value)
		t.FailNow()
	}
}

func TestLockingMap_PutAll(t *testing.T) {
	m := NewMap[string, int]()

	expectEntries := []Entry[string, int]{
		{Key: "Test1", Value: 1},
		{Key: "Test2", Value: 5},
		{Key: "Test3", Value: 9},
	}

	m.PutAll(expectEntries)

	entries := m.GetEntries()
	if !ContentEqualsEntry(expectEntries, entries) {
		t.Log("Entries did not match", expectEntries, entries)
		t.FailNow()
	}
}

func TestLockingMap_GetEntryPointers(t *testing.T) {
	m := NewMap[string, int]()
	PutJunkData(&m, 10)
}

func TestLockingMap_GetValuePointers(t *testing.T) {
	m := NewMap[string, int]()
	PutJunkData(&m, 10)
}

func TestLockingMap_GetPointer(t *testing.T) {
	s := struct {
		S string
	}{S: "Test"}

	m := NewMap[string, interface{}]()
	m.Put("Test1", s)

	value := m.GetPointer("Test1")
	if reflect.ValueOf(value).Elem().Interface() != s {
		t.Log("Pointer didn't match")
		t.FailNow()
	}

	value = m.GetPointer("Test2")
	if value != nil {
		t.Log("Expected 'Test2' to be nil but got", value)
		t.FailNow()
	}
}

func TestLockingMap_Remove(t *testing.T) {
	m := NewMap[string, int]()
	m.Put("Test1", 1)
	m.Remove("Test1")
	if m.Contains("Test1") {
		t.Log("Expected key 'Test1' to be removed but it still exists")
		t.FailNow()
	}
}

func TestNewMap(t *testing.T) {
	m := NewMap[string, int]()
	if m.Size() != 0 {
		t.FailNow()
	}
}

func TestLockingMap_RemoveUnless(t *testing.T) {
	m := NewMap[string, int]()
	PutJunkData(&m, 10)
	m.RemoveUnless(func(key string, value int) bool {
		return value < 3
	})
	if m.Size() != 3 {
		t.Logf("Expected 3 elements to remain after RemoveUnless() had %d remaining", m.Size())
		t.FailNow()
	}
}

func TestLockingMap_Get(t *testing.T) {
	m := NewMap[string, int]()
	m.Put("Test", 1)
	value, exists := m.Get("Test")
	if !exists {
		t.Log("Expected map to contain key 'Test'")
		t.FailNow()
	}
	if value != 1 {
		t.Logf("Expected key 'Test' to have the value of '1' but got '%d' instead", value)
	}
}

func TestLockingMap_GetOrDefault(t *testing.T) {
	m := NewMap[string, int]()
	m.Put("Test", 1)
	value := m.GetOrDefault("Test", 0)
	if value != 1 {
		t.Logf("Expected value of 'Test' to be '1' got '%d'", value)
		t.FailNow()
	}
	value = m.GetOrDefault("Test2", -1)
	if value != -1 {
		t.Logf("Expected value of 'Test2' to be '-1' got '%d'", value)
		t.FailNow()
	}
}

func TestLockingMap_RemoveIf(t *testing.T) {
	m := NewMap[string, int]()
	PutJunkData(&m, 10)
	m.RemoveIf(func(key string, value int) bool {
		return value < 3
	})
	if m.Size() != (10 - 3) {
		t.Logf("Expected 7 elements to remain after RemoveIf() had %d remaining", m.Size())
		t.FailNow()
	}
}

func TestLockingMap_Clear(t *testing.T) {
	m := NewMap[string, int]()
	PutJunkData(&m, 10)
	m.Clear()
	if m.Size() > 0 {
		t.Logf("Expected map to be empty after Clear() but had %d elements", m.Size())
		t.FailNow()
	}
}
