package base_test

import (
	"testing"

	"github.com/leidegre/datoms/internal/base"
	"github.com/leidegre/datoms/symbol"
	"github.com/leidegre/datoms/testutil"
)

func TestEntity(t *testing.T) {
	// The only way to make an EntityLike is to embed the Entity struct

	type test struct {
		base.Entity
	}

	t.Run("Id", func(t *testing.T) {
		const (
			expected int64 = 1234
		)

		// Internally we use the EntityLike interface to access Entity identities Id and Ident
		// Externally this isn't necessary because you have access the the fields via the type

		var entity = base.Entity{Id: expected}
		var like base.EntityLike = &test{Entity: entity}
		var actual, _ = base.EntityIdentities(like)

		testutil.AreEqual(t, expected, actual)
	})

	t.Run("Ident", func(t *testing.T) {
		var (
			expected = symbol.For(":test")
		)

		var entity = base.Entity{Ident: expected}
		var like base.EntityLike = &test{Entity: entity}
		var _, actual = base.EntityIdentities(like)

		testutil.AreEqual(t, expected, actual)
	})

	t.Run("Entid", func(t *testing.T) {
		var (
			a base.Entid = test{}
			b base.Entid = test{}
		)

		testutil.AreEqual(t, a, b)
	})

	t.Run("Entid (value equality)", func(t *testing.T) {
		// Here we demonstrate that we can use the value as a unique identity

		var (
			a base.Entid = test{Entity: base.Entity{Id: 1234}}
			b base.Entid = test{Entity: base.Entity{Id: 1234}}
		)

		testutil.AreEqual(t, a, b)

		m := make(map[base.Entid]bool)

		m[a] = true
		m[b] = true

		testutil.AreEqual(t, 1, len(m))
	})

	t.Run("Entid (pointer equality)", func(t *testing.T) {
		// Here we demonstrate that we can use the address as a distinct identity

		var (
			a base.Entid = &test{Entity: base.Entity{Id: 1234}}
			b base.Entid = &test{Entity: base.Entity{Id: 1234}}
		)

		testutil.NotEqual(t, a, b)

		m := make(map[base.Entid]bool)

		m[a] = true
		m[b] = true

		testutil.AreEqual(t, 2, len(m))
	})
}
