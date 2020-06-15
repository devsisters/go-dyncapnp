package schema

import (
	"bytes"
	"encoding/json"
	"runtime"
	"unsafe"

	"github.com/WKBae/go-dyncapnp/schema/proto"
)

type Freer func(ptr unsafe.Pointer)

func NoFree(_ unsafe.Pointer) {}

func New(ptr unsafe.Pointer) *Schema {
	return NewWithFreer(ptr, releaseSchema)
}

func NewWithFreer(ptr unsafe.Pointer, free Freer) *Schema {
	s := &Schema{
		ptr:  ptr,
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
	b, err := schemaToJson(s.ptr)
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
	return New(schemaGetGeneric(s.ptr))
}

func (s *Schema) Struct() *Struct {
	return newStruct(schemaAsStruct(s.ptr))
}

func (s *Schema) Enum() *Enum {
	return newEnum(schemaAsEnum(s.ptr))
}

func (s *Schema) Interface() *Interface {
	return newInterface(schemaAsInterface(s.ptr))
}

func (s *Schema) Const() *Const {
	return newConst(schemaAsConst(s.ptr))
}

func (s *Schema) Release() {
	if s != s.self {
		panic("Schema should not be copied")
	}

	s.free(s.ptr)
	s.ptr = nil
	runtime.SetFinalizer(s, nil)
}

func mustPtr(ptr unsafe.Pointer, err error) unsafe.Pointer {
	if err != nil {
		panic(err)
	}
	return ptr
}
