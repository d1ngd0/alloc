package alloc

import (
	"errors"
	"unsafe"
)

var (
	ErrMemoryExhausted = errors.New("memory exhausted")
	ErrOutOfRange      = errors.New("offset out of range")
)

// Ptr returns a pointer to the underlying value. This pointer
// tracks the allocator used to provision it. When using an allocator
// which might move the underlying data, this abstraction makes sure
// you can always retrieve the data. You should hold onto and pass this
// around instead of passing around the value from Defer since that could
// change after subsiquent allocations
type Ptr[T any] struct {
	offset uintptr
	alloc  Allocator
}

// Defer return the underlying type as a pointer
func (p Ptr[T]) Deref() *T {
	ptr := p.alloc.Offset(p.offset)
	return (*T)(ptr)
}

// Allocators are used to create an allocation of the
type Allocator interface {
	// Alloc creates a new item in memory with a size defined by the parameter
	// it returns the offset within allocated memory to the location. If any
	// errors occured they will be returned
	Alloc(size uintptr) (offset uintptr, err error)

	// Offset takes the parameter offset, and returns the actual pointer to the
	// location.
	Offset(offset uintptr) (ptr unsafe.Pointer)

	// Available returns the memory remaining until the allocator is exhausted
	Available() uintptr
}

// New will create a new type in the allocator, and return a pointer
// to that type
func New[T any](a Allocator) (Ptr[T], error) {
	offset, err := a.Alloc(unsafe.Sizeof(*new(T)))
	return Ptr[T]{offset: offset, alloc: a}, err

}
