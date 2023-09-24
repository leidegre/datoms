package schema_test

import (
	"testing"

	"github.com/leidegre/datoms/internal/schema"
	"github.com/leidegre/datoms/testutil"
)

func TestNew(t *testing.T) {
	schema.New()
}

func TestSchemaBootstrappingPart(t *testing.T) {
	s := schema.New()

	var (
		ident, _          = s.Id(schema.DbIdent)
		valueType, _      = s.Id(schema.DbValueType)
		cardinality, _    = s.Id(schema.DbCardinality)
		doc, _            = s.Id(schema.DbDoc)
		typeKeyword, _    = s.Id(schema.DbTypeKeyword)
		typeRef, _        = s.Id(schema.DbTypeRef)
		typeString, _     = s.Id(schema.DbTypeString)
		cardinalityOne, _ = s.Id(schema.DbCardinalityOne)
	)

	if attr, ok := s.AttrKeyword(schema.DbIdent); ok {
		testutil.AreEqual(t, ident, attr.Id)
		testutil.AreEqual(t, schema.DbIdent, attr.Ident)
		testutil.AreEqual(t, typeKeyword, attr.ValueType)
		testutil.AreEqual(t, cardinalityOne, attr.Cardinality)
	} else {
		t.FailNow()
	}

	if attr, ok := s.AttrKeyword(schema.DbValueType); ok {
		testutil.AreEqual(t, valueType, attr.Id)
		testutil.AreEqual(t, schema.DbValueType, attr.Ident)
		testutil.AreEqual(t, typeRef, attr.ValueType)
		testutil.AreEqual(t, cardinalityOne, attr.Cardinality)
	} else {
		t.FailNow()
	}

	if attr, ok := s.AttrKeyword(schema.DbCardinality); ok {
		testutil.AreEqual(t, cardinality, attr.Id)
		testutil.AreEqual(t, schema.DbCardinality, attr.Ident)
		testutil.AreEqual(t, typeRef, attr.ValueType)
		testutil.AreEqual(t, cardinalityOne, attr.Cardinality)
	} else {
		t.FailNow()
	}

	if attr, ok := s.AttrKeyword(schema.DbDoc); ok {
		testutil.AreEqual(t, doc, attr.Id)
		testutil.AreEqual(t, schema.DbDoc, attr.Ident)
		testutil.AreEqual(t, typeString, attr.ValueType)
		testutil.AreEqual(t, cardinalityOne, attr.Cardinality)
	} else {
		t.Fatal("cannot find attribute :db/doc")
	}
}
