package database

import (
	"reflect"

	"github.com/leidegre/datoms/cow"
	"github.com/leidegre/datoms/internal/schema"
	"github.com/leidegre/datoms/symbol"
)

type ContractType struct {
	Kind   reflect.Kind
	Elem   *ContractType   // Elem is nil for all kinds except Pointer, Slice
	Fields []ContractField // Fields are nil for all kinds except Struct
}

type ContractField struct {
	Kind  reflect.Kind
	Type  *ContractType
	Index []int
	Ident symbol.Keyword
}

func contractType(t reflect.Type) *ContractType {
	// todo: recursive type
	// todo: validate that you don't reuse ident, each ident must be unique per entity
	kind := t.Kind()
	switch kind {
	case reflect.Struct:
		// some struct types are terminal, like time.Time and symbol.Keyword
		return &ContractType{Kind: kind, Fields: contractFields(nil, t, nil)}
	case reflect.Pointer, reflect.Slice:
		return &ContractType{Kind: kind, Elem: contractType(t.Elem())}
	default:
		return &ContractType{Kind: kind}
	}
}

func contractFields(fields []ContractField, st reflect.Type, index []int) []ContractField {
	for i, end := 0, st.NumField(); i < end; i++ {
		sf := st.Field(i)
		if !sf.IsExported() {
			continue
		}
		kind := sf.Type.Kind()
		identTag := sf.Tag.Get("ident")
		if 0 < len(identTag) {
			ident := symbol.For(identTag)
			switch ident {
			case schema.DbId, schema.DbIdent:
				continue // filter out these as they have special meaning and we know they will be accessible via the Entity interface
			}
			fields = append(fields, ContractField{
				Kind:  kind,
				Type:  contractType(sf.Type),
				Index: cow.Append(index, sf.Index...),
				Ident: ident})
		} else if kind == reflect.Struct && sf.Anonymous {
			fields = contractFields(fields, sf.Type, cow.Append(index, sf.Index...))
		}
	}
	return fields
}
