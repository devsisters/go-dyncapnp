package dyncapnp

/*
#cgo CXXFLAGS: -std=c++14 -stdlib=libc++ -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnpc -lcapnp-json

#include "json.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

func readByteArray(arr C.struct_byteArray) []byte {
	bs := C.GoBytes(unsafe.Pointer(arr.arr), C.int(arr.length))
	C.free(unsafe.Pointer(arr.arr))
	return bs
}

func (s Schema) Encode(json []byte) []byte {
	return readByteArray(C.jsonToBinary(s.ptr, (*C.char)(unsafe.Pointer(&json[0])), C.size_t(len(json))))
}

func (s Schema) EncodePacked(json []byte) []byte {
	return readByteArray(C.jsonToPacked(s.ptr, (*C.char)(unsafe.Pointer(&json[0])), C.size_t(len(json))))
}

func (s Schema) Decode(bin []byte) []byte {
	return readByteArray(C.binaryToJson(s.ptr, (*C.char)(unsafe.Pointer(&bin[0])), C.size_t(len(bin))))
}

func (s Schema) DecodePacked(bin []byte) []byte {
	return readByteArray(C.packedToJson(s.ptr, (*C.char)(unsafe.Pointer(&bin[0])), C.size_t(len(bin))))
}
