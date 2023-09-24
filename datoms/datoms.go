package datoms

import (
	"github.com/leidegre/datoms/internal/base"
)

type Index = base.Index

const (
	EAVT = base.EAVT // Entity-Attribute-Value (document like)
	AEVT = base.AEVT // Attribute-Entity-Value (column like)
	AVET = base.AVET // Attribute-Value-Entity (key-value like)
	VAET = base.VAET // Value-Attribute-Entity (graph like)
)

type (
	Entid      = base.Entid      // Entid is anything that can resolve to an entity ID, like temp ID or lookup ref
	Entity     = base.Entity     // Entity can be embedded to create an entity type to be used with Pull or Transact.
	EntityLike = base.EntityLike // EntityLike is any struct that embeds Entity

	TxData = base.TxData

	TempId = base.TempId
)
