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

type ModuleTagType int64

var _tmp ModuleTagType

const (
	BMF_TAG_NONE            ModuleTagType = 0x0
	BMF_TAG_DECODER                       = 0x01 << 0
	BMF_TAG_ENCODER                       = 0x01 << 1
	BMF_TAG_FILTER                        = 0x01 << 2
	BMF_TAG_MUXER                         = 0x01 << 3
	BMF_TAG_DEMUXER                       = 0x01 << 4
	BMF_TAG_IMAGE_PROCESSOR               = 0x01 << 5
	BMF_TAG_AUDIO_PROCESSOR               = 0x01 << 6
	BMF_TAG_VIDEO_PROCESSOR               = 0x01 << 7
	BMF_TAG_DEVICE_HWACCEL                = 0x01 << 8
	BMF_TAG_AI                            = 0x01 << 9
	BMF_TAG_UTILS                         = 0x01 << 10
	BMF_TAG_DONE                          = 0x01 << (unsafe.Sizeof(_tmp)*8 - 1)
)

type ModuleTag struct {
	ins C.bmf_ModuleTag
}

func NewModuleTag(tag ModuleTagType) *ModuleTag {
	o := &ModuleTag{
		ins: C.bmf_module_tag_make(C.int64_t(tag)),
	}
	if o.ins == nil {
		return nil
	}

	runtime.SetFinalizer(o, deleteModuleTag)
	return o
}

func deleteModuleTag(tag *ModuleTag) {
	if tag != nil && tag.ins != nil {
		C.bmf_module_tag_free(tag.ins)
		tag.ins = nil
	}
}

func (mt *ModuleTag) Free() {
	runtime.SetFinalizer(mt, nil)
	deleteModuleTag(mt)
}

func (mt *ModuleTag) GetIns() C.bmf_ModuleTag {
	return mt.ins
}
