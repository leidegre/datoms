// see https://lampwww.epfl.ch/papers/idealhashtrees.pdf
package hashmap

import (
	"fmt"
	"math/bits"
	"slices"
	"strings"

	"github.com/leidegre/datoms/cow"
)

// We have three types of internal nodes to contend with these implement this interface.
type node[K comparable, V comparable] interface {
	get(k K, hash, shift uint64) (v V, ok bool)
	set(k K, hash, shift uint64, v V, inserted *uint32) node[K, V]
	delete(k K, hash, shift uint64) node[K, V]
}

const (
	shiftInc        = 5 // some value between 1 and 6
	BranchingFactor = 1 << shiftInc
	mask            = BranchingFactor - 1
)

func bitchk(hash uint64, shift uint64) uint64 {
	return (hash >> shift) & mask
}

func bitpos(hash uint64, shift uint64) uint64 {
	return uint64(1) << bitchk(hash, shift)
}

func index(bitmap uint64, bit uint64) int {
	return bits.OnesCount64(bitmap & (bit - 1))
}

// Persistent is a immutable hash map based on Phil Bagwell's paper Ideal Hash Trees.
type Persistent[K comparable, V comparable] struct {
	// Since the hashmap is passed around by value it's a good thing that it has a small footprint.

	root node[K, V]
	size uint32
}

// What would be idiomatic Go? Len, Size or Count? Go hasn't historically been keen on a containers...
func (m Persistent[K, V]) Len() int {
	return int(m.size)
}

func (m Persistent[K, V]) Get(k K, h uint64) (v V, ok bool) {
	if m.root == nil {
		return
	}
	v, ok = m.root.get(k, h, 0)
	return
}

func (m Persistent[K, V]) Set(k K, h uint64, v V) Persistent[K, V] {
	if m.root == nil {
		m.root = &bitmapNode[K, V]{}
	}

	var inserted uint32

	m.root = m.root.set(k, h, 0, v, &inserted)
	m.size = m.size + inserted

	return m
}

func (m Persistent[K, V]) Delete(k K, h uint64) Persistent[K, V] {
	if m.root == nil {
		return m
	}
	m.root = m.root.delete(k, h, 0)
	m.size = m.size - 1
	return m
}

func debugString[K comparable, V comparable](lvl int, root node[K, V]) (s string) {
	switch n := root.(type) {
	case *bitmapNode[K, V]:
		s += fmt.Sprintf("(:bitmap 0x%x [", n.bitmap)
		for _, n := range n.data {
			s += "\n"
			s += strings.Repeat(" ", 2*(lvl+1))
			s += debugString(lvl+1, n)
		}
		s += "\n" + strings.Repeat(" ", 2*lvl) + "])"
	case *valueNode[K, V]:
		s += fmt.Sprintf("{%#v %v %#v}", n.k, n.h, n.v)
	case *bucketNode[K, V]:
		s += "(:bucket ["
		for i, collision := range n.collisions {
			if 0 < i {
				s += " "
			}
			s += debugString(lvl+1, collision)
		}
		s += "])"
	}
	return
}

func DebugString[K comparable, V comparable](m Persistent[K, V]) string {
	return "\n" + debugString(0, m.root)
}

type bitmapNode[K comparable, V comparable] struct {
	bitmap uint64
	data   []node[K, V]
}

func (n *bitmapNode[K, V]) get(k K, hash, shift uint64) (v V, ok bool) {
	var (
		bitmap = n.bitmap
		bit    = bitpos(hash, shift)
	)
	if bitmap&bit == bit {
		var idx = index(bitmap, bit)

		v, ok = n.data[idx].get(k, hash, shift+shiftInc)
	}
	return
}

func (n *bitmapNode[K, V]) set(k K, hash uint64, shift uint64, v V, inserted *uint32) node[K, V] {
	var (
		bitmap = n.bitmap
		bit    = bitpos(hash, shift)
		idx    = index(bitmap, bit)
	)

	if bitmap&bit == 0 {
		// This bit position is vacant
		*inserted = 1
		var newNode node[K, V] = &valueNode[K, V]{k, hash, v}
		return &bitmapNode[K, V]{bitmap: bitmap | bit, data: cow.Insert(n.data, idx, newNode)}
	} else {
		var ndi = n.data[idx]
		newNode := ndi.set(k, hash, shift+shiftInc, v, inserted)
		if newNode == ndi {
			return n // no change
		}
		return &bitmapNode[K, V]{bitmap: bitmap, data: cow.Update(n.data, idx, newNode)}
	}
}

