package database_test

import (
	"testing"

	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/internal/database"
	"github.com/leidegre/datoms/internal/schema"
)

func TestTransactAdd(t *testing.T) {
	db := database.NewTestDatabase()

	t1 := base.NewTempId(schema.DbPartUser)
	t2 := base.NewTempId(schema.DbPartUser)
	t3 := base.NewTempId(schema.DbPartUser)

	tx, err := db.With([]base.TxData{
		database.Add(t1, schema.DbDoc, "foo"),
		database.Add(t2, schema.DbDoc, "bar"),
		database.Add(t3, schema.DbDoc, "baz"),
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(tx)

	t.FailNow()
}

func TestTransact(t *testing.T) {
	db := database.NewTestDatabase()

	t1 := base.NewTempId(schema.DbPartUser)
	t2 := base.NewTempId(schema.DbPartUser)
	t3 := base.NewTempId(schema.DbPartUser)

	type qux struct {
		base.Entity
		Doc string `ident:":db/doc"`
	}

	tx, err := db.With([]base.TxData{
		database.Add(t1, schema.DbDoc, "foo"),
		database.Add(t2, schema.DbDoc, "bar"),
		database.Add(t3, schema.DbDoc, "baz"),
		database.Map(nil, qux{Doc: "qux"}),
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(tx)

	t.FailNow()
}
