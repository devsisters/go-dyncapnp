package schema

/*
#cgo CXXFLAGS: -std=c++14 -stdlib=libc++ -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnp-json

#include <stdlib.h>
#include "schema.h"
*/
import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"unsafe"

	"github.com/WKBae/go-dyncapnp/schema/proto"
)

type Freer func(ptr unsafe.Pointer)

func NoFree(_ unsafe.Pointer) {}

func releaseSchema(ptr unsafe.Pointer) {
	C.releaseSchema(ptr)
}

func New(ptr unsafe.Pointer) *Schema {
	return NewWithFreer(ptr, releaseSchema)
}

func NewWithFreer(ptr unsafe.Pointer, free Freer) *Schema {
	s := &Schema{
		ptr: ptr,
		free: free,
	}
	s.self = s
	runtime.SetFinalizer(s, (*Schema).Release)
	return s
}

type Schema struct {
	self *Schema
	ptr  unsafe.Pointer
	free Freer
}

func (s *Schema) Proto() (proto.Proto, error) {
	b, err := readByteArray(C.schemaToJson(s.ptr))
	if err != nil {
		return nil, err
	}
	var m proto.Proto
	dec := json.NewDecoder(bytes.NewBuffer(b))
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Schema) Generic() *Schema {
	return New(C.schemaGetGeneric(s.ptr))
}

func (s *Schema) Struct() *Struct {
	return newStruct(C.schemaAsStruct(s.ptr))
}

func (s *Schema) Enum() *Enum {
	return newEnum(C.schemaAsEnum(s.ptr))
}

func (s *Schema) Interface() *Interface {
	return newInterface(C.schemaAsInterface(s.ptr))
}

func (s *Schema) Const() *Const {
	return newConst(C.schemaAsConst(s.ptr))
}

func (s *Schema) Release() {
	if s != s.self {
		panic("Schema should not be copied")
	}

	s.free(s.ptr)
	s.ptr = nil
	runtime.SetFinalizer(s, nil)
}

func mustPtr(res C.pointer_result) unsafe.Pointer {
	if res.err != nil {
		msg := C.GoString(res.err)
		C.free(unsafe.Pointer(res.err))
		panic(msg)
	}
	return res.ptr
}

func readByteArray(res C.byteArray_result) ([]byte, error) {
	if res.err != nil {
		err := fmt.Errorf(C.GoString(res.err))
		C.free(unsafe.Pointer(res.err))
		return nil, err
	}
	bs := C.GoBytes(unsafe.Pointer(res.result.arr), C.int(res.result.length))
	C.free(unsafe.Pointer(res.result.arr))
	return bs, nil
}
