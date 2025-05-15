package alloc

import "unsafe"

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

// Set is just a shorthand to deref the value and set
// the underlying bytes, it looks a little nicer than
// *(ptr.Deref()) = v
func (p Ptr[T]) Set(v T) {
	*(p.Deref()) = v
}

// bytes returns the raw underlying bytes, this is used
// for testing
func (p Ptr[T]) bytes() []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(p.alloc.Offset(p.offset))), unsafe.Sizeof(*new(T)))
}

// IsNull returns if the pointer is null
func (p Ptr[T]) IsNull() bool {
	return p.alloc == nil
}

// Null sets the pointer to null
func (p Ptr[T]) Null() {
	p.alloc = nil
}
