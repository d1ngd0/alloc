package alloc

import (
	"iter"
	"unsafe"
)

// Array is a pointer to the underlying bytes, and it's defined length. This type is
// stored within the allocator. You can store this type for later use, as it doesn't
// reference raw pointers
type Array[T any] struct {
	data Ptr[T]
	len  int
}

// Slice returns the array as a slice value. Any changes to the values of the slice will
// be reflected in the slice, however if you append it will allocate a new slice and will
// no longer be in the allocator
func (s Array[T]) Slice() []T {
	return unsafe.Slice((*T)(unsafe.Pointer(s.data.Deref())), s.len)
}

// Length returns the length of the array
func (s Array[T]) Length() int {
	return s.len
}

// Expand creates a new Array with the new size specified, Copies the data
// into the new array, and returns it. The new locations will have uninitialized
// data in it.
func (s Array[T]) Expand(size int) (Array[T], error) {
	if s.len >= size {
		panic("new size must be larger than previous size")
	}

	b, err := NewArray[T](s.data.alloc, size)
	if err != nil {
		return Array[T]{}, nil
	}

	copy(b.Deref().Slice(), s.Slice())
	return *b.Deref(), nil
}

// Iter returns an iterator for the array
func (s Array[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, val := range s.Slice() {
			if !yield(val) {
				return
			}
		}
	}
}

// IterIndex creates an iterator which returns the index and value
func (s Array[T]) IterIndex() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for x, val := range s.Slice() {
			if !yield(x, val) {
				return
			}
		}
	}
}

// NewArray creates a new Array in the allocator and returns a pointer to the
// Array. Both the Underlying bytes, and the Array header are stored to the allocator
func NewArray[T any](a Allocator, len int) (Ptr[Array[T]], error) {
	// allocate the space for the raw bytes and the byte slice
	dataOffset, err := a.Alloc(unsafe.Sizeof(*new(T))*uintptr(len), unsafe.Alignof(*new(T)))
	if err != nil {
		return Ptr[Array[T]]{}, err
	}

	// we do this now to ensure all bytes are allocated before we start
	// referencing pointers. It is possible the data moves on each allocation
	sliceOffset, err := a.Alloc(unsafe.Sizeof(Array[T]{}), unsafe.Alignof(Array[T]{}))
	if err != nil {
		return Ptr[Array[T]]{}, err
	}

	s := (*Array[T])(unsafe.Pointer(a.Offset(sliceOffset)))
	s.data = Ptr[T]{
		offset: dataOffset,
		alloc:  a,
	}
	s.len = len

	return Ptr[Array[T]]{
		offset: sliceOffset,
		alloc:  a,
	}, nil
}
