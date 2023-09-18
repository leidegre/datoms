package hashmap

import (
	"encoding/binary"
	"hash/maphash"
)

var (
	seed = maphash.MakeSeed()
)

// Key represents a comparable key with a precomputed hash code.
type Key[K comparable] struct {
	Key  K
	Hash uint64
}

func Int(key int) Key[int] {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(key))
	return Key[int]{key, maphash.Bytes(seed, b[:])}
}

func String(key string) Key[string] {
	return Key[string]{key, maphash.String(seed, key)}
}