func (n *bitmapNode[K, V]) delete(k K, hash uint64, shift uint64) node[K, V] {
	var (
		bitmap = n.bitmap
		bit    = bitpos(hash, shift)
		idx    = index(bitmap, bit)
	)

	if bitmap&bit == 0 {
		return n // not found
	} else {
		// todo: degenerate case
		//       removing {1 13920995807293724370 1} will still
		//       maintain an unnecessary deep tree

		// (:bitmap 0x8040000 [
		// 	{1 13920995807293724370 1}
		// 	(:bitmap 0x10000 [
		// 		{0 7137908203259007515 0}
		// 	])
		// ])

		// (:bitmap 0x8000000 [
		// 	(:bitmap 0x10000 [
		// 		{0 7137908203259007515 0}
		// 	])
		// ])

		// If we have two disjoint bitmaps with a small enough combined population
		// count we should merge them.

		// To cheaply merge with parent we need to maintain a compatible bitmap
		// I don't think this is trivial.

		var ndi = n.data[idx]
		next := ndi.delete(k, hash, shift+shiftInc)
		if next != ndi {
			if next == nil {
				if len(n.data) == 1 {
					return nil
				}
				return &bitmapNode[K, V]{bitmap &^ bit, cow.Delete(n.data, idx, idx+1)}
			}
			return &bitmapNode[K, V]{bitmap, cow.Update(n.data, idx, next)}
		}
		return n
	}
}

type valueNode[K comparable, V comparable] struct {
	k K
	h uint64
	v V
}

func (n *valueNode[K, V]) get(k K, hash, _ uint64) (v V, ok bool) {
	if n.k == k {
		v, ok = n.v, true
	}
	return
}

func (n *valueNode[K, V]) set(k K, h, shift uint64, v V, inserted *uint32) node[K, V] {
	if n.h == h {
		if n.k == k {
			if n.v == v {
				return n // no change
			}
			return &valueNode[K, V]{k, h, v}
		}
		// hash collision
		*inserted = 1
		return &bucketNode[K, V]{[]*valueNode[K, V]{n, {k, h, v}}}
	}
	var split node[K, V] = &bitmapNode[K, V]{bitpos(n.h, shift), []node[K, V]{n}}
	return split.set(k, h, shift, v, inserted)
}

func (n *valueNode[K, V]) delete(k K, hash uint64, shift uint64) node[K, V] {
	if n.k == k {
		return nil
	}
	return n
}

type bucketNode[K comparable, V comparable] struct {
	collisions []*valueNode[K, V]
}

func (n *bucketNode[K, V]) get(k K, _, _ uint64) (v V, ok bool) {
	for _, collision := range n.collisions {
		if collision.k == k {
			v, ok = collision.v, true
			break
		}
	}
	return
}

func (n *bucketNode[K, V]) set(k K, h, s uint64, v V, inserted *uint32) node[K, V] {
	// a hash collision bucket is never empty
	collisionHash := n.collisions[0].h
	if collisionHash == h {
		for i, collision := range n.collisions {
			if collision.k == k {
				return &bucketNode[K, V]{cow.Update(n.collisions, i, collision)}
			}
		}
		*inserted = 1
		return &bucketNode[K, V]{cow.Append(n.collisions, &valueNode[K, V]{k, h, v})}
	}
	var split = &bitmapNode[K, V]{bitpos(collisionHash, s), []node[K, V]{n}}
	return split.set(k, h, s, v, inserted)
}

func (n *bucketNode[K, V]) delete(k K, hash uint64, shift uint64) node[K, V] {
	for i, collision := range n.collisions {
		if collision.k == k {
			if len(n.collisions) == 1 {
				return nil
			}
			return &bucketNode[K, V]{cow.Delete(n.collisions, i, i+1)}
		}
	}
	return n
}

type Transient[K comparable, V comparable] struct {
	Persistent[K, V]
}

func (t *Transient[K, V]) Set(k K, h uint64, v V) {
	var inserted uint32
	t.root = setTransient(t.root, k, h, v, 0, &inserted)
	t.size = t.size + inserted
}

func (v *Transient[K, V]) Immutable() Persistent[K, V] {
	var tmp = v.Persistent
	v.Persistent = Persistent[K, V]{} // prevent awkward bugs
	return tmp
}

func setTransient[K comparable, V comparable](n node[K, V], k K, h uint64, v V, s uint64, inserted *uint32) node[K, V] {
	switch n := n.(type) {
	case nil:
		return &bitmapNode[K, V]{bitpos(h, s), []node[K, V]{&valueNode[K, V]{k, h, v}}}
	case *bitmapNode[K, V]:
		var (
			bitmap = n.bitmap
			bit    = bitpos(h, s)
			idx    = index(bitmap, bit)
		)
		if bitmap&bit == 0 {
			*inserted = 1
			n.bitmap, n.data = bitmap|bit, slices.Insert(n.data, idx, node[K, V](&valueNode[K, V]{k, h, v}))
		} else {
			n.data[idx] = setTransient(n.data[idx], k, h, v, s+shiftInc, inserted)
		}
		return n
	case *valueNode[K, V]:
		if n.h == h {
			if n.k == k {
				n.v = v // update
				return n
			}
			// collision
			*inserted = 1
			newNode := &bucketNode[K, V]{[]*valueNode[K, V]{n, {k, h, v}}}
			return newNode
		}
		newNode := &bitmapNode[K, V]{bitpos(n.h, s), make([]node[K, V], 1, 2)}
		newNode.data[0] = n
		return setTransient(newNode, k, h, v, s, inserted)
	default:
		panic("uh-oh!")
	}
}
