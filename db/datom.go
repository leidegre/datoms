package db

type Datom struct {
	E int64
	A int64
	V interface{}
	T int64
}

// Op reports whether this is an assertion (true) or retraction (false).
func (d *Datom) Op() bool {
	return (d.T & 1) == 1
}
