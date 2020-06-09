#include <capnp/schema-parser.h>

#include "parser.h"

void* parseSchemaFromFiles(const struct capnpFile* files, size_t filesLen, const char** paths, size_t pathsLen) {
	auto dir = kj::newInMemoryDirectory(kj::nullClock());
	for (size_t i = 0; i < filesLen; i++) {
		auto path = kj::Path::parse(files[i].path);
		auto arr = kj::Array<const char>(files[i].content, files[i].contentLen, kj::NullArrayDisposer::instance);
		dir->openFile(path, kj::WriteMode::CREATE | kj::WriteMode::CREATE_PARENT)
			->writeAll(arr.asBytes());
	}

	auto schemas = new void*[pathsLen];
	auto p = new capnp::SchemaParser;
	kj::FixedArray<const kj::ReadableDirectory*, 1> importPath;
	importPath[0] = dir.get();
	for (size_t i = 0; i < pathsLen; i++) {
		auto schema = new capnp::ParsedSchema;
		*schema = p->parseFromDirectory(*dir, kj::Path::parse(paths[i]), importPath);
		schemas[i] = static_cast<void*>(schema);
	}

	return static_cast<void*>(schemas);
}

void* findStructSchema(void* schemaPtr, char* name) {
	auto schema = static_cast<capnp::ParsedSchema*>(schemaPtr);
	KJ_IF_MAYBE(ptr, schema->findNested(name)) {
		auto child = new capnp::ParsedSchema;
		*child = *ptr;
		return static_cast<void*>(child);
	} else {
		return nullptr;
	}
}

void releaseSchemas(void* schemasPtr, size_t schemasLen) {
	auto schemas = static_cast<void**>(schemasPtr);
	for (int i = 0; i < schemasLen; i++) {
		releaseSchema(schemas[i]);
	}
	delete schemas;
}

void releaseSchema(void* schemaPtr) {
	auto schema = static_cast<capnp::ParsedSchema*>(schemaPtr);
	delete schema;
}
