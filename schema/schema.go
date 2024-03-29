package schema

import (
	"bytes"
	"encoding/json"
	"runtime"
	"unsafe"

	"github.com/devsisters/go-dyncapnp/schema/proto"
)

type Freer func(ptr unsafe.Pointer)

func New(ptr unsafe.Pointer) *Schema {
	return NewWithFreer(ptr, releaseSchema)
}

func NewWithFreer(ptr unsafe.Pointer, free Freer) *Schema {
	s := &Schema{
		ptr:  ptr,
		free: free,
	}
	if free != nil {
		runtime.SetFinalizer(s, (*Schema).release)
	}
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

func (s *Schema) ShortDisplayName() string {
	return schemaGetShortDisplayName(s.ptr)
}

func (s *Schema) release() {
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
