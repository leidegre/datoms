package schema

import (
	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/symbol"
)

type bootId uint32

// These values cannot change once a database is created
const (
	_ bootId = iota

	dbPartDb
	dbPartTx
	dbPartUser

	dbInstallValueType
	dbInstallPartition
	dbInstallAttribute

	dbIdent

	dbValueType
	dbTypeBool
	dbTypeFloat64
	dbTypeString
	dbTypeInt64
	dbTypeRef
	dbTypeKeyword
	dbTypeTime

	dbCardinality
	dbCardinalityOne
	dbCardinalityMany

	dbUnique
	dbUniqueValue
	dbUniqueIdentity

	dbIsComponent

	dbTxInstant

	dbDoc
)

var (
	DbId = symbol.For(":db/id") // entity ID

	DbPartDb   = symbol.For(":db.part/db")
	DbPartTx   = symbol.For(":db.part/tx")
	DbPartUser = symbol.For(":db.part/user")

	DbInstallPartition = symbol.For(":db.install/partition")
	DbInstallValueType = symbol.For(":db.install/valueType")
	DbInstallAttribute = symbol.For(":db.install/attribute")

	DbIdent = symbol.For(":db/ident")

	DbValueType   = symbol.For(":db/valueType")
	DbTypeBoolean = symbol.For(":db.type/boolean") // bool
	DbTypeDouble  = symbol.For(":db.type/double")  // float64
	DbTypeString  = symbol.For(":db.type/string")  // string
	DbTypeLong    = symbol.For(":db.type/long")    // int64
	DbTypeRef     = symbol.For(":db.type/ref")     // int64
	DbTypeKeyword = symbol.For(":db.type/keyword") // symbol.Keyword
	DbTypeInstant = symbol.For(":db.type/instant") // time.Time

	DbCardinality     = symbol.For(":db/cardinality")
	DbCardinalityOne  = symbol.For(":db.cardinality/one")
	DbCardinalityMany = symbol.For(":db.cardinality/many")

	DbUnique         = symbol.For(":db/unique")
	DbUniqueValue    = symbol.For(":db.unique/value")
	DbUniqueIdentity = symbol.For(":db.unique/identity")

	DbIsComponent = symbol.For(":db/isComponent")

	DbTx        = symbol.For(":db/tx") // ref to transaction entity when transacting
	DbTxInstant = symbol.For(":db/txInstant")

	DbDoc = symbol.For(":db/doc")
)

type bootstrappingPart struct {
	data []base.Datom
	base bootId
}

func (b *bootstrappingPart) add(e bootId, a bootId, v interface{}) {
	// all attributes of type ref must be offset from base...
	if v2, ok := v.(bootId); ok {
		v = int64(v2) // ref???
	}
	b.data = append(b.data, base.Datom{E: int64(e), A: int64(a), V: v, T: 1})
}

func (b *bootstrappingPart) defineEntity(id bootId, entIdent symbol.Keyword) {
	b.add(id, dbIdent, entIdent)
}

func (b *bootstrappingPart) definePartition(id bootId, partIdent symbol.Keyword) {
	b.defineEntity(id, partIdent)

	b.add(dbPartDb, dbInstallPartition, id)
}

func (b *bootstrappingPart) defineValueType(id bootId, typeIdent symbol.Keyword) {
	b.defineEntity(id, typeIdent)

	b.add(dbPartDb, dbInstallValueType, id)
}

func (b *bootstrappingPart) defineAttribute(id bootId, attrIdent symbol.Keyword, attrType bootId, attrCardinality bootId) {
	b.defineEntity(id, attrIdent)

	b.add(id, dbValueType, attrType)
	b.add(id, dbCardinality, attrCardinality)

	b.add(dbPartDb, dbInstallAttribute, id)
}

func BootstrappingPart(baseT int64) []base.Datom {
	var part bootstrappingPart

	part.definePartition(dbPartDb, DbPartDb)
	part.definePartition(dbPartTx, DbPartTx)
	part.definePartition(dbPartUser, DbPartUser)

	part.defineValueType(dbTypeBool, DbTypeBoolean)
	part.defineValueType(dbTypeFloat64, DbTypeDouble)
	part.defineValueType(dbTypeString, DbTypeString)
	part.defineValueType(dbTypeInt64, DbTypeLong)
	part.defineValueType(dbTypeRef, DbTypeRef)
	part.defineValueType(dbTypeKeyword, DbTypeKeyword)
	part.defineValueType(dbTypeTime, DbTypeInstant)

	part.defineAttribute(dbInstallValueType, DbInstallValueType, dbTypeRef, dbCardinalityMany)
	part.defineAttribute(dbInstallPartition, DbInstallPartition, dbTypeRef, dbCardinalityMany)
	part.defineAttribute(dbInstallAttribute, DbInstallAttribute, dbTypeRef, dbCardinalityMany)
	part.defineAttribute(dbIdent, DbIdent, dbTypeKeyword, dbCardinalityOne)
	part.defineAttribute(dbValueType, DbValueType, dbTypeRef, dbCardinalityOne)
	part.defineAttribute(dbCardinality, DbCardinality, dbTypeRef, dbCardinalityOne)
	part.defineAttribute(dbUnique, DbUnique, dbTypeRef, dbCardinalityOne)
	part.defineAttribute(dbIsComponent, DbIsComponent, dbTypeBool, dbCardinalityOne)
	part.defineAttribute(dbDoc, DbDoc, dbTypeString, dbCardinalityOne)
	part.defineAttribute(dbTxInstant, DbTxInstant, dbTypeTime, dbCardinalityOne)

	part.defineEntity(dbCardinalityOne, DbCardinalityOne)
	part.defineEntity(dbCardinalityMany, DbCardinalityMany)

	part.defineEntity(dbUniqueValue, DbUniqueValue)
	part.defineEntity(dbUniqueIdentity, DbUniqueIdentity)

	part.add(dbIdent, dbUnique, dbUniqueIdentity)

	return part.data
}
