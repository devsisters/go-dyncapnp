#include <capnp/schema.h>
#include "schema.h"

void releaseConstSchema(void* schemaPtr) {
	auto schema = static_cast<capnp::ConstSchema*>(schemaPtr);
	delete schema;
}
