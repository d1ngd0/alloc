package alloc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	arena := NewExpandingAllocator(4096)

	s, _ := NewArray[int](&arena, 10)
	b := s.Deref().Slice()
	for x := range b {
		b[x] = x
	}

	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, s.Deref().Slice())
}
