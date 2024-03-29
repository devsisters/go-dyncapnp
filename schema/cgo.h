#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>


typedef struct pointerList {
	void** ptr;
	size_t len;
} pointerList;

typedef struct pointer_result {
	void* ptr;
	const char* err;
} pointer_result;

typedef struct byteArray {
	char* arr;
	size_t length;
} byteArray;

typedef struct byteArray_result {
	struct byteArray result;
	const char* err;
} byteArray_result;


// Schema
void* schemaGetGeneric(void* schemaPtr);
void* schemaAsStruct(void* schemaPtr);
void* schemaAsEnum(void* schemaPtr);
void* schemaAsInterface(void* schemaPtr);
void* schemaAsConst(void* schemaPtr);
const char* schemaGetShortDisplayName(void* schemaPtr);
byteArray_result schemaToJson(void* schemaPtr);
void releaseSchema(void* schemaPtr);

// StructSchema
pointerList structGetFields(void* structSchemaPtr);
pointerList structGetUnionFields(void* structSchemaPtr);
pointerList structGetNonUnionFields(void* structSchemaPtr);
pointer_result structFindFieldByName(void* structSchemaPtr, const char* name);
byteArray_result structJsonToBinary(void* structSchemaPtr, const char* json, size_t len);
byteArray_result structJsonToPacked(void* structSchemaPtr, const char* json, size_t len);
byteArray_result structBinaryToJson(void* structSchemaPtr, const char* binary, size_t len);
byteArray_result structPackedToJson(void* structSchemaPtr, const char* binary, size_t len);
void releaseStructSchema(void* structSchemaPtr);

// StructSchema::Field
void* structFieldGetContainingStruct(void* structFieldPtr);
int structFieldGetIndex(void* structFieldPtr);
void* structFieldGetType(void* structFieldPtr);
byteArray_result structFieldToJson(void* structFieldPtr);
void releaseStructSchemaField(void* structFieldPtr);

// Type
pointer_result typeFromPrimitive(uint16_t which);
pointer_result typeFromStructSchema(void* structSchemaPtr);
pointer_result typeFromEnumSchema(void* enumSchemaPtr);
pointer_result typeFromInterfaceSchema(void* interfaceSchemaPtr);
pointer_result typeFromListSchema(void* listSchemaPtr);
uint16_t typeWhich(void* typePtr);
pointer_result typeAsStruct(void* typePtr);
pointer_result typeAsEnum(void* typePtr);
pointer_result typeAsInterface(void* typePtr);
pointer_result typeAsList(void* typePtr);
void releaseType(void* typePtr);

// ListSchema
void* listGetElementType(void* listPtr);
uint16_t listWhichElementType(void* listPtr);
pointer_result listGetStructElementType(void* listPtr);
pointer_result listGetEnumElementType(void* listPtr);
pointer_result listGetInterfaceElementType(void* listPtr);
pointer_result listGetListElementType(void* listPtr);
void releaseListSchema(void* listPtr);

// EnumSchema
pointerList enumGetEnumerants(void* enumPtr);
pointer_result enumFindEnumerantByName(void* enumPtr, const char* name);
void releaseEnum(void* enumPtr);
uint16_t enumerantGetOrdinal(void* enumerantPtr);
void releaseEnumerant(void* enumerantPtr);

// ConstSchema
void releaseConstSchema(void* schemaPtr);

// InterfaceSchema
void releaseInterfaceSchema(void* schemaPtr);

#ifdef __cplusplus
}
#endif
