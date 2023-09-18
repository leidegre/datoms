package symbol

import (
	"sync"
	"sync/atomic"
)

var (
	atom uint64
)

// Symbol is a unique atom. The name associated with a symbol doesn't mean anything.
//
//	foo := symbol.New("foo") // create like this
type Symbol struct {
	name string
	atom uint64
}

func (s Symbol) String() string { return s.name }

// Create a new symbol. Every symbol is unique. Here the name passed is just for user friendliness.
//
//	a := symbol.New("foo")
//	b := symbol.New("foo") // a != b
func New(name string) Symbol {
	return Symbol{name, atomic.AddUint64(&atom, 1)}
}

// Keyword is a distinct symbol type.
//
//	foo := symbol.For(":foo") // create like this
type Keyword struct {
	Symbol
}

var (
	kwTblLock sync.RWMutex
	kwTbl     = make(map[string]Keyword)
)

func getKeyword(name string) (kw Keyword, ok bool) {
	kwTblLock.RLock()
	defer kwTblLock.RUnlock()
	kw, ok = kwTbl[name]
	return
}

func newKeyword(name string) Keyword {
	kwTblLock.Lock()
	defer kwTblLock.Unlock()
	// Check for existing...
	if kw, ok := kwTbl[name]; ok {
		return kw
	}
	kw := Keyword{New(name)}
	kwTbl[name] = kw
	return kw
}

// Create a keyword. Safe for concurrent use by multiple goroutines.
//
//	var (
//	  foo1 = symbol.For(":foo")
//	  foo2 = symbol.For(":foo") // foo1 == foo2
//	  bar  = symbol.For(":bar") // foo1 != bar && foo2 != bar
//	)
func For(name string) Keyword {
	if kw, ok := getKeyword(name); ok {
		return kw
	}
	return newKeyword(name)
}
