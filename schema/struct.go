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

func newStruct(ptr unsafe.Pointer) *Struct {
	return &Struct{
		Schema: NewWithFreer(ptr, releaseStructSchema),
	}
}

func releaseStructSchema(ptr unsafe.Pointer) {
	C.releaseStructSchema(ptr)
}

type Struct struct {
	*Schema
}

func (s *Struct) Fields() []*StructField {
	list := C.structGetFields(s.ptr)
	l := (*[1 << 30]unsafe.Pointer)(unsafe.Pointer(list.ptr))[:list.len:list.len]
	res := make([]*StructField, len(l))
	for i, ptr := range l {
		res[i] = newStructField(ptr)
	}
	C.free(unsafe.Pointer(list.ptr))
	return res
}

func (s *Struct) Encode(json []byte) ([]byte, error) {
	return readByteArray(C.structJsonToBinary(s.ptr, (*C.char)(unsafe.Pointer(&json[0])), C.size_t(len(json))))
}

func (s *Struct) EncodePacked(json []byte) ([]byte, error) {
	return readByteArray(C.structJsonToPacked(s.ptr, (*C.char)(unsafe.Pointer(&json[0])), C.size_t(len(json))))
}

func (s *Struct) Decode(bin []byte) ([]byte, error) {
	return readByteArray(C.structBinaryToJson(s.ptr, (*C.char)(unsafe.Pointer(&bin[0])), C.size_t(len(bin))))
}

func (s *Struct) DecodePacked(bin []byte) ([]byte, error) {
	return readByteArray(C.structPackedToJson(s.ptr, (*C.char)(unsafe.Pointer(&bin[0])), C.size_t(len(bin))))
}

func newStructField(ptr unsafe.Pointer) *StructField {
	s := &StructField{
		ptr: ptr,
	}
	s.self = s
	runtime.SetFinalizer(s, (*StructField).Release)
	return s
}

type StructField struct{
	self *StructField
	ptr unsafe.Pointer
}

func (f *StructField) Parent() *Struct {
	return newStruct(C.structFieldGetContainingStruct(f.ptr))
}

func (f *StructField) Index() int {
	return int(C.structFieldGetIndex(f.ptr))
}

func (f *StructField) Type() *Type {
	return newType(C.structFieldGetType(f.ptr))
}

func (f *StructField) Release() {
	if f != f.self {
		panic("Schema should not be copied")
	}
	C.releaseStructSchemaField(f.ptr)
	f.ptr = nil
	runtime.SetFinalizer(f, nil)
}
