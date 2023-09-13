package vector_test

import (
	"testing"

	"github.com/leidegre/datoms/immutable/vector"
	"github.com/leidegre/datoms/testutil"
)

func TestAppend(t *testing.T) {
	run := func(t *testing.T, n, bits int) {
		var v = vector.Make[int](bits)
		for i := 0; i < n; i++ {
			v = v.Append(i)
			t.Logf("%v", v)
		}
		v = v.Append(n)
		t.Logf("%v", v)
		for i := 0; i <= n; i++ {
			actual := v.Get(i)
			if !(i == actual) {
				t.Fatalf("expected %v, actual %v", i, actual)
			}
		}
		// t.Fail()
	}

	t.Run("2-1", func(t *testing.T) {
		run(t, 2, 1)
	})

	t.Run("6-1", func(t *testing.T) {
		run(t, 6, 1)
	})

	t.Run("10-1", func(t *testing.T) {
		run(t, 10, 1)
	})

	t.Run("100-1", func(t *testing.T) {
		run(t, 100, 1)
	})

	t.Run("2000-2", func(t *testing.T) {
		run(t, 2000, 2)
	})

	t.Run("3000-3", func(t *testing.T) {
		run(t, 3000, 3)
	})
}

func BenchmarkAppend(b *testing.B) {
	bench := func(n, bits int) {
		var v = vector.Make[int](bits)
		for i := 0; i < n; i++ {
			v = v.Append(i)
		}
	}

	b.Run("1", func(b *testing.B) {
		bench(b.N, 1)
	})

	b.Run("2", func(b *testing.B) {
		bench(b.N, 2)
	})

	b.Run("3", func(b *testing.B) {
		bench(b.N, 3)
	})

	b.Run("4", func(b *testing.B) {
		bench(b.N, 4)
	})

	b.Run("5", func(b *testing.B) {
		bench(b.N, 5)
	})

	b.Run("6", func(b *testing.B) {
		bench(b.N, 6)
	})

	b.Run("7", func(b *testing.B) {
		bench(b.N, 7)
	})
}

func TestTransientAppend(t *testing.T) {
	run := func(t *testing.T, n int) {
		var v vector.Transient[int]
		for i := 0; i < n; i++ {
			v.Append(i)
		}
		for i := 0; i < n; i++ {
			actual := v.Get(i)
			if !(i == actual) {
				t.Fatalf("expected %v, actual %v", i, actual)
			}
		}
	}
	for i := 0; i <= 64*1024; i += 33 {
		run(t, i)
	}
}

func BenchmarkTransient(b *testing.B) {
	b.Run("Persistent", func(b *testing.B) {
		var v vector.Vector[int]
		for i := 0; i < b.N; i++ {
			v = v.Append(i)
		}
	})

	b.Run("Transient", func(b *testing.B) {
		var v vector.Transient[int]
		for i := 0; i < b.N; i++ {
			v.Append(i)
		}
	})
}

func TestIterator(t *testing.T) {
	t.Run("small", func(t *testing.T) {
		var v vector.Vector[int]

		v = v.Append(1)
		v = v.Append(2)
		v = v.Append(3)

		it := v.Range(0, v.Len())
		if !it.Valid() {
			t.Fatal("!it.Valid()")
		}
		testutil.AreEqual(t, 1, it.Value())
		it.Next()
		if !it.Valid() {
			t.Fatal("!it.Valid()")
		}
		testutil.AreEqual(t, 2, it.Value())
		it.Next()
		if !it.Valid() {
			t.Fatal("!it.Valid()")
		}
		testutil.AreEqual(t, 3, it.Value())
		it.Next()
		if it.Valid() {
			t.Fatal("it.Valid()")
		}
	})

	t.Run("big", func(t *testing.T) {
		var v vector.Vector[int]

		for i := 0; i < 64; i++ {
			v = v.Append(i)
		}

		j := 0
		for it := v.Range(0, v.Len()); it.Valid(); it.Next() {
			testutil.AreEqual(t, j, it.Value())
			j++
		}
	})
}
