package hashmap_test

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/leidegre/datoms/hash"
	"github.com/leidegre/datoms/immutable/hashmap"
	"github.com/leidegre/datoms/testutil"
)

func TestAdd(t *testing.T) {
	var m hashmap.Persistent[string, string]
	m = m.Set("foo", hash.String("foo"), "bar")
	testutil.AreEqual(t, 1, m.Len())
}

func TestCollision(t *testing.T) {
	var (
		m hashmap.Persistent[int, string]
	)

	m = m.Set(1, 0, "foo")
	t.Log(hashmap.DebugString(m))
	m = m.Set(2, 0, "bar") // collision
	t.Log(hashmap.DebugString(m))
	m = m.Set(3, 0, "baz") // collision 2x
	t.Log(hashmap.DebugString(m))
	m = m.Set(4, hashmap.BranchingFactor, "qux") // split
	t.Log(hashmap.DebugString(m))

	// t.FailNow()

	if v, ok := m.Get(1, 0); !(ok && v == "foo") {
		t.FailNow()
	}

	if v, ok := m.Get(2, 0); !(ok && v == "bar") {
		t.FailNow()
	}

	if v, ok := m.Get(3, 0); !(ok && v == "baz") {
		t.FailNow()
	}
}

func TestDepth(t *testing.T) {
	var m hashmap.Persistent[int, int]

	for i := 0; i < hashmap.BranchingFactor*hashmap.BranchingFactor+1; i++ {
		m = m.Set(i, uint64(i), i*i)
	}

	for i := 0; i < hashmap.BranchingFactor*hashmap.BranchingFactor+1; i++ {
		if j, ok := m.Get(i, uint64(i)); !(ok && j == i*i) {
			t.FailNow()
		}
	}

	// t.Log(hashmap.DebugString(m))
	// t.FailNow()
}

func TestDelete(t *testing.T) {
	var m hashmap.Persistent[int, int]

	for i := 0; i < 16; i++ {
		m = m.Set(i, hash.Int(i), i*i)
	}

	if !(m.Len() == 16) {
		t.FailNow()
	}

	t.Log(hashmap.DebugString(m))

	for i := 15; 0 <= i; i-- {
		m = m.Delete(i, hash.Int(i))
		t.Log(hashmap.DebugString(m))
	}

	if !(m.Len() == 0) {
		t.FailNow()
	}
}

