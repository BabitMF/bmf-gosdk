package hmp

/*
#include <stdint.h>
#include <stdlib.h>
#include "hmp_capi.h"
*/
import "C"
import "runtime"

type Ctype_hmp_Scalar = C.hmp_Scalar

type Scalar struct {
	p   C.hmp_Scalar
	own bool
}

func deleteScalar(o *Scalar) {
	if o != nil && o.own && o.p != nil {
		C.hmp_scalar_free(o.p)
		o.p = nil
	}
}

func WrapScalar(p C.hmp_Scalar, own bool) *Scalar {
	o := &Scalar{p: p, own: own}
	runtime.SetFinalizer(o, deleteScalar)
	return o
}

func NewFloatScalar(value float64) *Scalar {
	p := C.hmp_scalar_float(C.double(value))
	return WrapScalar(p, true)
}

func NewBoolScalar(value bool) *Scalar {
	p := C.hmp_scalar_bool(C.bool(value))
	return WrapScalar(p, true)
}

func NewIntScalar(value int64) *Scalar {
	p := C.hmp_scalar_int(C.int64_t(value))
	return WrapScalar(p, true)
}

func (self *Scalar) Free() {
	runtime.SetFinalizer(self, nil)
	deleteScalar(self)
}

func (self *Scalar) Pointer() C.hmp_Scalar {
	return self.p
}
