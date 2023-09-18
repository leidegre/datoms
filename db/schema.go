package db

import "github.com/leidegre/datoms/symbol"

type Schema interface {
	Id(kw symbol.Keyword) EntityId
	Attr(attrId EntityId) Attribute
}
