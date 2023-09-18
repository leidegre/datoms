package db

import "github.com/leidegre/datoms/symbol"

var (
	Id    = symbol.For(":db/id")
	Ident = symbol.For(":db/ident")
)

type Attribute struct {
}

type TxData interface{}

type EntityId int64

type Index int

const (
	EAVT Index = iota
	AEVT
	AVET
	VAET
)

type Seq interface {
	//...
}

type Database interface {
	T() (baseT, nextT int64)
	Schema() Schema
	SeekDatoms(index Index, components ...[]interface{}) Seq
	With(txData []TxData) Database
}

type Log interface {
	//...
}

type Transaction struct {
	DbBefore Database
	DbAfter  Database
	TxData   []Datom
	TempId   map[int64]int64
}

type Connection interface {
	Log() Log
	Database() Database
	Transact(txData []TxData) Transaction
}
