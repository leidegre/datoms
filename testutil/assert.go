package testutil

import (
	"testing"
	"unsafe"
)

func AreEqual[T comparable](t *testing.T, expected, actual T) {
	if !(expected == actual) {
		t.Fatalf("expected %v actual %v", expected, actual)
	}
}

func sliceRange[T any](s []T) (start, end uintptr) {
	start = uintptr(unsafe.Pointer(unsafe.SliceData(s)))
	end = start + unsafe.Sizeof(s[0])*uintptr(cap(s))
	return
}

func isOverlapping[T comparable](a, b []T) bool {
	var (
		x, y = sliceRange(a)
		z, w = sliceRange(b)
	)

	return (x < w) && (z < y)
}

func AreDistinctSlice[T comparable](t *testing.T, expected, actual []T) {
	if isOverlapping(expected, actual) {
		t.Fatal("slices are not distinct (references same memory)")
	}
}

func AreEqualSlice[T comparable](t *testing.T, expected, actual []T) {
	if !(len(expected) == len(actual)) {
		t.Fatalf("expected %v actual %v", len(expected), len(actual))
	}
	for i, expected := range expected {
		if !(expected == actual[i]) {
			t.Fatalf("[%v]: expected %v actual %v", i, expected, actual[i])
		}
	}
}
