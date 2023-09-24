package iterutil

import (
	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/internal/sort"
	"github.com/leidegre/datoms/iter"
)

func Live(seq iter.Seq[base.Datom]) iter.Seq[base.Datom] {
	var r base.Datom
	return func(yield func(v base.Datom) bool) {
		seq(func(d base.Datom) bool {
			if d.Assertion() {
				if d.E == r.E && d.A == r.A && sort.CompareValue(d.V, r.V) == 0 {
					r = base.Datom{}
					return true // continue
				}
				yield(d)
			} else {
				r = d
			}
			return true
		})
	}
}
