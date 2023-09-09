package imm

import (
	"testing"

	"github.com/leidegre/datoms/testutil"
)

func TestVectorDataLayout(t *testing.T) {
	// Here we are using a very small bit size for testing purposes

	t.Run("2", func(t *testing.T) {
		var v Vector[int]

		v.tail = []int{2, 3}
		v.count = 2
		v.bit = 2

		testutil.AreEqual(t, v.Get(0), 2)
		testutil.AreEqual(t, v.Get(1), 3)
	})

	t.Run("5", func(t *testing.T) {
		var v Vector[int]

		v.root = []vectorNode[int]{{tail: []int{2, 3, 5, 7}}}
		v.tail = []int{11}
		v.count = 5
		v.shift = 2
		v.bit = 2

		testutil.AreEqual(t, v.Get(0), 2)  // 0b00_00
		testutil.AreEqual(t, v.Get(1), 3)  // 0b00_01
		testutil.AreEqual(t, v.Get(2), 5)  // 0b00_10
		testutil.AreEqual(t, v.Get(3), 7)  // 0b00_11
		testutil.AreEqual(t, v.Get(4), 11) // 0b01_00
	})
}

func TestVectorAppend(t *testing.T) {
	var v = MakeVector[int](2)

	v = v.Append(1)
	v = v.Append(2)
	v = v.Append(3)
	v = v.Append(4)
	v = v.Append(5)

	//...
}
