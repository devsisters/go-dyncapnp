package schema

import (
	"runtime"
	"unsafe"
)

func newEnum(ptr unsafe.Pointer) *Enum {
	return &Enum{
		Schema: NewWithFreer(ptr, releaseEnum),
	}
}

type Enum struct {
	*Schema
}

func (e *Enum) Type() *Type {
	return newType(mustPtr(typeFromEnumSchema(e.ptr)))
}

func (e *Enum) Enumerants() []*Enumerant {
	list := enumGetEnumerants(e.ptr)
	res := make([]*Enumerant, len(list))
	for i, ptr := range list {
		res[i] = newEnumerant(ptr)
	}
	return res
}

func (e *Enum) FindByName(name string) (*Enumerant, bool) {
	ptr, err := enumFindEnumerantByName(e.ptr, name)
	if err != nil {
		panic(err)
	}
	if ptr == nil {
		return nil, false
	}

	return newEnumerant(ptr), true
}

func newEnumerant(ptr unsafe.Pointer) *Enumerant {
	s := &Enumerant{
		ptr: ptr,
	}
	runtime.SetFinalizer(s, (*Enumerant).release)
	return s
}

type Enumerant struct {
	ptr unsafe.Pointer

	noCopy noCopy
}

func (e *Enumerant) Ordinal() uint16 {
	return enumerantGetOrdinal(e.ptr)
}

func (e *Enumerant) release() {
	releaseEnumerant(e.ptr)
	e.ptr = nil
	runtime.SetFinalizer(e, nil)
}
