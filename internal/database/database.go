package database

import (
	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/internal/schema"
	"github.com/leidegre/datoms/iter"
)

type Transaction struct {
	DbBefore Interface
	DbAfter  Interface
	TxData   []base.Datom
	TempIds  map[int64]int64
}

type Interface interface {
	T() (baseT, nextT int64)

	Schema() schema.Interface

	SeekDatoms(index base.Index, components ...any) iter.Seq[base.Datom]

	Datoms(index base.Index, components ...any) iter.Seq[base.Datom]

	With(txData []base.TxData) (Transaction, error)
}
