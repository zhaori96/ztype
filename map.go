package ztype

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
)

// JSON is a convenience alias for Map with string keys and any values,
// representing a JSON-like generic map.
//
// Example:
//
//	var data JSON = NewMap(map[string]any{"name": "Alice", "age": 30})
//	fmt.Println(data.String()) // Output: {"name":"Alice","age":30}
type JSON = Map[string, any]

// Map is a generic type that wraps a map with keys of type K and values of type V.
// It tracks validity (null state) and whether it has been unmarshaled from JSON.
//
// Example:
//
//	m := NewMap(map[string]int{"one": 1, "two": 2})
//	fmt.Println(m.IsNull()) // Output: false
//	fmt.Println(m.Len())    // Output: 2
type Map[K comparable, V any] struct {
	value       map[K]V
	valid       bool
	unmarshaled bool
}

// NewMap creates a new Map with the given map value and marks it as valid.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1, "b": 2})
func NewMap[K comparable, V any](value map[K]V) Map[K, V] {
	return Map[K, V]{value: value, valid: true}
}

// NewNullMap creates a new Map that is marked as null (invalid).
//
// Example:
//
//	m := NewNullMap[string, int]()
func NewNullMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{valid: false}
}

// NewNullMapIfZero creates a new Map that is null if the input map is empty,
// otherwise returns a valid Map.
//
// Example:
//
//	m := NewNullMapIfZero(map[string]int{}) // null Map
//	m2 := NewNullMapIfZero(map[string]int{"a": 1}) // valid Map
func NewNullMapIfZero[K comparable, V any](value map[K]V) Map[K, V] {
	if len(value) == 0 {
		return NewNullMap[K, V]()
	}
	return NewMap(value)
}

// Get returns the underlying map value.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	v := m.Get() // map[string]int{"a": 1}
func (m Map[K, V]) Get() map[K]V {
	return m.value
}

// Set sets the internal map value and marks the Map as valid.
//
// Example:
//
//	var m Map[string]int
//	m.Set(map[string]int{"a": 1})
func (m *Map[K, V]) Set(value map[K]V) {
	m.value = value
	m.valid = true
}

// GetItem returns the value associated with the given key, and a boolean indicating existence.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	val, ok := m.GetItem("a") // val=1, ok=true
func (m Map[K, V]) GetItem(key K) (V, bool) {
	item, ok := m.value[key]
	return item, ok
}

// SetItem sets the value for the given key and marks the Map as valid.
//
// Example:
//
//	m := NewMap(map[string]int{})
//	m.SetItem("a", 42)
func (m *Map[K, V]) SetItem(key K, value V) {
	m.value[key] = value
	m.valid = true
}

// SetItemIf sets the value for the given key only if the condition is true.
//
// Example:
//
//	m := NewMap(map[string]int{})
//	m.SetItemIf("a", 42, true)  // sets
//	m.SetItemIf("b", 13, false) // does nothing
func (m *Map[K, V]) SetItemIf(key K, value V, condition bool) {
	if condition {
		m.SetItem(key, value)
	}
}

// DeleteItem removes the item with the given key and returns its value and true,
// or zero value and false if key does not exist.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	val, ok := m.DeleteItem("a") // val=1, ok=true
func (m *Map[K, V]) DeleteItem(key K) (V, bool) {
	if item, ok := m.GetItem(key); ok {
		delete(m.value, key)
		return item, true
	}
	var zero V
	return zero, false
}

// SetNull marks the Map as null and clears its content.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	m.SetNull()
func (m *Map[K, V]) SetNull() {
	m.value = map[K]V{}
	m.valid = false
}

// IsNull returns true if the Map is null (invalid).
//
// Example:
//
//	m := NewNullMap[string, int]()
//	if m.IsNull() { /* true */ }
func (m Map[K, V]) IsNull() bool {
	return !m.valid
}

// IsZero returns true if the internal map is empty.
//
// Example:
//
//	m := NewMap(map[string]int{})
//	fmt.Println(m.IsZero()) // true
func (m Map[K, V]) IsZero() bool {
	return len(m.value) == 0
}

// Len returns the number of items in the internal map.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1, "b": 2})
//	fmt.Println(m.Len()) // 2
func (m Map[K, V]) Len() int {
	return len(m.value)
}

// Unmarshaled returns true if the Map has been unmarshaled from JSON.
//
// Example:
//
//	var m Map[string]int
//	json.Unmarshal([]byte(`{"a":1}`), &m)
//	fmt.Println(m.Unmarshaled()) // true
func (m Map[K, V]) Unmarshaled() bool {
	return m.unmarshaled
}

