package iterutil_test

import (
	"slices"
	"testing"

	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/internal/iterutil"
	"github.com/leidegre/datoms/internal/sort"
	"github.com/leidegre/datoms/iter"
	"github.com/leidegre/datoms/testutil"
)

func datom(e int64, a int64, v interface{}, t int64, op bool) (d base.Datom) {
	d.E = e
	d.A = a
	d.V = v
	d.T = t << 1
	if op {
		d.T |= 1
	}
	return
}

func TestLive(t *testing.T) {
	// Datomic has these 4 indexes

	// EAVT
	// AEVT
	// AVET
	// VAET

	// We use the lsb bit of the T value to store whether this is an assertion (0x1) or retraction (0x0).

	// The same value cannot be both asserted and retracted in the same transaction.

	// The T component is sorted in descending order.

	// Thus the live iterator will see the retraction of a value first and
	// retract the value from the iteration.

	hist := []base.Datom{
		datom(0, 1, "foo", 0, true),
		datom(0, 1, "foo", 1, false),
		datom(0, 1, "bar", 1, true),
		datom(0, 1, "bar", 2, false),
		datom(0, 1, "baz", 2, true),
		datom(1, 1, "qux", 3, true),
	}

	slices.SortFunc(hist, sort.CompareHistory(base.EAVT))

	live := iter.Slice(iterutil.Live(iter.Forward(hist)))

	testutil.AreEqual(t, datom(0, 1, "baz", 2, true), live[0])
	testutil.AreEqual(t, datom(1, 1, "qux", 3, true), live[1])
}
