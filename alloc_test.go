package alloc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlloc(t *testing.T) {
	arena := NewPageAllocator()
	i1, _ := New[uint64](&arena)
	*(i1.Deref()) = math.MaxUint64
	i2, _ := New[uint64](&arena)
	*(i2.Deref()) = math.MaxUint32

	assert.Equal(t, arena.b[0:16], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0})
}

func TestPtr(t *testing.T) {
	arena := NewExpandingAllocator(8)
	// This fits and will not cause expansion
	i, _ := New[uint64](&arena)
	*(i.Deref()) = 100
	// this should be pointed to one location
	p1 := arena.Offset(i.offset)

	// This will cause a new byte slice to be created, so the pointer
	// will change, but i1 is still valid
	_, _ = New[uint64](&arena)

	p2 := arena.Offset(i.offset)
	assert.NotEqual(t, p1, p2)
	assert.Equal(t, *(i.Deref()), uint64(100))
}

func BenchmarkAlloc(b *testing.B) {
	b.Run("control", func(b *testing.B) {
		b.ReportAllocs()
		allocInt := func() *[1024]byte {
			var v [1024]byte
			return &v
		}

		for b.Loop() {
			var _ = allocInt()
		}
	})

	b.Run("page_allocator", func(b *testing.B) {
		b.ReportAllocs()
		allocInt := func(alloc Allocator) Ptr[[1024]byte] {
			v, _ := New[[1024]byte](alloc)
			return v
		}

		alloc := NewPageAllocator()
		for b.Loop() {
			var _ = allocInt(&alloc)
			alloc.Reset()
		}
	})
}

func BenchmarkExpandingAlloc(b *testing.B) {
	b.Run("control_1000", func(b *testing.B) {
		b.ReportAllocs()
		allocInt := func() *[1024]byte {
			var v [1024]byte
			return &v
		}

		for b.Loop() {
			for range 1000 {
				var _ = allocInt()
			}
		}
	})

	b.Run("expanding_allocator_1000", func(b *testing.B) {
		b.ReportAllocs()
		allocInt := func(alloc Allocator) Ptr[[1024]byte] {
			v, _ := New[[1024]byte](alloc)
			return v
		}

		alloc := NewExpandingAllocator(4096)
		for b.Loop() {
			for range 1000 {
				var _ = allocInt(&alloc)
			}
			alloc.Reset()
		}
	})
}
