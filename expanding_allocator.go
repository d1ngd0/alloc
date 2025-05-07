package alloc

import (
	"math"
	"unsafe"
)

// ExpandingAllocator creates a byte slice which it expands over time to hold
// all the data. Ince the byte slice is filled, it creates a new byte slice with
// double the size, copies the data over, and keeps going. Any Deref Ptr values
// will no longer be valid when we move the underlying data, so it is important
// to call Deref only when you want or need the underlying value
type ExpandingAllocator struct {
	b *[]byte
}

// ensure we implement allocator
var _ Allocator = &ExpandingAllocator{}

// NewExpandingAllocator will create a new Expanding allocator
func NewExpandingAllocator(size int) ExpandingAllocator {
	if size < allocatorAlignment {
		panic("allocator must be equal to or larger than 8")
	}

	b := make([]byte, 0, size)
	return ExpandingAllocator{&b}
}

// Alloc reserves a section of memory and returns the offset to it. If we are going
// to exhaust the memory, we create a new location for the memory with twice the size,
// copy the data over and then allocate
func (a *ExpandingAllocator) Alloc(size uintptr, alignment uintptr) (uintptr, error) {
	// find the start by aligning
	start := align(uintptr(len(*a.b)), alignment)
	// find the end
	end := start + size

	// if we are not large enough to hold the data we need to grow the underlying
	// bytes and move our data over
	if uintptr(cap(*a.b)) < end {
		// create the new location
		b := make([]byte, end, end*2+allocatorAlignment)
		// ensure the byte slice is aligned to the largest possible alignment
		b = align_slice(b, allocatorAlignment)
		// move the data over
		copy(b, *a.b)
		// switch the bytes over
		*a.b = b
	} else {
		// underlying array is large enough, so just increase the size of the array
		*a.b = (*a.b)[:end]
	}

	// return the offset to the newly allocated item
	return uintptr(start), nil
}

// Offset returns the actual uintptr
func (a *ExpandingAllocator) Offset(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(&(*a.b)[offset])
}

// Available always return MaxUInt64 since we will "never" run out of
// memory
func (a *ExpandingAllocator) Available() uintptr {
	return math.MaxUint64
}

// Reset sets the head back to 0, Any allocations relying on these
// bytes will be overwritten over time, only call this function if you
// *know* that all references to this data are gone
func (a *ExpandingAllocator) Reset() {
	*a.b = (*a.b)[:0]
}
