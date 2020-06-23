package schema

import (
	"unsafe"
)

func newInterface(ptr unsafe.Pointer) *Interface {
	return &Interface{
		Schema: NewWithFreer(ptr, releaseInterfaceSchema),
	}
}

type Interface struct {
	*Schema
}

func (s *Interface) Type() *Type {
	return newType(mustPtr(typeFromInterfaceSchema(s.ptr)))
}