// SetUnmarshaled sets the unmarshaled flag.
//
// Example:
//
//	var m Map[string]int
//	m.SetUnmarshaled(true)
func (m *Map[K, V]) SetUnmarshaled(value bool) {
	m.unmarshaled = value
}

// Has returns true if the key exists in the Map and the Map is valid.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	fmt.Println(m.Has("a")) // true
func (m Map[K, V]) Has(key K) bool {
	if !m.valid {
		return false
	}
	_, ok := m.value[key]
	return ok
}

// All returns a sequence of all key-value pairs.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1, "b": 2})
//	for pair := range m.All() { /* iterate pairs */ }
func (m Map[K, V]) All() iter.Seq2[K, V] {
	return maps.All(m.value)
}

// Insert adds all items from the given sequence to the Map and marks it valid.
//
// Example:
//
//	m := NewMap(map[string]int{})
//	m.Insert(iter.Of2([][2]interface{}{{"a", 1}, {"b", 2}}))
func (m *Map[K, V]) Insert(items iter.Seq2[K, V]) {
	maps.Insert(m.value, items)
	m.valid = true
}

// Keys returns a sequence of all keys.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	for key := range m.Keys() { fmt.Println(key) }
func (m Map[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.value)
}

// Values returns a sequence of all values.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	for value := range m.Values() { fmt.Println(value) }
func (m Map[K, V]) Values() iter.Seq[V] {
	return maps.Values(m.value)
}

// Collect creates a Map from the given sequence and marks it valid.
//
// Example:
//
//	var m Map[string]int
//	m.Collect(iter.Of2([][2]interface{}{{"a", 1}, {"b", 2}}))
func (m *Map[K, V]) Collect(items iter.Seq2[K, V]) {
	collected := maps.Collect(items)
	m.value = collected
	m.valid = true
}

// Filter returns a new Map containing only items where filter(key, value) is true.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1, "b": 2})
//	filtered := m.Filter(func(k string, v int) bool { return v > 1 })
func (m Map[K, V]) Filter(filter func(K, V) bool) Map[K, V] {
	result := map[K]V{}
	for key, value := range m.value {
		if filter(key, value) {
			result[key] = value
		}
	}
	m.value = result
	return m
}

// Merge merges other Maps into this Map, returning a new merged Map.
//
// Example:
//
//	m1 := NewMap(map[string]int{"a": 1})
//	m2 := NewMap(map[string]int{"b": 2})
//	merged := m1.Merge(m2)
func (m Map[K, V]) Merge(others ...Map[K, V]) Map[K, V] {
	merged := maps.Clone(m.value)
	for _, other := range others {
		maps.Copy(merged, other.value)
	}
	m.value = merged
	m.valid = true
	return m
}

// MergeRaw merges raw maps into this Map and returns a raw map.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	merged := m.MergeRaw(map[string]int{"b": 2})
func (m Map[K, V]) MergeRaw(others ...map[K]V) map[K]V {
	merged := maps.Clone(m.value)
	for _, other := range others {
		maps.Copy(merged, other)
	}
	return merged
}

// Clone returns a deep copy of the Map.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	c := m.Clone()
func (m Map[K, V]) Clone() Map[K, V] {
	m.value = maps.Clone(m.value)
	return m
}

// CloneRaw returns a deep copy of the underlying map.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	raw := m.CloneRaw()
func (m Map[K, V]) CloneRaw() map[K]V {
	return maps.Clone(m.value)
}

// EqualFunc returns true if this Map equals another Map using the provided equality function.
//
// Example:
//
//	m1 := NewMap(map[string]int{"a": 1})
//	m2 := NewMap(map[string]int{"a": 1})
//	equal := m1.EqualFunc(m2, func(a, b int) bool { return a == b })
func (m Map[K, V]) EqualFunc(other Map[K, V], equal func(V, V) bool) bool {
	return maps.EqualFunc(m.value, other.value, equal)
}

// EqualRawFunc returns true if this Map equals a raw map using the provided equality function.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	raw := map[string]int{"a": 1}
//	equal := m.EqualRawFunc(raw, func(a, b int) bool { return a == b })
func (m Map[K, V]) EqualRawFunc(other map[K]V, equal func(V, V) bool) bool {
	return maps.EqualFunc(m.value, other, equal)
}

// DeleteFunc deletes all items from the Map where the delete function returns true.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1, "b": 2})
//	m.DeleteFunc(func(k string, v int) bool { return v > 1 }) // removes "b"
func (m *Map[K, V]) DeleteFunc(delete func(K, V) bool) {
	maps.DeleteFunc(m.value, delete)
}

// JsonString returns a JSON string representation of the Map or "{}" if invalid.
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	s := m.JsonString() // "{\"a\":1}"
func (m Map[K, V]) JsonString() string {
	if !m.valid {
		return "{}"
	}
	data, erro := json.Marshal(m.value)
	if erro != nil {
		return ""
	}
	return string(data)
}

