package alloc

import "iter"

type Primitive[T any] interface {
	Primitive() T
}

// Object uses a linear search pattern to find the key specified.
// C is the underlying primitive golang type. When you call Get or
// iterate over the values, you likley want the golang primitive type, and
// not the type stored in alloc.
// K is the type for the key which is saved on the allocator. This type must
// implement Primitive to allow the key to be "cast" into a more familiar type
// T is the value type, and far less complicated :)
// example
//
//	func main() {
//	  arena :=
//	}
type Object[C comparable, K Primitive[C], T any] struct {
	keys Array[K]
	vals Array[T]
	len  int
}

// NewObject creates a new object on the heap. The Allocator passed in will store the
// underlying value, and the size is how many items you expect to store in the object.
// *the object can grow in size beyond the defined size* but it will cause a copy to a new
// location on the Allocator, and the previous bytes will not be cleaned up
func NewObject[C comparable, K Primitive[C], T any](alloc Allocator, size int) (Ptr[Object[C, K, T]], error) {
	obj, err := New[Object[C, K, T]](alloc)
	if err != nil {
		return Ptr[Object[C, K, T]]{}, err
	}

	keys, err := NewArray[K](alloc, size)
	if err != nil {
		return Ptr[Object[C, K, T]]{}, err
	}
	obj.Deref().keys = *keys.Deref()

	vals, err := NewArray[T](alloc, size)
	if err != nil {
		return Ptr[Object[C, K, T]]{}, err
	}
	obj.Deref().vals = *vals.Deref()

	// uninitialized data means there could be garbage in this value,
	// so we need to set the value to 0
	obj.Deref().len = 0

	return obj, nil
}

// index will return the index the key was found at. If the key was not
// found it will return -1
func (m Object[C, K, T]) index(key C) int {
	for x, val := range m.keys.Slice() {
		if x >= m.len {
			break
		}

		if val.Primitive() == key {
			return x
		}
	}

	return -1
}

// full returns if the object is full
func (m Object[C, K, T]) full() bool {
	return m.len >= m.keys.Length()
}

// grow increases the size of the object when full
func (m *Object[C, K, T]) grow() error {
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
func (m *Object[C, K, T]) Set(key K, val T) error {
	index := m.index(key.Primitive())
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
func (m Object[C, K, T]) Get(key C) (T, bool) {
	index := m.index(key)
	if index == -1 {
		return *new(T), false
	}

	return m.vals.Slice()[index], true
}

// Iter returns an iterator for the object enabling you to use this in a
// for each (range). The key will be the first value, and the type will be
// the second.
func (m Object[C, K, T]) Iter() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		for x, key := range m.keys.Slice() {
			val := m.vals.Slice()[x]
			if !yield(key, val) {
				return
			}
		}
	}
}

// Iter returns an iterator for the object enabling you to use this in a
// for each (range). The key will be the first value, and the type will be
// the second.
func (m Object[C, K, T]) IterPrimitive() iter.Seq2[C, T] {
	return func(yield func(C, T) bool) {
		for x, key := range m.keys.Slice() {
			val := m.vals.Slice()[x]
			if !yield(key.Primitive(), val) {
				return
			}
		}
	}
}

// Keys returns an interator of key values as primitives
func (m Object[C, K, T]) PrimitiveKeys() iter.Seq[C] {
	return func(yield func(C) bool) {
		for _, key := range m.keys.Slice() {
			if !yield(key.Primitive()) {
				return
			}
		}
	}
}

// Keys returns an interator of key values
func (m Object[C, K, T]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, key := range m.keys.Slice() {
			if !yield(key) {
				return
			}
		}
	}
}

// Values returns an interator of values
func (m Object[C, K, T]) Vals() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, val := range m.vals.Slice() {
			if !yield(val) {
				return
			}
		}
	}
}
