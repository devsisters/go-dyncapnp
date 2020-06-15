package schema

/*
#cgo CXXFLAGS: -std=c++14 -stdlib=libc++ -I${SRCDIR}/capnproto/c++/src
#cgo LDFLAGS: -lkj -lcapnp -lcapnp-json

#include <stdlib.h>
#include "schema.h"
*/
import "C"
import (
	"unsafe"
)

func newInterface(ptr unsafe.Pointer) *Interface {
	return &Interface{
		Schema: NewWithFreer(ptr, releaseInterfaceSchema),
	}
}

func releaseInterfaceSchema(ptr unsafe.Pointer) {
	C.releaseInterfaceSchema(ptr)
}

type Interface struct {
	*Schema
}
