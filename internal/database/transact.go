package database

import (
	"fmt"
	"time"

	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/internal/pack"
	"github.com/leidegre/datoms/internal/schema"
	"github.com/leidegre/datoms/symbol"
)

func Add(e base.Entid, a symbol.Keyword, v interface{}) base.TxData {
	return base.TxAdd{E: e, A: a, V: v}
}

func Retract(e base.Entid, a symbol.Keyword, v interface{}) base.TxData {
	return base.TxRetract{E: e, A: a, V: v}
}

func Map(id base.Entid, e base.EntityLike) base.TxData {
	return base.TxMap{Id: id, Entity: e}
}

func Transact(db Interface, txData []base.TxData) (baseT int64, nextT int64, data []base.Datom, tempIds map[int64]int64, err error) {
	var tx txBuilder

	tx.init(db)

	// These are optional, like :db.install/attribute
	tx.emitEntidAttr(base.NewTempId(schema.DbPartTx), schema.DbTxInstant, time.Now(), 1)

	for _, item := range txData {
		switch item := item.(type) {
		case base.TxAdd:
			err = tx.emitEntidAttr(item.E, item.A, item.V, 1)
			if err != nil {
				return
			}
		case base.TxRetract:
			err = tx.emitEntidAttr(item.E, item.A, item.V, 0)
			if err != nil {
				return
			}
		case base.TxMap:
			_, err = tx.txExpand(item.Id, item.Entity)
			if err != nil {
				return
			}
		default:
			panic(fmt.Sprintf("datoms: unknown type %T in transaction data", item))
		}
	}

	for i, d := range tx.data {
		if d.E < 0 {
			var (
				newId int64
				ok    bool
			)
			if newId, ok = tx.tempIds[d.E]; !ok {
				part, _ := pack.Unpack(d.E)
				newId = pack.EntityId(part, tx.nextT) // make a new entity
				tx.tempIds[d.E] = newId
				tx.nextT++
			}
			tx.data[i] = base.Datom{E: newId, A: d.A, V: d.V, T: d.T}
		}
	}

	// ---

	// We need to make a pass to detect missing implicit datoms

	// We need to make a pass to eliminate redundant datoms

	// For this we sort of need the b-tree...

	// ---

	baseT, nextT, data, tempIds = nextT, tx.nextT, tx.data, tx.tempIds
	return
}
