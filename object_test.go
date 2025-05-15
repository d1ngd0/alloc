package alloc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	arena := NewExpandingAllocator(pageSize)
	obj, err := NewObject[int, int](&arena, 10)
	if !assert.NoError(t, err) {
		return
	}

	for x := range 10 {
		obj.Deref().Set(x, x)
	}

	for x := range 10 {
		val, ok := obj.Deref().Get(x)
		assert.Equal(t, x, val)
		assert.Equal(t, true, ok)
	}

	var x int
	for key, val := range obj.Deref().Iter() {
		assert.Equal(t, x, key)
		assert.Equal(t, x, val)
		x++
	}

	x = 0
	for key := range obj.Deref().Keys() {
		assert.Equal(t, x, key)
		x++
	}

	x = 0
	for val := range obj.Deref().Vals() {
		assert.Equal(t, x, val)
		x++
	}
}
