package database

import (
	"reflect"

	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/internal/pack"
	"github.com/leidegre/datoms/internal/schema"
	"github.com/leidegre/datoms/symbol"
)

type entity struct {
	identities []int64
}

type txBuilder struct {
	db           Interface
	baseT, nextT int64
	schema       schema.Interface
	entities     map[int64]*entity
	tempIds      map[int64]int64
	data         []base.Datom
}

func (tx *txBuilder) init(db Interface) {
	_, nextT := db.T()

	tx.db = db
	tx.baseT = nextT
	tx.nextT = nextT
	tx.schema = db.Schema()
	tx.entities = nil
	tx.tempIds = make(map[int64]int64)
	tx.data = nil
}

func (tx *txBuilder) resolveEntid(id base.Entid) (int64, error) {
	switch id := id.(type) {
	case base.TempId:
		partId, ok := tx.schema.Id(id.Part.Ident)
		if !ok {
			return 0, base.ErrCannotResolvePartition
		}
		return pack.TempId(partId, id.TempId), nil
	default:
		return 0, base.ErrCannotResolve
	}
}

func (tx *txBuilder) resolveAttr(ident symbol.Keyword) (schema.Attr, error) {
	if attr, ok := tx.schema.AttrKeyword(ident); ok {
		return attr, nil
	}
	return schema.Attr{}, base.ErrAttributeNotFound
}

func (tx *txBuilder) emitEntidAttr(e base.Entid, a symbol.Keyword, v interface{}, op int64) (err error) {
	entid, err := tx.resolveEntid(e)
	if err != nil {
		return
	}

	attr, err := tx.resolveAttr(a)
	if err != nil {
		return
	}

	tx.emit(entid, attr, v, op)
	return
}

func (tx *txBuilder) emit(e int64, attr schema.Attr, v interface{}, op int64) (err error) {
	// todo: canonicalize value
	// todo: what if attr is identity

	var d = base.NewDatom(e, attr.Id, v, tx.baseT, op)

	tx.data = append(tx.data, d)
	return
}

func (tx *txBuilder) txExpand(entid base.Entid, e base.EntityLike) (id int64, err error) {
	if e == nil {
		panic("datoms: a top-level transaction map cannot be nil") // or do we silently ignore this?
	}

	if entid == nil {
		entid = e
	}

	if entid.Zero() {
		entid = base.NewTempId(schema.DbPartUser)
	}

	id, err = tx.resolveEntid(entid)
	if err != nil {
		return 0, err
	}

	v := reflect.ValueOf(e)
	t := contractType(v.Type())

	for _, tf := range t.Fields {
		vf := v.FieldByIndex(tf.Index)
		switch tf.Kind {
		case reflect.Pointer, reflect.Slice:
			if vf.IsNil() {
				continue
			}
		case reflect.String:
			if len(vf.String()) == 0 {
				continue
			}
		}
		attr, err := tx.resolveAttr(tf.Ident)
		if err != nil {
			return 0, err
		}
		tx.emit(id, attr, vf.Interface(), 1)
	}

	return id, nil
}
