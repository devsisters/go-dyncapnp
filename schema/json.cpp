#include <algorithm>
#include <kj/string.h>
#include <capnp/schema-parser.h>
#include <capnp/message.h>
#include <capnp/serialize.h>
#include <capnp/serialize-packed.h>
#include <capnp/compat/json.h>
#include "schema.h"

byteArray cloneArray(const kj::ArrayPtr<capnp::byte>& ar) {
	auto arr = ar.asChars();
	byteArray b;
	b.length = arr.size();
	b.arr = static_cast<char*>(malloc(sizeof(char) * b.length)); // to release in Go with C.free()
	std::copy(arr.begin(), arr.end(), b.arr);
	return b;
}

byteArray_result schemaToJson(void* schemaPtr) {
	try {
		auto schema = static_cast<capnp::Schema*>(schemaPtr);

		capnp::JsonCodec codec;
		auto str = codec.encode(schema->getProto());

		return {cloneArray(str.asBytes()), nullptr};
	} catch(const std::exception &e) {
		return {{}, strdup(e.what())};
	}
}

byteArray_result structJsonToBinary(void* schemaPtr, const char* json, size_t len) {
	try {
		auto schema = static_cast<capnp::StructSchema*>(schemaPtr);
		kj::Array<const char> arr(json, len, kj::NullArrayDisposer::instance);

		capnp::MallocMessageBuilder builder;
		auto root = builder.initRoot<capnp::DynamicStruct>(*schema);

		capnp::JsonCodec codec;
		codec.decode(arr, root);

		kj::VectorOutputStream os;
		capnp::writeMessage(os, builder);

		return {cloneArray(os.getArray()), nullptr};
	} catch(const std::exception &e) {
		return {{}, strdup(e.what())};
	}
}

byteArray_result structJsonToPacked(void* schemaPtr, const char* json, size_t len) {
	try {
		auto schema = static_cast<capnp::StructSchema*>(schemaPtr);
		kj::Array<const char> arr(json, len, kj::NullArrayDisposer::instance);

		capnp::MallocMessageBuilder builder;
		auto root = builder.initRoot<capnp::DynamicStruct>(*schema);

		capnp::JsonCodec codec;
		codec.decode(arr, root);

		kj::VectorOutputStream os;
		capnp::writePackedMessage(os, builder);

		return {cloneArray(os.getArray()), nullptr};
	} catch(const std::exception &e) {
		return {{}, strdup(e.what())};
	}
}

byteArray_result structBinaryToJson(void* schemaPtr, const char* binary, size_t len) {
	try {
		auto schema = static_cast<capnp::StructSchema*>(schemaPtr);

		kj::Array<const char> arr(binary, len, kj::NullArrayDisposer::instance);
		kj::ArrayInputStream is(arr.asBytes());
		capnp::InputStreamMessageReader message(is);
		auto root = message.getRoot<capnp::DynamicStruct>(*schema);

		capnp::JsonCodec codec;
		auto str = codec.encode(root);

		return {cloneArray(str.asBytes()), nullptr};
	} catch(const std::exception &e) {
		return {{}, strdup(e.what())};
	}
}

byteArray_result structPackedToJson(void* schemaPtr, const char* binary, size_t len) {
	try {
		auto schema = static_cast<capnp::StructSchema*>(schemaPtr);

		kj::Array<const char> arr(binary, len, kj::NullArrayDisposer::instance);
		kj::ArrayInputStream is(arr.asBytes());
		capnp::PackedMessageReader message(is);
		auto root = message.getRoot<capnp::DynamicStruct>(*schema);

		capnp::JsonCodec codec;
		auto str = codec.encode(root);

		return {cloneArray(str.asBytes()), nullptr};
	} catch(const std::exception &e) {
		return {{}, strdup(e.what())};
	}
}
