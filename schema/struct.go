package schema

import (
	"bytes"
	"encoding/json"
	"runtime"
	"unsafe"

	"github.com/WKBae/go-dyncapnp/schema/proto"
)

func newStruct(ptr unsafe.Pointer) *Struct {
	return &Struct{
		Schema: NewWithFreer(ptr, releaseStructSchema),
	}
}

type Struct struct {
	*Schema
}

func (s *Struct) Fields() []*StructField {
	list := structGetFields(s.ptr)
	res := make([]*StructField, len(list))
	for i, ptr := range list {
		res[i] = newStructField(ptr)
	}
	return res
}

func (s *Struct) UnionFields() []*StructField {
	list := structGetUnionFields(s.ptr)
	res := make([]*StructField, len(list))
	for i, ptr := range list {
		res[i] = newStructField(ptr)
	}
	return res
}

func (s *Struct) NonUnionFields() []*StructField {
	list := structGetNonUnionFields(s.ptr)
	res := make([]*StructField, len(list))
	for i, ptr := range list {
		res[i] = newStructField(ptr)
	}
	return res
}

func (s *Struct) Field(name string) (*StructField, bool) {
	fieldPtr, err := structFindFieldByName(s.ptr, name)
	if err != nil {
		panic(err)
	}
	if fieldPtr == nil {
		return nil, false
	}
	return newStructField(fieldPtr), true
}

func (s *Struct) Encode(json []byte) ([]byte, error) {
	return structJsonToBinary(s.ptr, json)
}

func (s *Struct) EncodePacked(json []byte) ([]byte, error) {
	return structJsonToPacked(s.ptr, json)
}

func (s *Struct) Decode(bin []byte) ([]byte, error) {
	return structBinaryToJson(s.ptr, bin)
}

func (s *Struct) DecodePacked(bin []byte) ([]byte, error) {
	return structPackedToJson(s.ptr, bin)
}

func newStructField(ptr unsafe.Pointer) *StructField {
	s := &StructField{
		ptr: ptr,
	}
	s.self = s
	runtime.SetFinalizer(s, (*StructField).Release)
	return s
}

type StructField struct {
	self *StructField
	ptr  unsafe.Pointer
}

func (f *StructField) Proto() (proto.Field, error) {
	b, err := structFieldToJson(f.ptr)
	if err != nil {
		return nil, err
	}
	var m proto.Field
	dec := json.NewDecoder(bytes.NewBuffer(b))
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func (f *StructField) Parent() *Struct {
	return newStruct(structFieldGetContainingStruct(f.ptr))
}

func (f *StructField) Index() int {
	return structFieldGetIndex(f.ptr)
}

func (f *StructField) Type() *Type {
	return newType(structFieldGetType(f.ptr))
}

func (f *StructField) Release() {
	if f != f.self {
		panic("Schema should not be copied")
	}
	releaseStructSchemaField(f.ptr)
	f.ptr = nil
	runtime.SetFinalizer(f, nil)
}
