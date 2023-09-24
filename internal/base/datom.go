package base

import (
	"fmt"

	"github.com/leidegre/datoms/internal/pack"
)

type Datom struct {
	E int64
	A int64
	V any
	T int64
}

func NewDatom(E int64, A int64, V any, T int64, op int64) Datom {
	return Datom{E, A, V, (T << 1) | (op & 1)}
}

func (d Datom) String() string {
	part, ent := pack.Unpack(d.E)
	return fmt.Sprintf("{%v:%v %v %v %v %v}", part, ent, d.A, d.V, d.T>>1, d.Assertion())
}

func (d Datom) Int() int64 {
	switch v := d.V.(type) {
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	default:
		panic("uh-oh")
	}
}

func (d *Datom) Tx() int64 {
	return d.T >> 1
}

func (d *Datom) Assertion() bool {
	return (d.T & 1) == 1
}

func (d *Datom) Retraction() bool {
	return (d.T & 1) == 0
}
