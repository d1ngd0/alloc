package alloc

import "unsafe"

const pageSize = 4096

type PageAllocator struct {
	ref uintptr
	b   [pageSize]byte
}

func NewPageAllocator() PageAllocator {
	return PageAllocator{}
}

func (a *PageAllocator) Alloc(size uintptr) (uintptr, error) {
	start := a.ref
	end := start + size

	if pageSize < int(end) {
		return 0, ErrMemoryExhausted
	}

	a.ref = end
	return uintptr(unsafe.Pointer(&a.b)) + start, nil
}

func (a *PageAllocator) Available() uintptr {
	return pageSize - a.ref
}

// Reset sets the head back to 0, Any allocations relying on these
// bytes will be overwritten over time, only call this function if you
// *know* that all references to this data are gone
func (a *PageAllocator) Reset() {
	a.ref = 0
}
