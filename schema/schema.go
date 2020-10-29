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
	runtime.SetFinalizer(s, (*Schema).Release)
	return s
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type Schema struct {
	ptr  unsafe.Pointer
	free Freer

	noCopy noCopy
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

func (s *Schema) AsStruct() *Struct {
	return newStruct(schemaAsStruct(s.ptr))
}

func (s *Schema) AsEnum() *Enum {
	return newEnum(schemaAsEnum(s.ptr))
}

func (s *Schema) AsInterface() *Interface {
	return newInterface(schemaAsInterface(s.ptr))
}

func (s *Schema) AsConst() *Const {
	return newConst(schemaAsConst(s.ptr))
}

func (s *Schema) Release() {
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
