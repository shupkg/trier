package strs

import "unsafe"

type (
	stringHeader struct {
		Data unsafe.Pointer
		Len  int
	}
	sliceHeader struct {
		Data unsafe.Pointer
		Len  int
		Cap  int
	}
)

// UnsafeString returns an unsafe string reference of b.
// The caller must treat the input slice as immutable.
//
// WARNING: Use carefully. The returned result must not leak to the end user
// unless the input slice is provably immutable.
func UnsafeString(b []byte) (s string) {
	src := (*sliceHeader)(unsafe.Pointer(&b))
	dst := (*stringHeader)(unsafe.Pointer(&s))
	dst.Data = src.Data
	dst.Len = src.Len
	return s
}

// UnsafeBytes returns an unsafe bytes slice reference of s.
// The caller must treat returned slice as immutable.
//
// WARNING: Use carefully. The returned result must not leak to the end user.
func UnsafeBytes(s string) (b []byte) {
	src := (*stringHeader)(unsafe.Pointer(&s))
	dst := (*sliceHeader)(unsafe.Pointer(&b))
	dst.Data = src.Data
	dst.Len = src.Len
	dst.Cap = src.Len
	return b
}
