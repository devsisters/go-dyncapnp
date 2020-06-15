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

func newConst(ptr unsafe.Pointer) *Const {
	return &Const{
		Schema: NewWithFreer(ptr, releaseConstSchema),
	}
}

func releaseConstSchema(ptr unsafe.Pointer) {
	C.releaseConstSchema(ptr)
}

type Const struct {
	*Schema
}
