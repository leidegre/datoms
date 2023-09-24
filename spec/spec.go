package spec

// import (
// 	"testing"

// 	"github.com/leidegre/datoms/datoms"
// 	"github.com/leidegre/datoms/testutil"
// )

// func Run(t *testing.T, newDatabase func() datoms.Database) {
// 	t.Run("NewDatabase", func(t *testing.T) {
// 		db := newDatabase()

// 		baseT, _ := db.T()

// 		testutil.AreEqual(t, 1000, baseT)
// 	})

// 	t.Run("With foo bar baz", func(t *testing.T) {
// 		db := newDatabase()

// 		t1 := datoms.NewTempId(datoms.PartUser)
// 		t2 := datoms.NewTempId(datoms.PartUser)
// 		t3 := datoms.NewTempId(datoms.PartUser)

// 		tx := db.With([]datoms.TxData{
// 			datoms.Add(t1, datoms.Doc, "foo"),
// 			datoms.Add(t2, datoms.Doc, "bar"),
// 			datoms.Add(t3, datoms.Doc, "baz")})

// 		t.Log(tx.TempIds)
// 		t.FailNow()
// 	})
// }
