package dyncapnp

/*
#cgo CXXFLAGS: -std=c++14 -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnpc

#include "parser.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func parseSchemaFromFiles(files map[string][]byte, imports map[string][]byte, paths []string) (unsafe.Pointer, []unsafe.Pointer, error) {
	cFiles, freeFiles := filesToCapnpFiles(files)
	defer freeFiles()

	cImports, freeImports := filesToCapnpFiles(imports)
	defer freeImports()

	cPathSlice := make([]*C.char, len(paths))
	for i, path := range paths {
		cPathSlice[i] = C.CString(path)
	}
	defer func() {
		for _, cPath := range cPathSlice {
			C.free(unsafe.Pointer(cPath))
		}
	}()

	res := C.parseSchemaFromFiles(cFiles, C.size_t(len(files)), cImports, C.size_t(len(imports)), &cPathSlice[0], C.size_t(len(cPathSlice)))
	if res.err != nil {
		err := fmt.Errorf(C.GoString(res.err))
		C.free(unsafe.Pointer(res.err))
		return nil, nil, err
	}

	// copy array of void* into slice
	cSchemasSlice := (*[1 << 30]unsafe.Pointer)(unsafe.Pointer(res.schemas))[:len(paths):len(paths)]
	ptrs := make([]unsafe.Pointer, len(cSchemasSlice))
	copy(ptrs, cSchemasSlice)

	// free array created in cgo
	C.free(unsafe.Pointer(res.schemas))

	return res.parser, ptrs, nil
}

func filesToCapnpFiles(files map[string][]byte) (*C.struct_capnpFile, func()) {
	cFiles := (*C.struct_capnpFile)(C.malloc(C.sizeof_struct_capnpFile * C.size_t(len(files))))

	cFileSlice := (*[1 << 30]C.struct_capnpFile)(unsafe.Pointer(cFiles))[:len(files):len(files)]
	i := 0
	for path, content := range files {
		cFileSlice[i].path = C.CString(path)
		cFileSlice[i].content = (*C.char)(C.CBytes(content))
		cFileSlice[i].contentLen = C.size_t(len(content))
		i++
	}

	return cFiles, func() {
		for _, f := range cFileSlice {
			C.free(unsafe.Pointer(f.content))
			C.free(unsafe.Pointer(f.path))
		}
		C.free(unsafe.Pointer(cFiles))
	}
}

func findNested(ptr unsafe.Pointer, name string) (unsafe.Pointer, error) {
	st := C.CString(name)
	defer C.free(unsafe.Pointer(st))

	res := C.findNested(ptr, st)
	if res.err != nil {
		err := fmt.Errorf(C.GoString(res.err))
		C.free(unsafe.Pointer(res.err))
		return nil, err
	}

	return res.schema, nil
}

func releaseParsedSchema(ptr unsafe.Pointer) {
	C.releaseParsedSchema(ptr)
}

func releaseParser(ptr unsafe.Pointer) {
	C.releaseParser(ptr)
}
