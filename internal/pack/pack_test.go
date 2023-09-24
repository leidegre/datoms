package pack_test

import (
	"testing"

	"github.com/leidegre/datoms/internal/pack"
	"github.com/leidegre/datoms/testutil"
)

type test struct {
	packed, part, ent int64
}

var tests = []test{
	{1, 1, 1},
	{2, 1, 2},
	{3, 1, 3},

	{4398046511104 + 1, 2, 1},
	{4398046511104 + 2, 2, 2},
	{4398046511104 + 3, 2, 3},

	{(3-1)*4398046511104 + 1, 3, 1},
	{(3-1)*4398046511104 + 2, 3, 2},
	{(3-1)*4398046511104 + 3, 3, 3},
}

func TestEntityId(t *testing.T) {
	for _, test := range tests {
		testutil.AreEqual(t, test.packed, pack.EntityId(test.part, test.ent))
	}
}

func TestTempId(t *testing.T) {
	for _, test := range tests {
		tempId := pack.TempId(test.part, test.ent)
		if !(tempId < 0) {
			t.Fatal("temp ID should be negative, non-zero value")
		}
		part, ent := pack.Unpack(tempId)
		testutil.AreEqual(t, test.part, part)
		testutil.AreEqual(t, test.ent, ent)
	}
}

func TestUnpack(t *testing.T) {
	for _, test := range tests {
		part, ent := pack.Unpack(test.packed)
		testutil.AreEqual(t, test.part, part)
		testutil.AreEqual(t, test.ent, ent)
	}
}

func TestOutOfRange(t *testing.T) {
	assertPanic := func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}

	t.Run("00", func(t *testing.T) {
		defer assertPanic()

		pack.EntityId(0, 0)
	})

	t.Run("01", func(t *testing.T) {
		defer assertPanic()

		pack.EntityId(0, 1)
	})

	t.Run("10", func(t *testing.T) {
		defer assertPanic()

		pack.EntityId(1, 0)
	})

	t.Run("x1", func(t *testing.T) {
		defer assertPanic()

		pack.EntityId((1<<20)+1, 1)
	})

	t.Run("1x", func(t *testing.T) {
		defer assertPanic()

		pack.EntityId(1, 1<<42)
	})
}
