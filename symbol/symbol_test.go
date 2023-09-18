package symbol_test

import (
	"testing"

	"github.com/leidegre/datoms/symbol"
)

func TestSymbol(t *testing.T) {
	foo := symbol.New("foo")
	bar := symbol.New("foo")

	if !(foo == foo) {
		t.Fatal("symbols should equal themselves")
	}

	if !(bar == bar) {
		t.Fatal("symbols should equal themselves")
	}

	if !(foo != bar) {
		t.Fatal("symbols should be unique")
	}
}

func TestKeyword(t *testing.T) {
	a := symbol.For("foo")
	b := symbol.For("foo")
	c := symbol.For("bar")

	if !(a == b) {
		t.Fatal("keywords should equal themselves")
	}

	if !(a != c) {
		t.Fatal("different keywords should not equal themselves")
	}
}
