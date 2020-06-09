#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>

struct byteArray {
	char* arr;
	size_t length;

};

struct byteArray jsonToBinary(void* schemaPtr, const char* json, size_t len);

struct byteArray jsonToPacked(void* schemaPtr, const char* json, size_t len);

struct byteArray binaryToJson(void* schemaPtr, const char* binary, size_t len);

struct byteArray packedToJson(void* schemaPtr, const char* binary, size_t len);

#ifdef __cplusplus
}
#endif
