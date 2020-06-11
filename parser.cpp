#include <exception>

#include <capnp/schema-parser.h>

#include "parser.h"

parseSchemaFromFiles_result parseSchemaFromFiles(const struct capnpFile* files, size_t filesLen, const struct capnpFile* imports, size_t importsLen, const char** paths, size_t pathsLen) {
	try {
		auto dir = kj::newInMemoryDirectory(kj::nullClock());
		for (size_t i = 0; i < filesLen; i++) {
			auto path = kj::Path::parse(files[i].path);
			auto arr = kj::Array<const char>(files[i].content, files[i].contentLen, kj::NullArrayDisposer::instance);
			dir->openFile(path, kj::WriteMode::CREATE | kj::WriteMode::CREATE_PARENT)
				->writeAll(arr.asBytes());
		}

		auto importDir = kj::newInMemoryDirectory(kj::nullClock());
		for (size_t i = 0; i < importsLen; i++) {
			auto path = kj::Path::parse(imports[i].path);
			auto arr = kj::Array<const char>(imports[i].content, imports[i].contentLen, kj::NullArrayDisposer::instance);
			importDir->openFile(path, kj::WriteMode::CREATE | kj::WriteMode::CREATE_PARENT)
				->writeAll(arr.asBytes());
		}
		kj::FixedArray<const kj::ReadableDirectory*, 1> importPath;
		importPath[0] = importDir.get();

		auto schemas = new void*[pathsLen];
		auto p = new capnp::SchemaParser;
		for (size_t i = 0; i < pathsLen; i++) {
			auto schema = new capnp::ParsedSchema;
			*schema = p->parseFromDirectory(*dir, kj::Path::parse(paths[i]), importPath);
			schemas[i] = static_cast<void*>(schema);
		}

		return {schemas, nullptr};
	} catch(const std::exception &e) {
		return {nullptr, strdup(e.what())};
	}
}

findStructSchema_result findStructSchema(void* schemaPtr, char* name) {
	try {
		auto schema = static_cast<capnp::ParsedSchema*>(schemaPtr);
		KJ_IF_MAYBE(ptr, schema->findNested(name)) {
			auto child = new capnp::ParsedSchema;
			*child = *ptr;
			return {static_cast<void*>(child), nullptr};
		} else {
			return {nullptr, nullptr};
		}
	} catch(const std::exception &e) {
		return {nullptr, strdup(e.what())};
	}
}

void releaseSchemas(void** schemasPtr, size_t schemasLen) {
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
