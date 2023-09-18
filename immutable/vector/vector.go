package vector

import (
	"github.com/leidegre/datoms/cow"
)

const (
	DefaultBranchingFactor = 5
)

// union type. either tail or next is used never both.
type node[T any] struct {
	// I benchmarked using interface{} and using type assertions to access as either []T or []node[T]
	// it didn't really make a performance difference but it did increase memory pressure by some 20-60%

	tail []T
	next []node[T]
}

// Persistent vector
type Persistent[T any] struct {
	root  []node[T]
	tail  []T
	count uint32
	shift uint16 // height of tree, starts at bits then is increased by bits for each new tree
	bits  uint16
}

// Make an empty vector with a specific branching factor. Ideally bits should be a value between 4 and 6.
func Make[T any](bits int) Persistent[T] {
	return Persistent[T]{
		shift: uint16(bits),
		bits:  uint16(bits),
	}
}

func (v *Persistent[T]) mask() uint32 {
	return (1 << v.bits) - 1
}

// Find the slice that has the element at index (lsb bits)
func (v *Persistent[T]) slice(idx uint32) []T {
	mask := v.mask()
	cutoff := (v.count - 1) &^ mask
	if cutoff <= idx {
		return v.tail
	}
	var (
		node  = v.root
		shift = v.shift
	)
	for ; v.bits < shift; shift -= v.bits {
		node = node[(idx>>shift)&mask].next
	}
	return node[(idx>>shift)&mask].tail
}

func (v *Persistent[T]) Len() int {
	return int(v.count)
}

func (v *Persistent[T]) Get(index int) T {
	return v.slice(uint32(index))[uint32(index)&v.mask()]
}

func (v *Persistent[T]) Append(value T) Persistent[T] {
	// https://github.com/clojure/clojure/blob/56d37996b18df811c20f391c840e7fd26ed2f58d/src/jvm/clojure/lang/PersistentVector.java#L222-L247

	// Make the zero value useful
	if v.bits == 0 {
		v.shift = DefaultBranchingFactor
		v.bits = DefaultBranchingFactor
	}

	if len(v.tail) < (1 << v.bits) {
		return Persistent[T]{
			root:  v.root,
			tail:  cow.Append(v.tail, value),
			count: v.count + 1,
			shift: v.shift,
			bits:  v.bits,
		}
	}

	// As long as there are vacant leaf nodes we want to use appendTail to fill these.
	// We don't need to recursively search for these since we compute the overflow condition numerically.

	if (v.count >> v.bits) <= (1 << v.shift) {
		return Persistent[T]{
			root:  v.appendTail(v.root, v.tail, v.count, v.shift),
			tail:  []T{value},
			count: v.count + 1,
			shift: v.shift,
			bits:  v.bits,
		}
	} else {
		var newRoot = []node[T]{
			{next: v.root},
			v.makePath(v.tail, v.shift),
		}
		return Persistent[T]{
			root:  newRoot,
			tail:  []T{value},
			count: v.count + 1,
			shift: v.shift + v.bits,
			bits:  v.bits,
		}
	}
}

func (v *Persistent[T]) makePath(tail []T, shift uint16) node[T] {
	if shift == 0 {
		return node[T]{tail: tail}
	} else {
		return node[T]{next: []node[T]{v.makePath(tail, shift-v.bits)}}
	}
}

func (v *Persistent[T]) appendTail(parent []node[T], tail []T, count uint32, shift uint16) []node[T] {
	// https://github.com/clojure/clojure/blob/56d37996b18df811c20f391c840e7fd26ed2f58d/src/jvm/clojure/lang/PersistentVector.java#L249-L270

	// The initial condition for the PersistentVector (as implemented in Clojure)
	// has a pre-allocated root with 32 entries.
	//
	// It does null checks instead of checking if the index is out of range.

	idx := int(((count - 1) >> shift) & ((1 << v.bits) - 1))

	if shift == v.bits {
		// invariant
		if !(idx == len(parent)) {
			panic("uh-oh")
		}
		return cow.Append(parent, node[T]{tail: tail})
	} else {
		if idx < len(parent) {
			var newTail = v.appendTail(parent[idx].next, tail, count, shift-v.bits)

			return cow.Update(parent, idx, node[T]{next: newTail})
		} else {
			var newTail = v.makePath(tail, shift-v.bits)

			return cow.Append(parent, newTail)
		}
	}
}

// Iterator protocol

type Iter[T any] struct {
	v        *Persistent[T]
	tail     []T
	idx, end uint32
}

// Reports whether the iterator is valid
func (it *Iter[T]) Valid() bool {
	return it.idx < it.end
}

// Move the iterator forward
func (it *Iter[T]) Next() {
	it.idx++
	if (it.idx & it.v.mask()) == 0 {
		it.tail = it.v.slice(it.idx)
	}
}

// Get the value associated with iterator. Requires: Valid()
func (it *Iter[T]) Value() T {
	return it.tail[it.idx&it.v.mask()]
}

func (v *Persistent[T]) Range(min, max int) Iter[T] {
	return Iter[T]{v: v, tail: v.slice(uint32(min)), idx: uint32(min), end: uint32(max)}
}

// Transient

type Transient[T any] struct {
	Persistent[T]
}

func (v *Transient[T]) Append(value T) {
	if v.bits == 0 {
		v.shift = DefaultBranchingFactor
		v.bits = DefaultBranchingFactor
	}

	if len(v.tail) < (1 << v.bits) {
		v.tail = append(v.tail, value)
		v.count++
		return
	}

	if (v.count >> v.bits) <= (1 << v.shift) {
		v.root = v.appendTail(v.root, v.tail, v.count, v.shift)
		v.tail = []T{value}
		v.count++
	} else {
		v.root = []node[T]{{next: v.root}, v.makePath(v.tail, v.shift)}
		v.tail = []T{value}
		v.count++
		v.shift += v.bits
	}
}

func (v *Transient[T]) makePath(tail []T, shift uint16) node[T] {
	if shift == 0 {
		return node[T]{tail: tail}
	} else {
		return node[T]{next: []node[T]{v.makePath(tail, shift-v.bits)}}
	}
}

func (v *Transient[T]) appendTail(parent []node[T], tail []T, count uint32, shift uint16) []node[T] {
	idx := int(((count - 1) >> shift) & ((1 << v.bits) - 1))
	if shift == v.bits {
		// invariant
		if !(idx == len(parent)) {
			panic("uh-oh")
		}
		return append(parent, node[T]{tail: tail})
	} else {
		if idx < len(parent) {
			var newTail = v.appendTail(parent[idx].next, tail, count, shift-v.bits)
			parent[idx] = node[T]{next: newTail}
			return parent
		} else {
			var newTail = v.makePath(tail, shift-v.bits)
			return append(parent, newTail)
		}
	}
}

// Note that a transient vector should be treated as a temporary. After calling Immutable the transient vector is no longer valid.
func (v *Transient[T]) Immutable() Persistent[T] {
	var tmp = v.Persistent
	v.Persistent = Persistent[T]{} // prevent awkward bugs
	return tmp
}
