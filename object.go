package alloc

import "iter"

// Object uses a linear
type Object[K comparable, T any] struct {
	keys Array[K]
	vals Array[T]
	len  int
}

// NewObject creates a new object on the heap. The Allocator passed in will store the
// underlying value, and the size is how many items you expect to store in the object.
// *the object can grow in size beyond the defined size* but it will cause a copy to a new
// location on the Allocator, and the previous bytes will not be cleaned up
func NewObject[K comparable, T any](alloc Allocator, size int) (Ptr[Object[K, T]], error) {
	obj, err := New[Object[K, T]](alloc)
	if err != nil {
		return Ptr[Object[K, T]]{}, err
	}

	keys, err := NewArray[K](alloc, size)
	if err != nil {
		return Ptr[Object[K, T]]{}, err
	}
	obj.Deref().keys = *keys.Deref()

	vals, err := NewArray[T](alloc, size)
	if err != nil {
		return Ptr[Object[K, T]]{}, err
	}
	obj.Deref().vals = *vals.Deref()

	// uninitialized data means there could be garbage in this value,
	// so we need to set the value to 0
	obj.Deref().len = 0

	return obj, nil
}

// index will return the index the key was found at. If the key was not
// found it will return -1
func (m Object[K, T]) index(key K) int {
	for x, val := range m.keys.Slice() {
		if x >= m.len {
			break
		}

		if val == key {
			return x
		}
	}

	return -1
}

// full returns if the object is full
func (m Object[K, T]) full() bool {
	return m.len >= m.keys.Length()
}

// grow increases the size of the object when full
func (m *Object[K, T]) grow() error {
	newlen := m.len * 2
	if newlen == 0 {
		newlen = 10
	}

	keys, err := m.keys.Expand(newlen)
	if err != nil {
		return err
	}

	vals, err := m.vals.Expand(newlen)
	if err != nil {
		return err
	}

	m.len = newlen
	m.keys = keys
	m.vals = vals

	return nil
}

// Set stores a value in the object. It will check to make sure there is enough space
// in the object and re-allocate the map if needed to make space by doubling the size
// of the object. Allocation errors can be returned when the underlying object is grown
// if you do not cause an expansion there will be no errors.
func (m *Object[K, T]) Set(key K, val T) error {
	index := m.index(key)
	if index != -1 {
		m.vals.Slice()[index] = val
		return nil
	}

	// make sure we can fit this new value into the object
	if m.full() {
		err := m.grow()
		if err != nil {
			return err
		}
	}

	m.keys.Slice()[m.len] = key
	m.vals.Slice()[m.len] = val
	m.len++
	return nil
}

// Get returns the value from the map, if no value exists we return the empty
// value of T and false
func (m Object[K, T]) Get(key K) (T, bool) {
	index := m.index(key)
	if index == -1 {
		return *new(T), false
	}

	return m.vals.Slice()[index], true
}

// Iter returns an iterator for the object enabling you to use this in a
// for each (range). The key will be the first value, and the type will be
// the second.
func (m Object[K, T]) Iter() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		for x, key := range m.keys.Slice() {
			val := m.vals.Slice()[x]
			if !yield(key, val) {
				return
			}
		}
	}
}

// Keys returns an interator of key values
func (m Object[K, T]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, key := range m.keys.Slice() {
			if !yield(key) {
				return
			}
		}
	}
}

// Values returns an interator of values
func (m Object[K, T]) Vals() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, val := range m.vals.Slice() {
			if !yield(val) {
				return
			}
		}
	}
}
