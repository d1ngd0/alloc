package alloc

import "unsafe"

const pageSize = 4096

// PageAllocator is an allocator with only 4096 bytes. This is the
// size of a page in linux. It will return ErrMemoryExhausted when full
type PageAllocator struct {
	ref uintptr
	b   [pageSize]byte
}

// ensure we implement the allocator
var _ Allocator = &PageAllocator{}

// NewPageAllocator will create a new page allocator
func NewPageAllocator() PageAllocator {
	return PageAllocator{}
}

// Alloc reserves the location in memory and returns the offset the
// new allocation occured at. If the page can not fit the size required
// ErrMemoryExhausted is returned.
func (a *PageAllocator) Alloc(size uintptr, alignment uintptr) (uintptr, error) {
	start := align(a.ref, alignment)
	end := start + size

	if pageSize < int(end) {
		return 0, ErrMemoryExhausted
	}

	a.ref = end
	return start, nil
}

// offset returns the pointer to the offset supplied
func (a *PageAllocator) Offset(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(&(a.b)[offset])
}

// Available returns the amount of memory left in the page which can
// be allocated to.
func (a *PageAllocator) Available() uintptr {
	return pageSize - a.ref
}

// Reset sets the head back to 0, Any allocations relying on these
// bytes will be overwritten over time, only call this function if you
// *know* that all references to this data are gone
func (a *PageAllocator) Reset() {
	a.ref = 0
}
