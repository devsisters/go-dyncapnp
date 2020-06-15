package schema

/*
#cgo CXXFLAGS: -std=c++14 -stdlib=libc++ -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnp-json

#include "schema.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func newList(ptr unsafe.Pointer) *List {
	l := &List{
		ptr: ptr,
	}
	l.self = l
	runtime.SetFinalizer(l, (*List).Release)
	return l
}

type List struct {
	self *List
	ptr unsafe.Pointer
}

func (l *List) ElementType() *Type {
	return newType(C.listGetElementType(l.ptr))
}

func (l *List) ElementTypeWhich() TypeWhich {
	return TypeWhich(C.listWhichElementType(l.ptr))
}

func (l *List) StructElementType() *Struct {
	return newStruct(mustPtr(C.listGetStructElementType(l.ptr)))
}

func (l *List) EnumElementType() *Enum {
	return newEnum(mustPtr(C.listGetEnumElementType(l.ptr)))
}

func (l *List) InterfaceElementType() *Interface {
	return newInterface(mustPtr(C.listGetInterfaceElementType(l.ptr)))
}

func (l *List) ListElementType() *List {
	return newList(mustPtr(C.listGetListElementType(l.ptr)))
}

func (l *List) Release() {
	if l != l.self {
		panic("Schema should not be copied")
	}
	C.releaseListSchema(l.ptr)
	l.ptr = nil
	runtime.SetFinalizer(l, nil)
}
