#include <exception>
#include <capnp/schema.h>
#include "schema.h"

void* listGetElementType(void* listPtr) {
	auto list = static_cast<capnp::ListSchema*>(listPtr);
	auto type = new capnp::Type;
	*type = list->getElementType();
	return static_cast<void*>(type);
}

uint16_t listWhichElementType(void* listPtr) {
	auto list = static_cast<capnp::ListSchema*>(listPtr);
	return static_cast<uint16_t>(list->whichElementType());
}

pointer_result listGetStructElementType(void* listPtr) {
	try {
		auto list = static_cast<capnp::ListSchema*>(listPtr);
		auto schema = new capnp::StructSchema;
		*schema = list->getStructElementType();
		return {static_cast<void*>(schema), nullptr};
	} catch(const std::exception &e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result listGetEnumElementType(void* listPtr) {
	try {
		auto list = static_cast<capnp::ListSchema*>(listPtr);
		auto schema = new capnp::EnumSchema;
		*schema = list->getEnumElementType();
		return {static_cast<void*>(schema), nullptr};
	} catch(const std::exception &e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result listGetInterfaceElementType(void* listPtr) {
	try {
		auto list = static_cast<capnp::ListSchema*>(listPtr);
		auto schema = new capnp::InterfaceSchema;
		*schema = list->getInterfaceElementType();
		return {static_cast<void*>(schema), nullptr};
	} catch(const std::exception &e) {
		return {nullptr, strdup(e.what())};
	}
}

pointer_result listGetListElementType(void* listPtr) {
	try {
		auto list = static_cast<capnp::ListSchema*>(listPtr);
		auto schema = new capnp::ListSchema;
		*schema = list->getListElementType();
		return {static_cast<void*>(schema), nullptr};
	} catch(const std::exception &e) {
		return {nullptr, strdup(e.what())};
	}
}

void releaseListSchema(void* listPtr) {
	auto list = static_cast<capnp::ListSchema*>(listPtr);
	delete list;
}
