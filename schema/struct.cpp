#include <cstdlib>
#include <capnp/schema.h>
#include "schema.h"

struct pointerList structGetFields(void* structSchemaPtr) {
	auto schema = static_cast<capnp::StructSchema*>(structSchemaPtr);
	auto fields = schema->getFields();
	auto list = static_cast<void**>(malloc(sizeof(void*) * fields.size()));
	for (int i = 0; i < fields.size(); i++) {
		auto field = new capnp::StructSchema::Field;
		*field = fields[i];
		list[i] = static_cast<void*>(field);
	}
	return {list, fields.size()};
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
