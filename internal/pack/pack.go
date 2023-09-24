package pack

const (
	shift        = 42
	maxEntity    = (1 << shift) - 1
	maxPartition = (1 << (62 - shift)) - 1
	tempId       = -4611686018427387904 // 0xc000000000000000
)

func EntityId(part, ent int64) int64 {
	// Assuming each entity is just 1 datom of 25 bytes each
	// the footprint of a full database is at least ~100 TiB

	// s t pppppppppppppppppppp eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee

	if !(1 <= part && part <= maxPartition+1) {
		panic("datoms: partition ID is out of range")
	}

	if !(1 <= ent && ent <= maxEntity) {
		panic("datoms: entity ID is out of range")
	}

	// :db.part/db will have an entity ID of 1 but
	// we want this to be the zeroth partition
	// so that all the part bits are zero

	return ((part - 1) << shift) | ent
}

func TempId(part, ent int64) int64 {
	return tempId | EntityId(part, ent)
}

// Extract the partition and entity from packed
func Unpack(packed int64) (part int64, ent int64) {
	return ((packed >> shift) & maxPartition) + 1, packed & maxEntity
}
