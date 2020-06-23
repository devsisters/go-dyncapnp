#include <cstdlib>
#include <exception>
#include <capnp/schema.h>
#include "schema.h"

pointerList enumGetEnumerants(void* enumPtr) {
	auto enu = static_cast<capnp::EnumSchema*>(enumPtr);
	auto ens = enu->getEnumerants();
	auto arr = static_cast<void**>(malloc(sizeof(void*) * ens.size()));
	for (int i = 0; i < ens.size(); i++) {
		auto en = new capnp::EnumSchema::Enumerant;
		*en = ens[i];
		arr[i] = en;
	}
	return {arr, ens.size()};
}

pointer_result enumFindEnumerantByName(void* enumPtr, const char* name) {
	try {
		auto enu = static_cast<capnp::EnumSchema*>(enumPtr);
		KJ_IF_MAYBE(ptr, enu->findEnumerantByName(name)) {
			auto enumerant = new capnp::EnumSchema::Enumerant;
			*enumerant = *ptr;
			return {static_cast<void*>(enumerant), nullptr};
		} else {
			return {nullptr, nullptr};
		}
	} catch(const std::exception &e) {
		return {nullptr, strdup(e.what())};
	}
}

void releaseEnum(void* enumPtr) {
	auto enu = static_cast<capnp::EnumSchema*>(enumPtr);
	delete enu;
}

uint16_t enumerantGetOrdinal(void* enumerantPtr) {
	auto enumerant = static_cast<capnp::EnumSchema::Enumerant*>(enumerantPtr);
	return enumerant->getOrdinal();
}

void releaseEnumerant(void* enumerantPtr) {
	auto enumerant = static_cast<capnp::EnumSchema::Enumerant*>(enumerantPtr);
	delete enumerant;
}
