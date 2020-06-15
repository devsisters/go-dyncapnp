package schema

import (
	"runtime"
	"unsafe"
)

func newType(ptr unsafe.Pointer) *Type {
	t := &Type{
		ptr: ptr,
	}
	t.self = t
	runtime.SetFinalizer(t, (*Type).Release)
	return t
}

type Type struct {
	self *Type
	ptr  unsafe.Pointer
}

func (t *Type) Which() TypeWhich {
	return TypeWhich(typeWhich(t.ptr))
}

func (t *Type) Struct() *Struct {
	return newStruct(mustPtr(typeAsStruct(t.ptr)))
}

func (t *Type) Enum() *Enum {
	return newEnum(mustPtr(typeAsEnum(t.ptr)))
}

func (t *Type) List() *List {
	return newList(mustPtr(typeAsList(t.ptr)))
}

func (t *Type) Interface() *Interface {
	return newInterface(mustPtr(typeAsInterface(t.ptr)))
}

func (t *Type) Release() {
	if t != t.self {
		panic("Schema should not be copied")
	}
	releaseType(t.ptr)
	t.ptr = nil
	runtime.SetFinalizer(t, nil)
}

type TypeWhich uint16

const (
	TypeVoid TypeWhich = iota
	TypeBool
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
	TypeFloat32
	TypeFloat64
	TypeText
	TypeData
	TypeList
	TypeEnum
	TypeStruct
	TypeInterface
	TypeAnyPointer
)

func (t TypeWhich) String() string {
	switch t {
	case TypeVoid:
		return "void"
	case TypeBool:
		return "bool"
	case TypeInt8:
		return "int8"
	case TypeInt16:
		return "int16"
	case TypeInt32:
		return "int32"
	case TypeInt64:
		return "int64"
	case TypeUint8:
		return "uint8"
	case TypeUint16:
		return "uint16"
	case TypeUint32:
		return "uint32"
	case TypeUint64:
		return "uint64"
	case TypeFloat32:
		return "float32"
	case TypeFloat64:
		return "float64"
	case TypeText:
		return "text"
	case TypeData:
		return "data"
	case TypeList:
		return "list"
	case TypeEnum:
		return "enum"
	case TypeStruct:
		return "struct"
	case TypeInterface:
		return "interface"
	case TypeAnyPointer:
		return "anyPointer"
	default:
		return ""
	}
}
