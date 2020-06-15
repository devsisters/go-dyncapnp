package schema

/*
#cgo CXXFLAGS: -std=c++14 -stdlib=libc++ -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnp-json

#include <stdlib.h>
#include "schema.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func newEnum(ptr unsafe.Pointer) *Enum {
	return &Enum{
		Schema: NewWithFreer(ptr, releaseEnum),
	}
}

func releaseEnum(ptr unsafe.Pointer) {
	C.releaseEnum(ptr)
}

type Enum struct {
	*Schema
}

func (e *Enum) Enumerants() []*Enumerant {
	list := C.enumGetEnumerants(e.ptr)
	l := (*[1 << 30]unsafe.Pointer)(unsafe.Pointer(list.ptr))[:list.len:list.len]
	res := make([]*Enumerant, len(l))
	for i, ptr := range l {
		res[i] = newEnumerant(ptr)
	}
	C.free(unsafe.Pointer(list.ptr))
	return res
}

func (e *Enum) FindByName(name string) (*Enumerant, bool) {
	st := C.CString(name)
	defer C.free(unsafe.Pointer(st))

	res := C.enumFindEnumerantByName(e.ptr, st)
	if res.err != nil {
		msg := C.GoString(res.err)
		C.free(unsafe.Pointer(res.err))
		panic(msg)
	}
	if res.ptr == nil {
		return nil, false
	}

	return newEnumerant(res.ptr), true
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
	return uint16(C.enumerantGetOrdinal(e.ptr))
}

func (e *Enumerant) Release() {
	if e != e.self {
		panic("Schema should not be copied")
	}
	C.releaseEnumerant(e.ptr)
	e.ptr = nil
	runtime.SetFinalizer(e, nil)
}
