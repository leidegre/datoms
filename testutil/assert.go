package testutil

import "testing"

func AreEqual[T comparable](t *testing.T, expected, actual T) {
	if !(expected == actual) {
		t.Fatalf("expected %v actual %v", expected, actual)
	}
}
