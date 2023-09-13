package cow_test

import (
	"testing"
	"unsafe"

	"github.com/leidegre/datoms/cow"
	"github.com/leidegre/datoms/testutil"
)

// Assert Go behavior, what is nil and what is not
func TestSlice(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var s []int // nil pointer
		p := unsafe.SliceData(s)
		if p != nil {
			t.Fail()
		}
	})

	t.Run("empty", func(t *testing.T) {
		var s = []int{} // not a nil pointer
		p := unsafe.SliceData(s)
		if p == nil {
			t.Fail()
		}
	})

	t.Run("make0", func(t *testing.T) {
		s := make([]int, 0) // not a nil pointer
		p := unsafe.SliceData(s)
		if p == nil {
			t.Fail()
		}
	})
}

func TestShallowCopy(t *testing.T) {
	xs := []int{1, 2, 3}
	ys := cow.ShallowCopy(xs)
	ys[0] = 4
	ys[1] = 5
	ys[2] = 6

	testutil.AreEqual(t, 1, xs[0])
	testutil.AreEqual(t, 2, xs[1])
	testutil.AreEqual(t, 3, xs[2])

	testutil.AreEqual(t, 4, ys[0])
	testutil.AreEqual(t, 5, ys[1])
	testutil.AreEqual(t, 6, ys[2])
}

func TestAppend(t *testing.T) {
	a := cow.Append[int](nil, 1)
	b := cow.Append[int](a, 2)
	c := cow.Append[int](b, 3)

	testutil.AreDistinctSlice(t, a, b)
	testutil.AreDistinctSlice(t, b, c)
	testutil.AreDistinctSlice(t, a, c)

	testutil.AreEqualSlice(t, []int{1}, a)
	testutil.AreEqualSlice(t, []int{1, 2}, b)
	testutil.AreEqualSlice(t, []int{1, 2, 3}, c)
}
