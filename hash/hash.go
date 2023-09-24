package hash

import (
	"encoding/binary"
	"hash/maphash"
)

var (
	seed = maphash.MakeSeed()
)

func Int(v int) uint64 {
	return Uint64(uint64(v))
}

func Bytes(b []byte) uint64 {
	return maphash.Bytes(seed, b)
}

func String(s string) uint64 {
	return maphash.String(seed, s)
}

func Uint16(v uint16) uint64 {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], v)
	return Bytes(b[:])
}

func Uint32(v uint32) uint64 {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	return Bytes(b[:])
}

func Uint64(v uint64) uint64 {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	return Bytes(b[:])
}
