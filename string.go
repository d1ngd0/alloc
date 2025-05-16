package alloc

import (
	"cmp"
	"unsafe"
)

type String Array[byte]

// NewString returns a String
func NewString(alloc Allocator, s string) (Ptr[String], error) {
	return NewStringFromBytes(alloc, []byte(s))
}

// NewStringFromBytes returns a Ptr to a string
func NewStringFromBytes(alloc Allocator, b []byte) (Ptr[String], error) {
	arr, err := NewArray[byte](alloc, len(b))
	if err != nil {
		return Ptr[String]{}, nil
	}

	copy(arr.Deref().Slice(), b)
	return Ptr[String](arr), nil
}

// String reutrns the underlying value as a safe golang string
func (s String) String() string {
	return string(Array[byte](s).Slice())
}

// Cast returns the underlying values as an unsafe golang
// string. There is no copying of bytes in the string, but it can change
// if the underlying bytes change. You should only use this if you
// **know** the value will not chnage
func (s String) Cast() string {
	if Array[byte](s).Length() == 0 {
		return ""
	}

	return unsafe.String(
		&(Array[byte](s).Slice()[0]),
		Array[byte](s).Length(),
	)
}

// Cmp implements the Comparable interface which is used for object keys
func (s String) Cmp(val String) int {
	return cmp.Compare(s.Cast(), val.Cast())
}
