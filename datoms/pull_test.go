package datoms_test

import (
	"testing"
)

func TestPull(t *testing.T) {
	type Person struct {
		Id        int64  `kw:":db/id"`
		FirstName string `kw:":person/firstName"`
		LastName  string `kw:":person/lastName"`
	}

	// the properties of Person are directly accessible through Foo
	// while Person has an entity ID it will be the same as Foo
	type Foo struct {
		Id     int64 `kw:":db/id"`
		Person       // embedding
	}

	// Here person is used as a contract for a reference to another entity
	type Bar struct {
		Id       int64  `kw:":db/id"`
		Relative Person `kw:":person/relative"` // ref
	}
}
