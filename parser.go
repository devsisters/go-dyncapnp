package dyncapnp

//go:generate go run ./gen/

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/devsisters/go-dyncapnp/schema"
)

var ErrSchemaNotFound = fmt.Errorf("schema not found")

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// ParsedSchema of a Cap'n'proto type. Should not be copied.
type ParsedSchema struct {
	parser *schemaParser
	ptr    unsafe.Pointer
	*schema.Schema

	noCopy noCopy
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
		Schema: schema.NewWithFreer(ptr, nil),
	}
	runtime.SetFinalizer(sc, (*ParsedSchema).release)
	return sc, nil
}

func (s *ParsedSchema) release() {
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
			Schema: schema.NewWithFreer(schemaPtrs[i], nil),
		}
	}

	return pathSchemas, nil
}

type schemaParser struct {
	ptr      unsafe.Pointer
	refCount int

	noCopy noCopy
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
