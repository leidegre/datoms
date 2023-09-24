package database

import (
	"reflect"
	"testing"

	"github.com/leidegre/datoms/internal/base"
)

func TestMap(t *testing.T) {
	type Bar struct {
		base.Entity
		Bar string `ident:":bar"`
	}

	type Foo struct {
		base.Entity
		Doc string `ident:":db/doc"`
		Foo *Bar   `ident:":foo"`
	}

	var foo Foo

	foo.Doc = "foo"

	v := reflect.ValueOf(foo)

	c := contractType(v.Type())

	t.Fatal(c)
}
