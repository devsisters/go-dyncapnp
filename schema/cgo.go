package schema

/*
#cgo CXXFLAGS: -std=c++14 -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnp-json

#include <stdlib.h>
#include "cgo.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Schema

func schemaGetGeneric(schemaPtr unsafe.Pointer) unsafe.Pointer {
	return C.schemaGetGeneric(schemaPtr)
}

func schemaAsStruct(schemaPtr unsafe.Pointer) unsafe.Pointer {
	return C.schemaAsStruct(schemaPtr)
}

func schemaAsEnum(schemaPtr unsafe.Pointer) unsafe.Pointer {
	return C.schemaAsEnum(schemaPtr)
}

func schemaAsInterface(schemaPtr unsafe.Pointer) unsafe.Pointer {
	return C.schemaAsInterface(schemaPtr)
}

func schemaAsConst(schemaPtr unsafe.Pointer) unsafe.Pointer {
	return C.schemaAsConst(schemaPtr)
}

func schemaToJson(schemaPtr unsafe.Pointer) ([]byte, error) {
	return readByteArray(C.schemaToJson(schemaPtr))
}

func releaseSchema(schemaPtr unsafe.Pointer) {
	C.releaseSchema(schemaPtr)
}

// StructSchema

func structGetFields(structSchemaPtr unsafe.Pointer) []unsafe.Pointer {
	return readPointerList(C.structGetFields(structSchemaPtr))
}

func structGetUnionFields(structSchemaPtr unsafe.Pointer) []unsafe.Pointer {
	return readPointerList(C.structGetUnionFields(structSchemaPtr))
}

func structGetNonUnionFields(structSchemaPtr unsafe.Pointer) []unsafe.Pointer {
	return readPointerList(C.structGetNonUnionFields(structSchemaPtr))
}

func structFindFieldByName(structSchemaPtr unsafe.Pointer, name string) (unsafe.Pointer, error) {
	st := C.CString(name)
	defer C.free(unsafe.Pointer(st))
	return readPointerResult(C.structFindFieldByName(structSchemaPtr, st))
}

func structJsonToBinary(structSchemaPtr unsafe.Pointer, json []byte) ([]byte, error) {
	return readByteArray(C.structJsonToBinary(structSchemaPtr, (*C.char)(unsafe.Pointer(&json[0])), C.size_t(len(json))))
}

func structJsonToPacked(structSchemaPtr unsafe.Pointer, json []byte) ([]byte, error) {
	return readByteArray(C.structJsonToPacked(structSchemaPtr, (*C.char)(unsafe.Pointer(&json[0])), C.size_t(len(json))))
}

func structBinaryToJson(structSchemaPtr unsafe.Pointer, binary []byte) ([]byte, error) {
	return readByteArray(C.structBinaryToJson(structSchemaPtr, (*C.char)(unsafe.Pointer(&binary[0])), C.size_t(len(binary))))
}

func structPackedToJson(structSchemaPtr unsafe.Pointer, binary []byte) ([]byte, error) {
	return readByteArray(C.structPackedToJson(structSchemaPtr, (*C.char)(unsafe.Pointer(&binary[0])), C.size_t(len(binary))))
}

func releaseStructSchema(structSchemaPtr unsafe.Pointer) {
	C.releaseStructSchema(structSchemaPtr)
}

// StructSchema::Field

func structFieldGetContainingStruct(structFieldPtr unsafe.Pointer) unsafe.Pointer {
	return C.structFieldGetContainingStruct(structFieldPtr)
}

func structFieldGetIndex(structFieldPtr unsafe.Pointer) int {
	return int(C.structFieldGetIndex(structFieldPtr))
}

func structFieldGetType(structFieldPtr unsafe.Pointer) unsafe.Pointer {
	return C.structFieldGetType(structFieldPtr)
}

func structFieldToJson(structFieldPtr unsafe.Pointer) ([]byte, error) {
	return readByteArray(C.structFieldToJson(structFieldPtr))
}

func releaseStructSchemaField(structFieldPtr unsafe.Pointer) {
	C.releaseStructSchemaField(structFieldPtr)
}

// Type

func typeWhich(typePtr unsafe.Pointer) uint16 {
	return uint16(C.typeWhich(typePtr))
}

func typeAsStruct(typePtr unsafe.Pointer) (unsafe.Pointer, error) {
	return readPointerResult(C.typeAsStruct(typePtr))
}

func typeAsEnum(typePtr unsafe.Pointer) (unsafe.Pointer, error) {
	return readPointerResult(C.typeAsEnum(typePtr))
}

func typeAsInterface(typePtr unsafe.Pointer) (unsafe.Pointer, error) {
	return readPointerResult(C.typeAsInterface(typePtr))
}

func typeAsList(typePtr unsafe.Pointer) (unsafe.Pointer, error) {
	return readPointerResult(C.typeAsList(typePtr))
}

func releaseType(typePtr unsafe.Pointer) {
	C.releaseType(typePtr)
}

// ListSchema

func listGetElementType(listPtr unsafe.Pointer) unsafe.Pointer {
	return C.listGetElementType(listPtr)
}
func listWhichElementType(listPtr unsafe.Pointer) uint16 {
	return uint16(C.listWhichElementType(listPtr))
}
func listGetStructElementType(listPtr unsafe.Pointer) (unsafe.Pointer, error) {
	return readPointerResult(C.listGetStructElementType(listPtr))
}
func listGetEnumElementType(listPtr unsafe.Pointer) (unsafe.Pointer, error) {
	return readPointerResult(C.listGetEnumElementType(listPtr))
}
func listGetInterfaceElementType(listPtr unsafe.Pointer) (unsafe.Pointer, error) {
	return readPointerResult(C.listGetInterfaceElementType(listPtr))
}
func listGetListElementType(listPtr unsafe.Pointer) (unsafe.Pointer, error) {
	return readPointerResult(C.listGetListElementType(listPtr))
}
func releaseListSchema(listPtr unsafe.Pointer) {
	C.releaseListSchema(listPtr)
}

// EnumSchema

func enumGetEnumerants(enumPtr unsafe.Pointer) []unsafe.Pointer {
	return readPointerList(C.enumGetEnumerants(enumPtr))
}

func enumFindEnumerantByName(enumPtr unsafe.Pointer, name string) (unsafe.Pointer, error) {
	st := C.CString(name)
	defer C.free(unsafe.Pointer(st))
	return readPointerResult(C.enumFindEnumerantByName(enumPtr, st))
}

func releaseEnum(enumPtr unsafe.Pointer) {
	C.releaseEnum(enumPtr)
}

func enumerantGetOrdinal(enumerantPtr unsafe.Pointer) uint16 {
	return uint16(C.enumerantGetOrdinal(enumerantPtr))
}

func releaseEnumerant(enumerantPtr unsafe.Pointer) {
	C.releaseEnumerant(enumerantPtr)
}

// ConstSchema

func releaseConstSchema(schemaPtr unsafe.Pointer) {
	C.releaseConstSchema(schemaPtr)
}

// InterfaceSchema

func releaseInterfaceSchema(schemaPtr unsafe.Pointer) {
	C.releaseInterfaceSchema(schemaPtr)
}

// helper functions

func readPointerList(res C.pointerList) []unsafe.Pointer {
	l := (*[1 << 30]unsafe.Pointer)(unsafe.Pointer(res.ptr))[:res.len:res.len]
	ptrs := make([]unsafe.Pointer, len(l))
	copy(ptrs, l)
	C.free(unsafe.Pointer(res.ptr))
	return ptrs
}

func readPointerResult(res C.pointer_result) (unsafe.Pointer, error) {
	if res.err != nil {
		err := fmt.Errorf(C.GoString(res.err))
		C.free(unsafe.Pointer(res.err))
		return nil, err
	}
	return res.ptr, nil
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
