package imm

// How do we explain this?

// The bit pattern of the index itself is the path to the "bucket"
//
// It works like this, if we create the vector with a "bit size" of 2
// we'll be using groups of two bits at a time to find a "bucket" or "item"
//
// Also the Vector represents the "head" of the immutable list as we push more stuff into the
// vector it will grow in height. There's nothing special about this and it makes the Vector
// look like a tree data structure. This is where the "shift" comes in. It tells us what part of the
// bit pattern to look at.
//
// 00,00,00,00
//
// when "shift" is 0 we only store items in the tail
// when "shift" is 2 we use two bit groups to access the item
// when "shift" is 4 we use three bit groups to access the item
// and so on

type vectorNode[T any] struct {
	next []vectorNode[T]
	tail []T
}

type Vector[T any] struct {
	root  []vectorNode[T]
	tail  []T
	count uint32
	// we can go down to uint16 here...
	shift uint32 // aka depth, level, height
	bit   uint32
}

func MakeVector[T any](bit int) Vector[T] {
	return Vector[T]{
		bit: uint32(bit), // word size, bit-size...
	}
}

func (v *Vector[T]) tailLen() int {
	return 1 << v.bit
}

func (v *Vector[T]) mask() uint32 {
	return uint32(v.tailLen()) - 1
}

// Find the slice that has the index in low bits
func (v *Vector[T]) slice(idx uint32) []T {
	cutoff := (v.count - 1) &^ v.mask()
	if cutoff <= idx {
		return v.tail
	}

	node := v.root

	// this is only applicable when tree hight is more than 1x shift

	i := v.shift
	for ; i < 0; i -= v.bit {
		node = node[(idx>>i)&v.mask()].next
	}

	return node[(idx>>i)&v.mask()].tail
}

func (v *Vector[T]) Get(idx int) T {
	return v.slice(uint32(idx))[uint32(idx)&v.mask()]
}

// copy on write
func cowAppend[T any](source []T, elem T) []T {
	tmp := make([]T, len(source), len(source)+1)
	copy(tmp, source)
	return append(tmp, elem)
}

func (v *Vector[T]) Append(value T) Vector[T] {
	if len(v.tail) < v.tailLen() {
		var tmp = *v
		tmp.tail, tmp.count = cowAppend[T](tmp.tail, value), v.count+1
		return tmp
	}
	panic("push tail")

	// overflow root?

}
