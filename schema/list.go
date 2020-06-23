package schema

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
	ptr  unsafe.Pointer
}

func (s *List) Type() *Type {
	return newType(mustPtr(typeFromListSchema(s.ptr)))
}

func (l *List) ElementType() *Type {
	return newType(listGetElementType(l.ptr))
}

func (l *List) ElementTypeWhich() TypeWhich {
	return TypeWhich(listWhichElementType(l.ptr))
}

func (l *List) ElementTypeAsStruct() *Struct {
	return newStruct(mustPtr(listGetStructElementType(l.ptr)))
}

func (l *List) ElementTypeAsEnum() *Enum {
	return newEnum(mustPtr(listGetEnumElementType(l.ptr)))
}

func (l *List) ElementTypeAsInterface() *Interface {
	return newInterface(mustPtr(listGetInterfaceElementType(l.ptr)))
}

func (l *List) ElementTypeAsList() *List {
	return newList(mustPtr(listGetListElementType(l.ptr)))
}

func (l *List) Release() {
	if l != l.self {
		panic("Schema should not be copied")
	}
	releaseListSchema(l.ptr)
	l.ptr = nil
	runtime.SetFinalizer(l, nil)
}
