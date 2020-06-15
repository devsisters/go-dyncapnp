#include <capnp/schema.h>
#include "schema.h"

void* schemaGetGeneric(void* schemaPtr) {
	auto schema = static_cast<capnp::Schema*>(schemaPtr);
	auto newSchema = new capnp::Schema;
	*newSchema = schema->getGeneric();
	return static_cast<void*>(newSchema);
}

void* schemaAsStruct(void* schemaPtr) {
	auto schema = static_cast<capnp::Schema*>(schemaPtr);
	auto newSchema = new capnp::StructSchema;
	*newSchema = schema->asStruct();
	return static_cast<void*>(newSchema);
}

void* schemaAsEnum(void* schemaPtr) {
	auto schema = static_cast<capnp::Schema*>(schemaPtr);
	auto newSchema = new capnp::EnumSchema;
	*newSchema = schema->asEnum();
	return static_cast<void*>(newSchema);
}

void* schemaAsInterface(void* schemaPtr) {
	auto schema = static_cast<capnp::Schema*>(schemaPtr);
	auto newSchema = new capnp::InterfaceSchema;
	*newSchema = schema->asInterface();
	return static_cast<void*>(newSchema);
}

void* schemaAsConst(void* schemaPtr) {
	auto schema = static_cast<capnp::Schema*>(schemaPtr);
	auto newSchema = new capnp::ConstSchema;
	*newSchema = schema->asConst();
	return static_cast<void*>(newSchema);
}

void releaseSchema(void* schemaPtr) {
	auto schema = static_cast<capnp::Schema*>(schemaPtr);
	delete schema;
}
