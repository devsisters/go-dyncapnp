#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>

struct capnpFile {
	const char* path;
	char* content;
	size_t contentLen;
};

struct parseSchemaFromFiles_result {
	void** schemas;
	const char* err;
};

struct parseSchemaFromFiles_result parseSchemaFromFiles(const struct capnpFile* files, size_t filesLen, const struct capnpFile* imports, size_t importsLen, const char** paths, size_t pathsLen);

struct findStructSchema_result {
	void* schema;
	const char* err;
};

struct findStructSchema_result findStructSchema(void* schemaPtr, char* name);

void releaseSchemas(void** schemasPtr, size_t schemasLen);

void releaseSchema(void* schemaPtr);

#ifdef __cplusplus
}
#endif