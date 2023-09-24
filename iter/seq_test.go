package iter_test

import (
	"testing"

	"github.com/leidegre/datoms/iter"
	"github.com/leidegre/datoms/testutil"
)

func TestRange(t *testing.T) {
	var n, s int
	iter.Range(1, 4)(func(v int) bool {
		n, s = n+1, s+v
		return true
	})
	testutil.AreEqual(t, 3, n)
	testutil.AreEqual(t, 1+2+3, s)
}

func TestFilter(t *testing.T) {
	var n, s int
	iter.Filter(iter.Range(0, 10), func(v int) bool {
		return v%2 == 0
	})(func(v int) bool {
		n, s = n+1, s+v
		return true
	})
	testutil.AreEqual(t, 5, n)
	testutil.AreEqual(t, 0+2+4+6+8, s)
}

func TestSlice(t *testing.T) {
	testutil.AreEqualSlice(t, []int{1, 2, 3}, iter.Slice(iter.Range(1, 4)))
}
