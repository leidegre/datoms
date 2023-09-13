// Copy-on-write slice utilities
package cow

func ShallowCopy[T any](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}
	var tmp = make([]T, len(slice))
	copy(tmp, slice)
	return tmp
}

// Copy-on-write append
func Append[T any](slice []T, elems ...T) []T {
	var tmp = make([]T, len(slice), len(slice)+len(elems))
	copy(tmp, slice)
	return append(tmp, elems...)
}

// Copy-on-write set element at index
func Update[T any](slice []T, index int, elem T) []T {
	var tmp = ShallowCopy(slice)
	tmp[index] = elem
	return tmp
}
