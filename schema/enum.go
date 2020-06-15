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
	s.self = s
	runtime.SetFinalizer(s, (*Enumerant).Release)
	return s
}

type Enumerant struct {
	self *Enumerant
	ptr  unsafe.Pointer
}

func (e *Enumerant) Ordinal() uint16 {
	return enumerantGetOrdinal(e.ptr)
}

func (e *Enumerant) Release() {
	if e != e.self {
		panic("Schema should not be copied")
	}
	releaseEnumerant(e.ptr)
	e.ptr = nil
	runtime.SetFinalizer(e, nil)
}