func TestSpec(t *testing.T) {
	type opCode int

	const (
		opAdd opCode = iota
		opGet
		opUpdate
		opDelete
		opMax
	)

	type keyHash struct {
		k int
		h uint64
	}

	type inst struct {
		op opCode
		k  keyHash
		v  int
	}

	type testOutcome struct {
		stream []inst
		err    string // if the test failed
	}

	pick := func(m map[keyHash]int) (keyHash, int, bool) {
		if len(m) == 0 {
			return keyHash{}, 0, false
		}
		i := rand.Intn(len(m))
		j := 0
		for k, v := range m {
			if i == j {
				return k, v, true
			}
			j++
		}
		panic(nil)
	}

	apply := func(m map[keyHash]int, inst inst) {
		switch inst.op {
		case opAdd:
			m[inst.k] = inst.v
		case opGet:
			// ...
		case opUpdate:
			m[inst.k] = inst.v
		case opDelete:
			delete(m, inst.k)
		}
	}

	newKey := func(m int) keyHash {
		k := rand.Intn(m)
		return keyHash{k, hash.Int(k)}
	}

	// newTest will generate a test case with up to N instructions
	// using the magnitude M to generate keys in between 0 and M-1
	newTest := func(n, m int) []inst {
		var (
			test []inst
			base = make(map[keyHash]int)
		)

		// seed with rand(m) number of items
		for i, end := 0, rand.Intn(m); i < end; i++ {
			k, v := newKey(m), rand.Int()
			if _, ok := base[k]; ok {
				continue
			}
			test = append(test, inst{opAdd, k, v})
			apply(base, test[len(test)-1])
		}

		// make random mutation
		for i := 0; i < n; i++ {
			switch opCode(rand.Intn(int(opMax))) {
			case opAdd:
				k, v := newKey(m), rand.Int()
				if _, ok := base[k]; ok {
					continue
				}
				test = append(test, inst{opAdd, k, v})
			case opGet:
				k, v, ok := pick(base)
				if !ok {
					continue
				}
				test = append(test, inst{opGet, k, v})
			case opUpdate:
				v := rand.Int()
				k, _, ok := pick(base)
				if !ok {
					continue
				}
				test = append(test, inst{opUpdate, k, v})
			case opDelete:
				k, v, ok := pick(base)
				if !ok {
					continue
				}
				test = append(test, inst{opDelete, k, v})
			}
			apply(base, test[len(test)-1])
		}
		return test
	}

	runTest := func(test []inst) ([]inst, string) {
		var h hashmap.Persistent[int, int]
		for i, inst := range test {
			z := h.Len()
			switch inst.op {
			case opAdd:
				h = h.Set(inst.k.k, inst.k.h, inst.v)
				if v, ok := h.Get(inst.k.k, inst.k.h); !(ok && v == inst.v) {
					return test[:i+1], "key was not found after add"
				}
				if !(z+1 == h.Len()) {
					return test[:i+1], "size should be incremented after add"
				}
			case opGet:
				if v, ok := h.Get(inst.k.k, inst.k.h); !(ok && v == inst.v) {
					return test[:i+1], "key with value was not found"
				}
			case opUpdate:
				h = h.Set(inst.k.k, inst.k.h, inst.v)
				if v, ok := h.Get(inst.k.k, inst.k.h); !(ok && v == inst.v) {
					return test[:i+1], "key was not found with new value after update"
				}
				if !(z == h.Len()) {
					return test[:i+1], "size should be unchanged after update"
				}
			case opDelete:
				h = h.Delete(inst.k.k, inst.k.h)
				if _, ok := h.Get(inst.k.k, inst.k.h); ok {
					return test[:i+1], "key was found after delete"
				}
				if !(z-1 == h.Len()) {
					return test[:i+1], "size should be decremented after delete"
				}
			}
		}
		return nil, "" // ok
	}

	concurrency := runtime.NumCPU()

	var wg sync.WaitGroup

	wg.Add(concurrency)

	var ch = make(chan testOutcome)

	for i := 0; i < concurrency; i++ {
		go func() {
			for i := 0; i < 1000; i++ {
				stream, err := runTest(newTest(100, 100))
				if err != "" {
					ch <- testOutcome{stream, err}
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var (
		minLen = math.MaxInt
		min    testOutcome
	)

	for failure := range ch {
		if len(failure.stream) < minLen {
			min = failure
		}
	}

	if 0 < len(min.stream) {
		var s string = "var h hashmap.HashMap[int, int]\n"
		for _, inst := range min.stream {
			switch inst.op {
			case opAdd:
				s += fmt.Sprintf("h = h.Set(%v, %v, %v) // add\n", inst.k.k, inst.k.h, inst.v)
			case opGet:
				// s += fmt.Sprintf("h.Get(hashmap.Int(%v))\n", inst.k)
			case opUpdate:
				s += fmt.Sprintf("h = h.Set(%v, %v, %v) // update\n", inst.k.k, inst.k.h, inst.v)
			}
		}
		lst := min.stream[len(min.stream)-1]
		s += fmt.Sprintf("if v, ok := h.Get(%v, %v); !(ok && v == %v) {\n\tt.Fatal(%#v)\n}\n", lst.k.k, lst.k.h, lst.v, min.err)
		t.Log(s)
		t.Fatal(min.err)
	}
}

func BenchmarkInsert(b *testing.B) {
	// BranchingFactor: 64
	// 	BenchmarkInsert/map[int]int-32         	              11811010	       138.4 ns/op	      60 B/op	       0 allocs/op
	// 	BenchmarkInsert/hashmap.HashMap[int,_int]-32         	 1000000	        1721 ns/op	    3051 B/op	      10 allocs/op

	// BranchingFactor: 32
	// 	BenchmarkInsert/map[int]int-32         	              12001692	       139.7 ns/op	      59 B/op	       0 allocs/op
	// 	BenchmarkInsert/hashmap.HashMap[int,_int]-32         	 1000000	        1332 ns/op	    1909 B/op	      11 allocs/op

	b.Run("map[int]int", func(b *testing.B) {
		var m = make(map[int]int)
		for i := 0; i < b.N; i++ {
			m[i] = i * i
		}
	})

	b.Run("hashmap.HashMap[int, int]", func(b *testing.B) {
		var m hashmap.Persistent[int, int]
		for i := 0; i < b.N; i++ {
			m = m.Set(i, hash.Int(i), i*i)
		}
	})
}

func TestTransient(t *testing.T) {
	const n = 64

	var m hashmap.Transient[int, int]

	for i := 0; i < n; i++ {
		m.Set(i, hash.Int(i), i*i)
	}

	t.Log(hashmap.DebugString(m.Persistent))

	if !(m.Len() == n) {
		t.FailNow()
	}

	imm := m.Immutable()

	for i := 0; i < n; i++ {
		if v, ok := imm.Get(i, hash.Int(i)); !(ok && v == i*i) {
			t.FailNow()
		}
	}
}
