package alloc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlloc(t *testing.T) {
	arena := NewPageAllocator()
	i1, _ := New[uint64](&arena)
	*i1 = math.MaxUint64
	i2, _ := New[uint64](&arena)
	*i2 = math.MaxUint32

	assert.Equal(t, arena.b[0:16], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0})
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
		allocInt := func(alloc Allocator) *[1024]byte {
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
		allocInt := func(alloc Allocator) *[1024]byte {
			v, _ := New[[1024]byte](alloc)
			return v
		}

		alloc := NewExpandingAllocator()
		for b.Loop() {
			for range 1000 {
				var _ = allocInt(&alloc)
			}
			alloc.Reset()
		}
	})
}
