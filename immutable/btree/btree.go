package btree

import (
	"github.com/leidegre/datoms/immutable/vector"
	"github.com/leidegre/datoms/iter"
)

type Node[T any] interface {
	node()
}

// Persistent is not comparable because it must be initialized with a compare function. If you want to test for equality test the Roots of two sets for equality instead.
type Persistent[T any] struct {
	Root    Node[T] // root must not be nil
	compare func(a, b T) int
}

// Add
func (set Persistent[T]) Add(v T) Persistent[T] {
	return set
}

func (set Persistent[T]) Seek(v T) iter.Seq[T] {
	return func(yield func(v T) bool) {
		// ...
	}
}

// // Delete do we need delete?
// func (set Set[T]) Delete(v T) Set[T] {
// 	return set
// }

func New[T any](compare func(a, b T) int) Persistent[T] {
	return Persistent[T]{Root: &leafNode[T]{}, compare: compare}
}

type leafNode[T any] struct {
	vector.Persistent[T]
}

func (n *leafNode[T]) node() {}

type internalNode[T any] struct {
	vector.Persistent[Node[T]]
}

func (n *internalNode[T]) node() {}
