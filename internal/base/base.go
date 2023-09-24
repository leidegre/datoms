package base

import (
	"github.com/leidegre/datoms/symbol"
)

type Index uint8

const (
	EAVT Index = iota // Entity-Attribute-Value (document like)
	AEVT              // Attribute-Entity-Value (column like)
	AVET              // Attribute-Value-Entity (key-value like)
	VAET              // Value-Attribute-Entity (graph like)
)

type Entid interface {
	Zero() bool

	entid()
}

// todo: lookup refs?

// entity base...
type Entity struct {
	Id    int64          `ident:":db/id"`
	Ident symbol.Keyword `ident:":db/ident"`
}

type EntityLike interface {
	Entid

	entity() Entity

	// SetId(id int64)
	// SetIdent(ident symbol.Keyword)
}

func (e Entity) Zero() bool     { return e == Entity{} }
func (e Entity) entid()         {}
func (e Entity) entity() Entity { return e }

// func (e *Entity) SetId(id int64)                { e.Id = id }
// func (e *Entity) SetIdent(ident symbol.Keyword) { e.Ident = ident }

// This might not be needed?
func EntityIdentities(e EntityLike) (int64, symbol.Keyword) {
	ent := e.entity()
	return ent.Id, ent.Ident
}
