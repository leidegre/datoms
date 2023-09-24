package iter

// The Go authors are in the process of standardizing an iterator protocol

// https://github.com/golang/go/issues/61405
// https://github.com/golang/go/issues/61405#issuecomment-1638896606
// https://github.com/golang/go/issues/61897

// Once this has been merged we can use the range over func feature
// https://go-review.googlesource.com/c/go/+/510541
// https://go-review.googlesource.com/c/go/+/510541/16/src/cmd/compile/internal/rangefunc/rewrite.go

// Seq is a standard iterator and can be thought of as “push iterator”, which push values to the yield function.
type Seq[V any] func(yield func(v V) bool)

// Forward returns a forward iterator over a slice
func Forward[V any](s []V) Seq[V] {
	return func(yield func(v V) bool) {
		for i := 0; i < len(s); i++ {
			yield(s[i])
		}
	}
}

// Backward returns a backward iterator over a slice
func Backward[V any](s []V) Seq[V] {
	return func(yield func(v V) bool) {
		for i := len(s) - 1; 0 < i; i-- {
			yield(s[i])
		}
	}
}

// Range produces a sequence from the inclusive lower bound to the exclusive upper bound.
// Range(0, 10) is the equivalent of for i := 0; i < 10; i++ { ... }
func Range(i, j int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for ; i < j; i++ {
			if !yield(i) {
				break
			}
		}
	}
}

func Filter[V any](seq Seq[V], pred func(v V) bool) Seq[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if pred(v) {
				yield(v)
			}
			return true
		})
	}
}

func TakeWhile[V any](seq Seq[V], pred func(v V) bool) Seq[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if pred(v) {
				yield(v)
				return true
			}
			return false
		})
	}
}

// Slice concatenate all values in sequence to a slice
func Slice[V any](seq Seq[V]) []V {
	var tmp []V
	seq(func(v V) bool {
		tmp = append(tmp, v)
		return true
	})
	return tmp
}
