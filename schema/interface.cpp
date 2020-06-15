#include <capnp/schema.h>
#include "schema.h"

void releaseInterfaceSchema(void* schemaPtr) {
	auto schema = static_cast<capnp::InterfaceSchema*>(schemaPtr);
	delete schema;
}
