package dyncapnp

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
	"unsafe"
)

// Parsed .capnp schema files. MUST be .Release()'ed after use.
type ParsedSchemas struct {
	ptr   *unsafe.Pointer
	paths map[string]unsafe.Pointer
}

var ErrSchemaNotFound = fmt.Errorf("schema not found")

func (s *ParsedSchemas) Get(path string, structName string) (*Schema, error) {
	pt, ok := s.paths[path]
	if !ok {
		return nil, ErrSchemaNotFound
	}

	st := C.CString(structName)
	defer C.free(unsafe.Pointer(st))

	res := C.findStructSchema(pt, st)
	if res.err != nil {
		defer C.free(unsafe.Pointer(res.err))
		return nil, fmt.Errorf(C.GoString(res.err))
	}
	if res.schema == nil {
		return nil, ErrSchemaNotFound
	}

	return &Schema{
		ptr: res.schema,
	}, nil
}

func (s *ParsedSchemas) Release() {
	C.releaseSchemas(s.ptr, C.size_t(len(s.paths)))
	s.ptr = nil
	s.paths = nil
}

type Schema struct {
	ptr unsafe.Pointer
}

func (s *Schema) Release() {
	C.releaseSchema(s.ptr)
	s.ptr = unsafe.Pointer(nil)
}

func ParseFromFiles(files map[string][]byte, imports map[string][]byte, paths []string) (*ParsedSchemas, error) {
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
		return nil, err
	}
	cSchemasSlice := (*[1 << 30]unsafe.Pointer)(unsafe.Pointer(res.schemas))[:len(paths):len(paths)]

	pathSchemas := make(map[string]unsafe.Pointer, len(paths))
	for i := range paths {
		pathSchemas[paths[i]] = cSchemasSlice[i]
	}

	return &ParsedSchemas{
		ptr:   res.schemas,
		paths: pathSchemas,
	}, nil
}
