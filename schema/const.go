package schema

import (
	"unsafe"
)

func newConst(ptr unsafe.Pointer) *Const {
	return &Const{
		Schema: NewWithFreer(ptr, releaseConstSchema),
	}
}

type Const struct {
	*Schema
}