// MarshalJSON implements the json.Marshaler interface.
//
// Example:
//
//	json.Marshal(m)
func (n Map[K, V]) MarshalJSON() ([]byte, error) {
	if n.valid {
		return json.Marshal(n.value)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//
// Example:
//
//	json.Unmarshal(data, &m)
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	m.unmarshaled = true
	if bytes.Equal(data, []byte("null")) {
		m.valid = false
		m.value = map[K]V{}
		return nil
	}

	var result map[K]V
	if err := json.Unmarshal(data, &result); err != nil {
		m.valid = false
		return err
	}

	m.valid = true
	m.value = result
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
//
// Example:
//
//	m.MarshalText()
func (m Map[K, V]) MarshalText() ([]byte, error) {
	if m.valid {
		return json.Marshal(m.value)
	}
	return []byte("null"), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
//
// Example:
//
//	m.UnmarshalText(data)
func (m *Map[K, V]) UnmarshalText(data []byte) error {
	return m.UnmarshalJSON(data)
}

// Scan implements the sql.Scanner interface for database deserialization.
//
// Example:
//
//	var m Map[string]int
//	db.QueryRow(...).Scan(&m)
func (m *Map[K, V]) Scan(value any) error {
	if value == nil {
		m.valid = false
		m.value = map[K]V{}
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("invalid type: %T", value)
	}

	result := map[K]V{}
	if erro := json.Unmarshal(data, &result); erro != nil {
		m.valid = false
		return erro
	}

	m.valid = true
	m.value = result
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
//
// Example:
//
//	val, err := m.Value()
func (m Map[K, V]) Value() (driver.Value, error) {
	if !m.valid {
		return nil, nil
	}
	value, erro := json.Marshal(m.value)
	if erro != nil {
		return nil, erro
	}
	return string(value), nil
}

// String returns the JSON string representation of the Map.
// If the Map is invalid (null), it returns "{}".
//
// Example:
//
//	m := NewMap(map[string]int{"a": 1})
//	fmt.Println(m.String()) // Output: {"a":1}
func (m Map[K, V]) String() string {
	if !m.valid {
		return "null"
	}
	return fmt.Sprintf("%v", m.value)
}

// ComparableJSON is a convenience alias for MapComparable with string keys and any values,
// representing a JSON-like generic map with comparable values.
//
// Example:
//
//	var data ComparableJSON = MapComparable[string, any]{}
//	data.Set(map[string]any{"name": "Alice", "age": 30})
//	fmt.Println(data.String()) // Output: {"name":"Alice","age":30}
type ComparableJSON = MapComparable[string, any]

// MapComparable embeds Map[K, V] and adds methods
// useful when keys and values are comparable.
//
// Example:
//
//	m1 := MapComparable[string, int]{}
//	m1.Set(map[string]int{"a": 1, "b": 2})
//
//	m2 := MapComparable[string, int]{}
//	m2.Set(map[string]int{"a": 1, "b": 2})
//
//	equal := m1.Equal(m2) // true
type MapComparable[K comparable, V comparable] struct {
	Map[K, V]
}

// Equal returns true if m and other have exactly the same keys and values.
//
// Example:
//
//	equal := m1.Equal(m2)
func (m MapComparable[K, V]) Equal(other MapComparable[K, V]) bool {
	return maps.Equal(m.value, other.value)
}

// EqualRaw returns true if m and the raw map other have exactly the same keys and values.
//
// Example:
//
//	rawMap := map[string]int{"a": 1, "b": 2}
//	equal := m1.EqualRaw(rawMap)
func (m MapComparable[K, V]) EqualRaw(other map[K]V) bool {
	return maps.Equal(m.value, other)
}

// CompareAndSwap sets the value for key to new only if the current value is equal to old.
// Returns true if the swap was performed.
//
// Example:
//
//	swapped := m.CompareAndSwap("a", 1, 3) // true if current value is 1
func (m *MapComparable[K, V]) CompareAndSwap(key K, old, new V) bool {
	item, ok := m.GetItem(key)
	if !ok || item != old {
		return false
	}
	m.SetItem(key, new)
	return true
}

// DeleteIfEquals deletes the key only if its current value equals value.
// Returns true if the key was deleted.
//
// Example:
//
//	deleted := m.DeleteIfEquals("a", 3) // true if current value is 3
func (m *MapComparable[K, V]) DeleteIfEquals(key K, value V) bool {
	item, ok := m.GetItem(key)
	if ok && item == value {
		m.DeleteItem(key)
		return true
	}
	return false
}
