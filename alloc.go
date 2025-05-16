package alloc

import (
	"errors"
	"unsafe"
)

var (
	ErrMemoryExhausted = errors.New("memory exhausted")
	ErrOutOfRange      = errors.New("offset out of range")
)

// Allocators are used to create an allocation of the
type Allocator interface {
	// Alloc creates a new item in memory with a size defined by the parameter
	// it returns the offset within allocated memory to the location. If any
	// errors occured they will be returned
	Alloc(size uintptr, alignment uintptr) (offset uintptr, err error)

	// Offset takes the parameter offset, and returns the actual pointer to the
	// location.
	Offset(offset uintptr) (ptr unsafe.Pointer)

	// Available returns the memory remaining until the allocator is exhausted
	Available() uintptr
}

// New will create a new type in the allocator, and return a pointer
// to that type
func New[T any](a Allocator) (Ptr[T], error) {
	offset, err := a.Alloc(
		unsafe.Sizeof(*new(T)),
		unsafe.Alignof(*new(T)),
	)

	return Ptr[T]{offset: offset, alloc: a}, err
}

// Must wraps any allocatation functions and panics if an error occurs
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err.Error())
	}

	return v
}

// Primitive changes the type into it's underlying go primitive type, for
// instance, the special "String" type returns a string
type Primitive[T any] interface {
	Cast() T
}

// allocatorAlignment makes sure the byte slice is aligned to the larges possible size
// which is 8. Then when we copy things over everything stays aligned
const allocatorAlignment = 8

// align updates the underlying byte array so that it aligns to the largest size
func align_slice(b []byte, alignment uintptr) []byte {
	// grab the location of the byte data
	ptr := uintptr(unsafe.Pointer(&b[0]))
	// calculate the location for the aligned value
	alignedPtr := align(ptr, alignment)
	// update the byte slice to point to the aligned value
	return b[alignedPtr-ptr:]
}

// align will take a uintptr and a number and turn it into an aligned starting point
func align(index, alignment uintptr) uintptr {
	if index%alignment == 0 {
		return index
	}

	return index + alignment - (index % alignment)
}
