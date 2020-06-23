#include <capnp/schema.h>
#include "cgo.h"

void releaseInterfaceSchema(void* schemaPtr) {
	auto schema = static_cast<capnp::InterfaceSchema*>(schemaPtr);
	delete schema;
}
