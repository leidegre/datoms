package db

import "github.com/leidegre/datoms/immutable/vector"

type SimpleSchema struct {
	//...
}

type SimpleDatabase struct {
	eavt   vector.Persistent[Datom]
	aevt   vector.Persistent[Datom]
	avet   vector.Persistent[Datom]
	vaet   vector.Persistent[Datom]
	schema Schema
}

type SimpleLog struct {
	//...
}

type SimpleConnection struct {
	log SimpleLog
	db  SimpleDatabase
}
