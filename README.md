# Go Locking Map

![Go](https://img.shields.io/badge/Powered%20By-Go-29BEB0?style=for-the-badge)
![LINES OF CODE](https://img.shields.io/tokei/lines/github/jacobtread/go-locking-map?style=for-the-badge)
![LICENSE](https://img.shields.io/github/license/jacobtread/go-locking-map?style=for-the-badge)

A simple map struct which uses Go 1.18 generics to provide a map type that uses locks and has lots of helper functions
to do specific operations on the underlying map

Install with

```shell
$ go get github.com/jacobtread/go-locking-map
```

## Creating Maps

### Empty Map

The following code creates a new empty map

```go
m := NewMap[string, int]()
```

### Map from entries

The following code shows the creation of a map from existing entries

```go
entries := []Entry[string, int]{
{Key: "Test1", Value: 1},
{Key: "Test2", Value: 5},
{Key: "Test3", Value: 9},
}
m := NewMapOf(entries)
```

### Map from arrays

The following code shows the creation of a map from two arrays. One for keys and one for values

```go
keys := []string{"Test1", "Test2", "Test3"}
values := []int{1, 2, 3}
m := NewMapOfArrays(keys, values)
```

## Inserting Values

### Inserting one

To insert one key value pair you can use the "Put" function which takes in the key and value

```go
key := "Test"
value := 1
m.Put(key, value)
```

### Inserting many

You can insert lots of key value pairs by using an array of entries.

```go
entries := []Entry[string, int]{
{Key: "Test1", Value: 1},
{Key: "Test2", Value: 5},
{Key: "Test3", Value: 9},
}
m.PutAll(entries)
```

## Retrieving Values

### Retrieve value with key

You can retrieve a single value using its key with the `Get` function. This will return both the value and whether the
value exists in the underlying map

```go
key := "Test1"
value, exists := m.Get(key)
```

### Retrieving value with key and default

If you want to return another value if the key doesn't exist in the map you can use the `GetOrDefault` function which
takes a default value as the second argument

```go
key := "Test1"
value := m.GetOrDefault(key, 1)
```

### Retrieving value pointer

If you want to retrieve a pointer to the value instead of retrieving the actual value you can use `GetPointer` which
will return a pointer to the value or nil if the value doesn't exist

```go
key := "Test1"
value := m.GetPointer(key)
```

### Get with compute fallback

If you want to generate a new value if the key isn't present in the map and both insert and return that value you can
use `GetOrCompute` which takes a compute function as the second argument which will be called and set as the value if
the key doesn't exist

```go
key := "Test1"
value := m.GetOrCompute(key, func () {
return 1 // Test1 will be assigned to 1 if it's not already set
})
```

### Checking if a key exists

If you want to check if a key exists in the map without retrieving its value you can call the `Contains` function which
returns whether the key exists

```go
key := "Test1"
exists := m.Contains(key)
```

### Retrieving total number of entries

To retrieve the total number of entries contained within the map you can use the
`Size` function. This returns the number of entries

```go
size := m.Size()
```

### Retrieving all the keys

If you would like an array of all the keys in the map you can use the `GetKeys` function which will create an array of
all the keys

```go
keys := m.GetKeys()
```

### Retrieving all the values

If you would like an array of all the values in the map you can use the `GetValues` function which will create an array
of all the values

```go
values := m.GetValues()
```

### Retrieving pointers to all the values

If you would like an array of pointers to all the values in the map you can use the
`GetValuePointers` function which will create an array of pointers to each of the values in the map

```go
pointers := m.GetValuePointers()
```

### Retrieving pointers to all the values as entries

If you would like an array of pointers to all the values in the form of entries you can use the
`GetEntryPointers` function this will create an array of entries but the entries will have pointers to the values
instead of the values themselves

```go
entries := m.GetEntryPointers()
```

### Retrieving all the key value pairs

If you would like an array of all the keys and values in the map you can use the
`GetEntries` function which will create an array of entries for all the map key and value pairs

```go
entries := m.GetEntries()
```

### Checking for matches
If you would like to see if any entries in the map match a specific condition or multiple you can use the `AnyMatch` function
which will run the provided test function on each all the entries in the map. The function will return true as soon as one of
the tests pass
```go
matches := m.AnyMatch(func (key string, value int) bool {
	return value > 5
})
```

### Checking everything matches
If you want to make sure that every entry in the map matches a certain condition you can use the `AllMatch` function which will 
return false if any of the entry tests don't return true
```go
matches := m.AllMatch(func (key string, value int) bool {
	return value != 0
})
```

## Iteration

### For loop iteration

To iterate over the map like a for loop you can use the `ForEach` function which takes in an action function. This
action function will be called for every element in the loop. Note: You must not modify the underlying elements of the
map using this function or the lock will fail if you want to do that you should use the `ForEachSafe` function instead
which is identical in functionality, but it copies the entries before iteration to prevent concurrent read writes

```go
m.ForEach(func (key string, value int) {
fmt.Println(key, value)
})

m.ForEachSafe(func (key string, value int) {
fmt.Println(key, value)
})
```

### For loop with break

To iterate over the for loop with a breaking condition you can use the `ForEachUntil` function this is the same
as `ForEach` except the action function returns a boolean value. When the action function returns true the for loop
iteration is stopped

```go
m.ForEachUntil(func (key string, value int) {
fmt.Println(key, value)
return value < 0
})
```

### Finding Min Value

If you have data within your map which can be used in the form of an integer, and you wish to find the minimum value you
can use the `MinOf` function which iterates over the map key value pairs and maps them to the provided int value using
the action function provided.

```go
min := m.MinOf(func (key string, value int) int {
  return value
})
```

### Finding Max Value

If you have data within your map which can be used in the form of an integer, and you wish to find the maximum value you
can use the `MaxOf` function which iterates over the map key value pairs and maps them to the provided int value using
the action function provided.

```go
max := m.MaxOf(func (key string, value int) int {
  return value
})
```

### Finding the Sum

If you have data within your map which can be used in the form of an integer, and you wish to find the sum of all the
values you can use the `SumOf` function which iterates over the map key value pairs and counts the value mapped to the
provided int value using the action function provided.

```go
sum := m.SumOf(func (key string, value int) int {
  return value
})
```

## Removing entries

### Removing a key

You can remove a key with the `Remove` function

```go
key := "Test1"
m.Remove(key)
```

### Removing a key and getting its value

If you want to remove a key but get its value as well you can use the
`RemoveAndGet` function which returns the value of the key

```go
key := "Test1"
value := m.RemoveAndGet(key)
```

### Conditional Removal

If you would like to remove all values that match a custom predicate you can use the
`RemoveIf` function this takes an "action" function as the argument. This function is called on each entry in the map
and will remove any entries that return true. The inverse of this is `RemoveUnless` which removes all elements that
don't return true

```go
m.RemoveIf(func (key string, value int) bool {
  return value != 5
})

m.RemoveUnless(func (key string, value int) bool {
  return value == 5
})
```

### Removing Everything

To clear the map completely you can use the `Clear` function which deletes all the keys and values in the map

```go
m.Clear()
```

### Removing Everything with action

If you want to completely clear all the contents of the map, but you also want to run logic for each
entry that gets removed you can use the `ClearAnd` function which removes all the map keys and runs
the provided action function for each removed entry
```go
m.ClearAnd(func(key string, value int) {
	// Do something with the entry...
})
```
