package bmf

/*
#cgo CXXFLAGS:   --std=c++11
#cgo LDFLAGS:    -lbmf_module_sdk
#include <stdint.h>
#include <stdlib.h>
#include <bmf/sdk/bmf_capi.h>

const char *_sdk_version()
{
    return BMF_SDK_VERSION;
}

*/
import "C"

type OpaqueDataKey int

func LastError() string {
	cstr := C.GoString(C.bmf_last_error())
	return cstr
}

func SdkVersion() string {
	cstr := C.GoString(C._sdk_version())
	return cstr
}

const (
	kJsonParam OpaqueDataKey = 2
)
