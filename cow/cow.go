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

// Copy-on-write insert element at index
func Insert[T any](slice []T, index int, elem T) []T {
	var tmp = make([]T, len(slice)+1)
	copy(tmp, slice[:index])
	tmp[index] = elem
	copy(tmp[index+1:], slice[index:])
	return tmp
}

// Copy-on-write delete element at index
func Delete[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j] // bounds check

	var tmp = make([]E, len(s)-1)
	copy(tmp, s[:i])
	copy(tmp[i:], s[j:])
	return tmp
}
