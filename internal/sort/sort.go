package sort

import (
	"cmp"
	"time"

	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/symbol"
)

func bool2int(b bool) int {
	// https://cs.opensource.google/go/go/+/master:src/cmd/internal/obj/util.go;l=652?q=Bool2int&sq=&ss=go%2Fgo
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}

// bool is comparable but not ordered
func CompareBool[T ~bool](x, y T) int {
	if !x {
		if !y {
			return 0 // both false
		}
		return -1 // x is false, y is true
	}
	if !y {
		return 1 // x is true, y is false
	}
	return 0 // both true
}

func CompareOrdered[T ~uint32 | ~int32 | ~uint64 | ~int64 | ~float32 | ~float64 | ~string](x, y T) int {
	if x < y {
		return -1
	}
	if y < x {
		return +1
	}
	return 0
}

// CompareValue compares values of matching type. This is enforced by our
// index sort order, attributes always precedes values. There's never going
// to be a situation where the left type doesn't match the right type.
func CompareValue(x, y interface{}) int {
	switch x := x.(type) {
	case bool:
		y := y.(bool)
		return CompareBool(x, y)
	case float32:
		y := y.(float32)
		return cmp.Compare(x, y) // NaN handling
	case float64:
		y := y.(float64)
		return cmp.Compare(x, y) // NaN handling
	case string:
		y := y.(string)
		return CompareOrdered(x, y)
	case int32:
		y := y.(int32)
		return CompareOrdered(x, y)
	case int64:
		y := y.(int64)
		return CompareOrdered(x, y)
	case symbol.Keyword:
		y := y.(symbol.Keyword)
		return CompareOrdered(x.String(), y.String())
	case time.Time:
		y := y.(time.Time)
		return x.Compare(y)
	default:
		panic("value is not comparable")
	}
}

// Live
func CompareIndex(index base.Index) func(x, y base.Datom) int {
	switch index {
	case base.EAVT:
		return CompareEAV
	case base.AEVT:
		return CompareAEV
	case base.AVET:
		return CompareAVE
	case base.VAET:
		return CompareVAE
	default:
		panic("datoms: invalid index")
	}
}

func CompareEAV(x, y base.Datom) int {
	e := CompareOrdered(x.E, y.E)
	if e != 0 {
		return e
	}
	a := CompareOrdered(x.A, y.A)
	if a != 0 {
		return a
	}
	v := CompareValue(x.V, y.V)
	return v
}

func CompareAEV(x, y base.Datom) int {
	a := CompareOrdered(x.A, y.A)
	if a != 0 {
		return a
	}
	e := CompareOrdered(x.E, y.E)
	if e != 0 {
		return e
	}
	v := CompareValue(x.V, y.V)
	return v
}

func CompareAVE(x, y base.Datom) int {
	a := CompareOrdered(x.A, y.A)
	if a != 0 {
		return a
	}
	v := CompareValue(x.V, y.V)
	if v != 0 {
		return v
	}
	e := CompareOrdered(x.E, y.E)
	return e
}

func CompareVAE(x, y base.Datom) int {
	v := CompareValue(x.V, y.V)
	if v != 0 {
		return v
	}
	a := CompareOrdered(x.A, y.A)
	if a != 0 {
		return a
	}
	e := CompareOrdered(x.E, y.E)
	return e
}

func CompareHistory(index base.Index) func(x, y base.Datom) int {
	switch index {
	case base.EAVT:
		return func(x, y base.Datom) int {
			eav := CompareEAV(x, y)
			if eav != 0 {
				return eav
			}
			t := CompareOrdered(y.T, x.T) // descending
			return t
		}
	case base.AEVT:
		return func(x, y base.Datom) int {
			aev := CompareAEV(x, y)
			if aev != 0 {
				return aev
			}
			t := CompareOrdered(y.T, x.T) // descending
			return t
		}
	case base.AVET:
		return func(x, y base.Datom) int {
			ave := CompareAVE(x, y)
			if ave != 0 {
				return ave
			}
			t := CompareOrdered(y.T, x.T) // descending
			return t
		}
	case base.VAET:
		return func(x, y base.Datom) int {
			vae := CompareVAE(x, y)
			if vae != 0 {
				return vae
			}
			t := CompareOrdered(y.T, x.T) // descending
			return t
		}
	default:
		panic("uh-oh")
	}
}
