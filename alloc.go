package alloc

import (
	"errors"
	"unsafe"
)

var (
	ErrMemoryExhausted = errors.New("memory exhausted")
)

type Allocator interface {
	Alloc(s uintptr) (uintptr, error)
	Available() uintptr
}

func New[T any](a Allocator) (*T, error) {
	ptr, err := a.Alloc(unsafe.Sizeof(*new(T)))
	if err != nil {
		return nil, err
	}

	return (*T)(unsafe.Pointer(ptr)), nil
}
