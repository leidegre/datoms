package base

import "errors"

var (
	ErrCannotResolvePartition = errors.New("cannot resolve partition")
	ErrCannotResolve          = errors.New("cannot resolve Entid")
	ErrAttributeNotFound      = errors.New("attribute not found")
)
