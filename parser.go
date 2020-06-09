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
	"unsafe"
)

// Parsed .capnp schema files. MUST be .Release()'ed after use.
type ParsedSchemas struct {
	ptr   unsafe.Pointer
	paths map[string]unsafe.Pointer
}

func (s *ParsedSchemas) Get(path string, structName string) Schema {
	pt, ok := s.paths[path]
	if !ok {
		panic("given path not exists")
	}

	st := C.CString(structName)
	defer C.free(unsafe.Pointer(st))

	ptr := C.findStructSchema(pt, st)
	if ptr == nil {
		panic("given struct not found in the file")
	}

	return Schema{
		ptr: ptr,
	}
}

func (s *ParsedSchemas) Release() {
	C.releaseSchemas(s.ptr, C.size_t(len(s.paths)))
	s.ptr = unsafe.Pointer(nil)
	s.paths = nil
}

type Schema struct {
	ptr unsafe.Pointer
}

func (s *Schema) Release() {
	C.releaseSchema(s.ptr)
	s.ptr = unsafe.Pointer(nil)
}

func ParseFromFiles(files map[string][]byte, paths []string) *ParsedSchemas {
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

	cPathSlice := make([]*C.char, len(paths))
	for i, path := range paths {
		cPathSlice[i] = C.CString(path)
	}
	defer func() {
		for _, chs := range cPathSlice {
			C.free(unsafe.Pointer(chs))
		}
	}()

	cSchemas := C.parseSchemaFromFiles(cFiles, C.size_t(len(cFileSlice)), &cPathSlice[0], C.size_t(len(cPathSlice)))
	cSchemasSlice := (*[1 << 30]unsafe.Pointer)(cSchemas)[:len(paths):len(paths)]

	pathSchemas := make(map[string]unsafe.Pointer, len(paths))
	for i := range paths {
		pathSchemas[paths[i]] = cSchemasSlice[i]
	}

	return &ParsedSchemas{
		ptr:   cSchemas,
		paths: pathSchemas,
	}
}
