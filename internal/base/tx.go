package base

import "github.com/leidegre/datoms/symbol"

type TxData interface {
	txData()
}

type TxAdd struct {
	E Entid
	A symbol.Keyword
	V interface{}
}

func (TxAdd) txData() {}

type TxRetract struct {
	E Entid
	A symbol.Keyword
	V interface{}
}

func (TxRetract) txData() {}

type TxMap struct {
	Id     Entid
	Entity EntityLike
}

func (TxMap) txData() {}
