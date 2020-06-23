package dyncapnp

//go:generate go run ./gen/

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/WKBae/go-dyncapnp/schema"
)

var ErrSchemaNotFound = fmt.Errorf("schema not found")

// ParsedSchema of a Cap'n'proto type. MUST not be copied. Should be .Release()'ed after use.
type ParsedSchema struct {
	parser *schemaParser
	ptr    unsafe.Pointer
	*schema.Schema
}

// FindNested finds nested schema with given name. Returns ErrSchemaNotFound if nothing was found.
func (s *ParsedSchema) Nested(name string) (*ParsedSchema, error) {
	ptr, err := findNested(s.ptr, name)
	if err != nil {
		return nil, err
	}
	if ptr == nil {
		return nil, ErrSchemaNotFound
	}

	s.parser.incRef()
	sc := &ParsedSchema{
		parser: s.parser,
		ptr:    ptr,
		Schema: schema.NewWithFreer(ptr, schema.NoFree),
	}
	runtime.SetFinalizer(sc, (*ParsedSchema).Release)
	return sc, nil
}

func (s *ParsedSchema) Release() {
	s.Schema.Release()
	releaseParsedSchema(s.ptr)
	s.parser.decRef()
}

func ParseFromFiles(files map[string][]byte, imports map[string][]byte, paths []string) (map[string]*ParsedSchema, error) {
	// no need to parse if paths is empty
	if len(paths) == 0 {
		return nil, nil
	}

	// prepend standard imports
	importsWithStd := make(map[string][]byte, len(stdImports)+len(imports))
	for p, b := range stdImports {
		importsWithStd[p] = b
	}
	for p, b := range imports {
		importsWithStd[p] = b
	}

	parserPtr, schemaPtrs, err := parseSchemaFromFiles(files, importsWithStd, paths)
	if err != nil {
		return nil, err
	}

	parser := &schemaParser{
		ptr:      parserPtr,
		refCount: len(schemaPtrs),
	}

	pathSchemas := make(map[string]*ParsedSchema, len(paths))
	for i, path := range paths {
		pathSchemas[path] = &ParsedSchema{
			parser: parser,
			ptr:    schemaPtrs[i],
			Schema: schema.NewWithFreer(schemaPtrs[i], schema.NoFree),
		}
	}

	return pathSchemas, nil
}

type schemaParser struct {
	ptr      unsafe.Pointer
	refCount int
}

func (p *schemaParser) incRef() {
	if p == nil {
		return
	}
	p.refCount++
}

func (p *schemaParser) decRef() {
	if p == nil {
		return
	}
	p.refCount--
	if p.refCount < 0 {
		releaseParser(p.ptr)
	}
}
