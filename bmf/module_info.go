package bmf

/*
#include <stdint.h>
#include <stdlib.h>
#include <bmf/sdk/bmf_capi.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type ModuleInfo struct {
	ins C.bmf_ModuleInfo
}

func WrapModuleInfo(p C.bmf_ModuleInfo) ModuleInfo {
	return ModuleInfo{ins: p}
}

func NewModuleInfo() *ModuleInfo {
	o := &ModuleInfo{
		ins: C.bmf_module_info_make(),
	}
	if o.ins == nil {
		return nil
	}
	runtime.SetFinalizer(o, deleteModuleInfo)
	return o
}

func deleteModuleInfo(info *ModuleInfo) {
	if info != nil && info.ins != nil {
		C.bmf_module_info_free(info.ins)
		info.ins = nil
	}
}

func (mi *ModuleInfo) Free() {
	runtime.SetFinalizer(mi, nil)
	deleteModuleInfo(mi)
}

func (mi *ModuleInfo) SetModuleDescription(description string) {
	cs := C.CString(description)
	defer C.free(unsafe.Pointer(cs))

	C.bmf_module_info_set_description(mi.ins, cs)
}

func (mi *ModuleInfo) SetModuleTag(tag *ModuleTag) {
	C.bmf_module_info_set_tag(mi.ins, tag.GetIns())
}
