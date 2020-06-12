package dyncapnp

//go:generate go run ./gen/

/*
#cgo CXXFLAGS: -std=c++14 -stdlib=libc++ -I${SRCDIR}/capnproto/c++/src
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
	"runtime"
	"unsafe"
)

var ErrSchemaNotFound = fmt.Errorf("schema not found")

// ParsedSchema of a Cap'n'proto type. MUST not be copied. Should be .Release()'ed after use.
type ParsedSchema struct {
	ptr unsafe.Pointer
}

// FindNested finds nested schema with given name. Returns ErrSchemaNotFound if nothing was found.
func (s *ParsedSchema) Nested(name string) (*ParsedSchema, error) {
	st := C.CString(name)
	defer C.free(unsafe.Pointer(st))

	res := C.findNested(s.ptr, st)
	if res.err != nil {
		err := fmt.Errorf(C.GoString(res.err))
		C.free(unsafe.Pointer(res.err))
		return nil, err
	}
	if res.schema == nil {
		return nil, ErrSchemaNotFound
	}

	sc := &ParsedSchema{
		ptr: res.schema,
	}
	runtime.SetFinalizer(sc, (*ParsedSchema).Release)
	return sc, nil
}

func (s *ParsedSchema) Release() {
	if s.ptr == nil {
		return
	}
	C.releaseSchema(s.ptr)
	s.ptr = nil
	runtime.SetFinalizer(s, nil)
}

func ParseFromFiles(files map[string][]byte, imports map[string][]byte, paths []string) (map[string]*ParsedSchema, error) {
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

	importLen := len(stdImports) + len(imports)
	cImports := C.allocFiles(C.size_t(importLen))
	defer C.free(unsafe.Pointer(cImports))

	cImportSlice := (*[1 << 30]C.struct_capnpFile)(unsafe.Pointer(cImports))[:importLen:importLen]
	i = 0
	for path, content := range stdImports {
		cImportSlice[i].path = C.CString(path)
		cImportSlice[i].content = (*C.char)(unsafe.Pointer(&content[0]))
		cImportSlice[i].contentLen = C.size_t(len(content))
		i++
	}
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
		return nil, err
	}

	// iterate void** array and wrap them with ParsedSchema
	cSchemasSlice := (*[1 << 30]unsafe.Pointer)(unsafe.Pointer(res.schemas))[:len(paths):len(paths)]
	pathSchemas := make(map[string]*ParsedSchema, len(paths))
	for i := range paths {
		pathSchemas[paths[i]] = &ParsedSchema{
			ptr: cSchemasSlice[i],
		}
	}

	// release the returned array
	C.free(unsafe.Pointer(res.schemas))

	return pathSchemas, nil
}
