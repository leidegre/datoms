package sort

import "github.com/leidegre/datoms/internal/base"

// This function may have to do lookups to resolve some entity IDs
func ResolveEntid(v any) int64 {
	switch v := v.(type) {
	case uint32:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	default:
		panic("cannot resolve entity id")
	}
}

func Target(index base.Index, components []any) (target base.Datom) {
	switch index {
	case base.EAVT:
		switch len(components) {
		case 4:
			target.T = ResolveEntid(components[3])
			fallthrough
		case 3:
			target.V = components[2] // resolveValue?
			fallthrough
		case 2:
			target.A = ResolveEntid(components[1])
			fallthrough
		case 1:
			target.E = ResolveEntid(components[0])
		}
		return
	default:
		panic("todo")
	}
}

func TakeWhile(index base.Index, components []any) func(d base.Datom) bool {
	t := Target(index, components)
	switch index {
	case base.EAVT:
		switch len(components) {
		case 0:
			return func(d base.Datom) bool { return true }
		case 1:
			return func(d base.Datom) bool { return CompareOrdered(d.E, t.E) == 0 }
		case 2:
			return func(d base.Datom) bool {
				return CompareOrdered(d.E, t.E) == 0 && CompareOrdered(d.A, t.A) == 0
			}
		default:
			panic("uh-oh")
		}
	default:
		panic("uh-oh")
	}
}
