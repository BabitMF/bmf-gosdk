package builder

/*
#cgo CXXFLAGS:  --std=c++11
#cgo LDFLAGS:   -lhmp -lbmf_module_sdk -lengine

#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include "connector_capi.h"

*/
import "C"

func LastError() string {
	cstr := C.GoString(C.bmf_engine_last_error())
	return cstr
}
