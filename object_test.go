package alloc

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleObject(t *testing.T) {
	// create a new allocator
	arena := NewExpandingAllocator(pageSize)

	// create a new object in the allocator, we are only storing
	// a single value in here so lets keep it small
	obj := Must(NewObject[string, String, String](&arena, 1)).Deref()
	//     ^ This is just a helper function that panics if an error
	//     is returned
	//                                                      ^.Deref()
	//                                    is used to turn a Ptr[Object[...]]
	//                                    into just an Object[...]

	// set a value in the object
	obj.Set(
		// When calling Deref we actually point to the underlying bytes in
		// the allocator so there is no copying. AKA a pointer. So we need
		// to dereference the pointer hence the *
		*Must(NewString(&arena, "key")).Deref(),
		*Must(NewString(&arena, "value")).Deref(),
	)

	// Cool, setup done, now lets get the key
	v, ok := obj.Get("key")

	// woo hoo
	assert.Equal(t, "value", v.Cast())
	//                       ^ compare to primitive type, not the
	//                       underlying type

	// ok is the same as maps in golang, it will return false if we grab
	// something that doesn't exist.
	assert.Equal(t, true, ok)

	// here is the example of it not existing
	v, ok = obj.Get("other_key")
	assert.Equal(t, false, ok)

	// the value is then an empty value
	assert.Equal(t, "", v.Cast())
}

func TestObject(t *testing.T) {
	arena := NewExpandingAllocator(pageSize)
	obj, err := NewObject[string, String, int](&arena, 10)
	if !assert.NoError(t, err) {
		return
	}

	for x := range 10 {
		s, err := NewString(&arena, strconv.Itoa(x))
		if !assert.NoError(t, err) {
			return
		}

		obj.Deref().Set(*s.Deref(), x)
	}

	for x := range 10 {
		val, ok := obj.Deref().Get(strconv.Itoa(x))
		assert.Equal(t, x, val)
		assert.Equal(t, true, ok)
	}

	var x int
	for key, val := range obj.Deref().IterPrimitive() {
		assert.Equal(t, strconv.Itoa(x), key)
		assert.Equal(t, x, val)
		x++
	}

	x = 0
	for key := range obj.Deref().PrimitiveKeys() {
		assert.Equal(t, strconv.Itoa(x), key)
		x++
	}

	x = 0
	for val := range obj.Deref().Vals() {
		assert.Equal(t, x, val)
		x++
	}
}
