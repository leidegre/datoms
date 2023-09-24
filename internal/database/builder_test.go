package database

import (
	"testing"

	"github.com/leidegre/datoms/internal/base"
)

func TestBuilder(t *testing.T) {
	type foo struct {
		base.Entity
		Foo string `ident:":foo"`
	}

	var tx txBuilder

	tx.init(NewTestDatabase())

	tx.txExpand(nil, foo{Foo: "foo"})
}
