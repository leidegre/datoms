package base

import (
	"sync/atomic"

	"github.com/leidegre/datoms/symbol"
)

type TempId struct {
	Part   Entity
	TempId int64
}

func (tempId TempId) Zero() bool { return tempId == TempId{} }
func (TempId) entid()            {}

var (
	tempId int64 = 1 << 40 // peer and transactor should use different range
)

func NewTempId(part symbol.Keyword) TempId {
	return TempId{Part: Entity{Ident: part}, TempId: atomic.AddInt64(&tempId, 1)}
}
