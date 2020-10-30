package schema

import (
	"runtime"
	"unsafe"
)

func newList(ptr unsafe.Pointer) *List {
	l := &List{
		ptr: ptr,
	}
	runtime.SetFinalizer(l, (*List).release)
	return l
}

type List struct {
	ptr unsafe.Pointer

	noCopy noCopy
}

func (l *List) Type() *Type {
	return newType(mustPtr(typeFromListSchema(l.ptr)))
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

func (l *List) release() {
	releaseListSchema(l.ptr)
	l.ptr = nil
	runtime.SetFinalizer(l, nil)
}
