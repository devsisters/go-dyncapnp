#include <exception>
#include <capnp/schema.h>
#include "cgo.h"

pointer_result typeFromPrimitive(uint16_t which) {
	try {
		auto type = new capnp::Type;
		*type = capnp::Type(static_cast<capnp::schema::Type::Which>(which));
		return {static_cast<void*>(type), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result typeFromStructSchema(void* structSchemaPtr) {
	try {
		auto schema = static_cast<capnp::StructSchema*>(structSchemaPtr);
		auto type = new capnp::Type;
		*type = capnp::Type(*schema);
		return {static_cast<void*>(type), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result typeFromEnumSchema(void* enumSchemaPtr) {
	try {
		auto schema = static_cast<capnp::EnumSchema*>(enumSchemaPtr);
		auto type = new capnp::Type;
		*type = capnp::Type(*schema);
		return {static_cast<void*>(type), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result typeFromInterfaceSchema(void* interfaceSchemaPtr) {
	try {
		auto schema = static_cast<capnp::InterfaceSchema*>(interfaceSchemaPtr);
		auto type = new capnp::Type;
		*type = capnp::Type(*schema);
		return {static_cast<void*>(type), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result typeFromListSchema(void* listSchemaPtr) {
	try {
		auto schema = static_cast<capnp::ListSchema*>(listSchemaPtr);
		auto type = new capnp::Type;
		*type = capnp::Type(*schema);
		return {static_cast<void*>(type), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

uint16_t typeWhich(void* typePtr) {
	auto type = static_cast<capnp::Type*>(typePtr);
	return static_cast<uint16_t>(type->which());
}

pointer_result typeAsStruct(void* typePtr) {
	try {
		auto type = static_cast<capnp::Type*>(typePtr);
		auto newSchema = new capnp::StructSchema;
		*newSchema = type->asStruct();
		return {static_cast<void*>(newSchema), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result typeAsEnum(void* typePtr) {
	try {
		auto type = static_cast<capnp::Type*>(typePtr);
		auto newSchema = new capnp::EnumSchema;
		*newSchema = type->asEnum();
		return {static_cast<void*>(newSchema), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result typeAsInterface(void* typePtr) {
	try {
		auto type = static_cast<capnp::Type*>(typePtr);
		auto newSchema = new capnp::InterfaceSchema;
		*newSchema = type->asInterface();
		return {static_cast<void*>(newSchema), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result typeAsList(void* typePtr) {
	try {
		auto type = static_cast<capnp::Type*>(typePtr);
		auto newSchema = new capnp::ListSchema;
		*newSchema = type->asList();
		return {static_cast<void*>(newSchema), nullptr};
	} catch (const std::exception& e) {
		return {nullptr, strdup(e.what())};
	}
}

void releaseType(void* typePtr) {
	auto type = static_cast<capnp::Type*>(typePtr);
	delete type;
}
