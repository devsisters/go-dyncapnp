package dyncapnp

/*
#cgo CXXFLAGS: -std=c++14 -stdlib=libc++ -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnpc -lcapnp-json

#include "json.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func readByteArray(res C.struct_byteArray_result) ([]byte, error) {
	if res.err != nil {
		err := fmt.Errorf(C.GoString(res.err))
		C.free(unsafe.Pointer(res.err))
		return nil, err
	}
	bs := C.GoBytes(unsafe.Pointer(res.result.arr), C.int(res.result.length))
	C.free(unsafe.Pointer(res.result.arr))
	return bs, nil
}

func (s ParsedSchema) Encode(json []byte) ([]byte, error) {
	return readByteArray(C.jsonToBinary(s.ptr, (*C.char)(unsafe.Pointer(&json[0])), C.size_t(len(json))))
}

func (s ParsedSchema) EncodePacked(json []byte) ([]byte, error) {
	return readByteArray(C.jsonToPacked(s.ptr, (*C.char)(unsafe.Pointer(&json[0])), C.size_t(len(json))))
}

func (s ParsedSchema) Decode(bin []byte) ([]byte, error) {
	return readByteArray(C.binaryToJson(s.ptr, (*C.char)(unsafe.Pointer(&bin[0])), C.size_t(len(bin))))
}

func (s ParsedSchema) DecodePacked(bin []byte) ([]byte, error) {
	return readByteArray(C.packedToJson(s.ptr, (*C.char)(unsafe.Pointer(&bin[0])), C.size_t(len(bin))))
}
