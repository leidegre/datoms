package testutil

import "testing"

func TestIsOverlapping(t *testing.T) {
	// s := make([]int, 6)
	// a := s[0:3]
	// b := s[3:6]

	a := make([]int, 3)
	b := make([]int, 3)

	t.Run("not a, b", func(t *testing.T) {
		if isOverlapping(a, b) {
			t.Fail()
		}
	})

	t.Run("a, a", func(t *testing.T) {
		if !isOverlapping(a, a) {
			t.Fail()
		}
	})

	t.Run("b, b", func(t *testing.T) {
		if !isOverlapping(b, b) {
			t.Fail()
		}
	})
}
