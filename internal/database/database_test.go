package database

import (
	"slices"

	"github.com/leidegre/datoms/cow"
	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/internal/iterutil"
	"github.com/leidegre/datoms/internal/schema"
	"github.com/leidegre/datoms/internal/sort"
	"github.com/leidegre/datoms/iter"
)

// A basic and NOT scalable database implementation for testing.
type TestDatabase struct {
	baseT, nextT int64
	data         []base.Datom
	schema       *schema.Schema
}

func (db *TestDatabase) T() (int64, int64) { return db.baseT, db.nextT }

func (db *TestDatabase) Schema() schema.Interface { return db.schema }

func (db *TestDatabase) SeekDatoms(index base.Index, components ...any) iter.Seq[base.Datom] {
	cmp := sort.CompareHistory(index)
	data := cow.ShallowCopy(db.data)
	slices.SortFunc(data, cmp)
	i, _ := slices.BinarySearchFunc(data, sort.Target(index, components), cmp)
	return iterutil.Live(iter.Forward(data[i:]))
}

func (db *TestDatabase) Datoms(index base.Index, components ...any) iter.Seq[base.Datom] {
	return iter.TakeWhile(db.SeekDatoms(index, components...), sort.TakeWhile(index, components))
}

func (db *TestDatabase) With(txData []base.TxData) (Transaction, error) {
	baseT, nextT, data, tempIds, err := Transact(db, txData)

	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		DbBefore: db,
		DbAfter:  &TestDatabase{baseT, nextT, cow.Append(db.data, data...), db.schema.With(data)},
		TxData:   data,
		TempIds:  tempIds,
	}, nil
}

// Create a new basic database for testing.
func NewTestDatabase() Interface {
	data := schema.BootstrappingPart(0)

	var db TestDatabase

	db.baseT = 1000
	db.nextT = 1000
	db.data = data
	db.schema = (&schema.Schema{}).With(data)

	return &db
}
