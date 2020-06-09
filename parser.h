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

void* parseSchemaFromFiles(const struct capnpFile* files, size_t filesLen, const struct capnpFile* imports, size_t importsLen, const char** paths, size_t pathsLen);

void* findStructSchema(void* schemaPtr, char* name);

void releaseSchemas(void* schemasPtr, size_t schemasLen);

void releaseSchema(void* schemaPtr);

#ifdef __cplusplus
}
#endif