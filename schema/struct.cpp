#include <cstdlib>
#include <exception>
#include <capnp/schema.h>
#include "schema.h"

template <typename T>
pointerList listToPointers(T orig) {
	auto size = orig.size();
	auto list = static_cast<void**>(malloc(sizeof(void*) * size));
	for (int i = 0; i < size; i++) {
		auto item = new decltype(orig[i]);
		*item = orig[i];
		list[i] = static_cast<void*>(item);
	}
	return {list, size};
}

pointerList structGetFields(void* structSchemaPtr) {
	auto schema = static_cast<capnp::StructSchema*>(structSchemaPtr);
	return listToPointers(schema->getFields());
}

pointerList structGetUnionFields(void* structSchemaPtr) {
	auto schema = static_cast<capnp::StructSchema*>(structSchemaPtr);
	return listToPointers(schema->getUnionFields());
}

pointerList structGetNonUnionFields(void* structSchemaPtr) {
	auto schema = static_cast<capnp::StructSchema*>(structSchemaPtr);
	return listToPointers(schema->getNonUnionFields());
}

pointer_result structFindFieldByName(void* structSchemaPtr, const char* name) {
	try {
		auto schema = static_cast<capnp::StructSchema*>(structSchemaPtr);
		KJ_IF_MAYBE(ptr, schema->findFieldByName(name)) {
			auto field = new capnp::StructSchema::Field;
			*field = *ptr;
			return {static_cast<void*>(field), nullptr};
		} else {
			return {nullptr, nullptr};
		}
	} catch(const std::exception &e) {
		return {nullptr, strdup(e.what())};
	}
}

void releaseStructSchema(void* structSchemaPtr) {
	auto schema = static_cast<capnp::StructSchema*>(structSchemaPtr);
	delete schema;
}

void* structFieldGetContainingStruct(void* structFieldPtr) {
	auto field = static_cast<capnp::StructSchema::Field*>(structFieldPtr);
	auto schema = new capnp::StructSchema;
	*schema = field->getContainingStruct();
	return static_cast<void*>(schema);
}

int structFieldGetIndex(void* structFieldPtr) {
	auto field = static_cast<capnp::StructSchema::Field*>(structFieldPtr);
	return int(field->getIndex());
}

void* structFieldGetType(void* structFieldPtr) {
	auto field = static_cast<capnp::StructSchema::Field*>(structFieldPtr);
	auto type = new capnp::Type;
	*type = field->getType();
	return static_cast<void*>(type);
}

void releaseStructSchemaField(void* structFieldPtr) {
	auto field = static_cast<capnp::StructSchema::Field*>(structFieldPtr);
	delete field;
}
