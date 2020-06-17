package dyncapnp

/*
#cgo CXXFLAGS: -std=c++14 -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnpc

#include "parser.h"
#include <stdlib.h>

struct capnpFile* allocFiles(size_t s) {
	return (struct capnpFile*) malloc(sizeof(struct capnpFile) * s);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func parseSchemaFromFiles(files map[string][]byte, imports map[string][]byte, paths []string) (unsafe.Pointer, []unsafe.Pointer, error) {
	cFiles := C.allocFiles(C.size_t(len(files)))
	defer C.free(unsafe.Pointer(cFiles))

	cFileSlice := (*[1 << 30]C.struct_capnpFile)(unsafe.Pointer(cFiles))[:len(files):len(files)]
	i := 0
	for path, content := range files {
		cFileSlice[i].path = C.CString(path)
		cFileSlice[i].content = (*C.char)(unsafe.Pointer(&content[0]))
		cFileSlice[i].contentLen = C.size_t(len(content))
		i++
	}
	defer func() {
		for _, f := range cFileSlice {
			C.free(unsafe.Pointer(f.path))
		}
	}()

	cImports := C.allocFiles(C.size_t(len(imports)))
	defer C.free(unsafe.Pointer(cImports))

	cImportSlice := (*[1 << 30]C.struct_capnpFile)(unsafe.Pointer(cImports))[:len(imports):len(imports)]
	i = 0
	for path, content := range imports {
		cImportSlice[i].path = C.CString(path)
		cImportSlice[i].content = (*C.char)(unsafe.Pointer(&content[0]))
		cImportSlice[i].contentLen = C.size_t(len(content))
		i++
	}
	defer func() {
		for _, f := range cImportSlice {
			C.free(unsafe.Pointer(f.path))
		}
	}()

	cPathSlice := make([]*C.char, len(paths))
	for i, path := range paths {
		cPathSlice[i] = C.CString(path)
	}
	defer func() {
		for _, chs := range cPathSlice {
			C.free(unsafe.Pointer(chs))
		}
	}()

	res := C.parseSchemaFromFiles(cFiles, C.size_t(len(cFileSlice)), cImports, C.size_t(len(cImportSlice)), &cPathSlice[0], C.size_t(len(cPathSlice)))
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
