package alloc

import (
	"math"
	"unsafe"
)

type ExpandingAllocator struct {
	b *[]byte
}

func NewExpandingAllocator() ExpandingAllocator {
	b := make([]byte, 0, pageSize)
	return ExpandingAllocator{&b}
}

func (a *ExpandingAllocator) Alloc(size uintptr) (uintptr, error) {
	start := uintptr(len(*a.b))
	end := start + size

	if uintptr(cap(*a.b)) < end {
		b := make([]byte, end, end*2)
		copy(b, *a.b)
		*a.b = b
	} else {
		*a.b = (*a.b)[:end]
	}

	return uintptr(unsafe.Pointer(&(*a.b)[start])), nil
}

func (a *ExpandingAllocator) Available() uintptr {
	return math.MaxUint64
}

// Reset sets the head back to 0, Any allocations relying on these
// bytes will be overwritten over time, only call this function if you
// *know* that all references to this data are gone
func (a *ExpandingAllocator) Reset() {
	*a.b = (*a.b)[:0]
}
