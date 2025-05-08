package alloc

import "unsafe"

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
