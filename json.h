#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>

struct byteArray {
	char* arr;
	size_t length;
};

struct byteArray_result {
	struct byteArray result;
	const char* err;
};

struct byteArray_result schemaToJson(void* schemaPtr);

struct byteArray_result jsonToBinary(void* schemaPtr, const char* json, size_t len);

struct byteArray_result jsonToPacked(void* schemaPtr, const char* json, size_t len);

struct byteArray_result binaryToJson(void* schemaPtr, const char* binary, size_t len);

struct byteArray_result packedToJson(void* schemaPtr, const char* binary, size_t len);

#ifdef __cplusplus
}
#endif
