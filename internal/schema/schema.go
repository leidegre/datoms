package schema

import (
	"slices"

	"github.com/leidegre/datoms/cow"
	"github.com/leidegre/datoms/hash"
	hamt "github.com/leidegre/datoms/immutable/hashmap"
	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/internal/sort"
	"github.com/leidegre/datoms/symbol"
)

type Attr struct {
	Id          int64          `ident:":db/id"`
	Ident       symbol.Keyword `ident:":db/ident"`
	ValueType   int64          `ident:":db/valueType"`
	Cardinality int64          `ident:":db/cardinality"`
	Unique      int64          `ident:":db/unique"`
	IsComponent bool           `ident:":db/isComponent"`
	Doc         string         `ident:":db/doc"`
}

type Interface interface {
	Id(ident symbol.Keyword) (id int64, ok bool)
	Attr(attrId int64) (attr Attr, ok bool)
	AttrKeyword(attrIdent symbol.Keyword) (attr Attr, ok bool)
}

type Schema struct {
	idents hamt.Persistent[symbol.Keyword, int64]
	attrs  hamt.Persistent[int64, Attr]
}

func (s *Schema) Id(ident symbol.Keyword) (id int64, ok bool) {
	id, ok = s.idents.Get(ident, ident.Hash())
	return
}

// Attr looks up attribute by entid
func (s *Schema) Attr(attrId int64) (attr Attr, ok bool) {
	attr, ok = s.attrs.Get(attrId, hash.Uint64(uint64(attrId)))
	return
}

// AttrKeyword looks up attribute by ident
func (s *Schema) AttrKeyword(attrIdent symbol.Keyword) (attr Attr, ok bool) {
	var attrId int64
	if attrId, ok = s.Id(attrIdent); ok {
		attr, ok = s.attrs.Get(attrId, hash.Uint64(uint64(attrId)))
	}
	return
}

func (s *Schema) attrId(ident symbol.Keyword) uint32 {
	if id, ok := s.Id(ident); ok {
		return uint32(id)
	}
	panic("datoms: cannot find required attribute")
}

func (s *Schema) With(data []base.Datom) *Schema {
	data = cow.ShallowCopy(data)

	slices.SortFunc(data, sort.CompareIndex(base.EAVT))

	var (
		idents = s.idents
		attrs  = s.attrs
	)

	type elem struct {
		Id          int64
		Ident       symbol.Keyword
		ValueType   int64
		Cardinality int64
		Unique      int64
		IsComponent bool
	}

	var (
		el      elem
		els     []elem
		install map[int64]bootId
	)

	for _, d := range data {
		if el.Id != d.E {
			if el.Id != 0 {
				els = append(els, el)
			}
			el = elem{Id: d.E}
		}
		// We have to use bootstrapping partition IDs statically to get going
		switch bootId(d.A) {
		case dbInstallAttribute:
			// E will be :datoms.part/db
			// A will be :datoms.install/attribute
			// V will be entid of attribute
			if install == nil {
				install = make(map[int64]bootId)
			}
			install[d.V.(int64)] = dbInstallAttribute
		case dbIdent:
			ident := d.V.(symbol.Keyword)
			el.Ident = ident
			idents = idents.Set(ident, ident.Hash(), d.E)
		case dbValueType:
			el.ValueType = d.V.(int64)
		case dbCardinality:
			el.Cardinality = d.V.(int64)
		case dbUnique:
			el.Unique = d.V.(int64)
		case dbIsComponent:
			el.IsComponent = d.V.(bool)
		}
	}
	if el.Id != 0 {
		els = append(els, el)
	}

	if install != nil {
		for _, el := range els {
			k, h := el.Id, hash.Uint64(uint64(el.Id))
			if v, ok := install[k]; ok {
				switch v {
				case dbInstallAttribute:
					var attr Attr
					if attr, ok = s.attrs.Get(k, h); !ok {
						attr.Id = k
					}
					attr.Ident = el.Ident
					attr.ValueType = el.ValueType
					attr.Cardinality = el.Cardinality
					attr.Unique = el.Unique
					attr.IsComponent = el.IsComponent
					attrs = attrs.Set(k, h, attr)
				}
			}
		}
	}

	if s.idents != idents || s.attrs != attrs {
		return &Schema{idents: idents, attrs: attrs}
	}

	return s
}

func New() *Schema {
	return (&Schema{}).With(BootstrappingPart(0))
}
