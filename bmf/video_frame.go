package bmf

/*
#include <stdint.h>
#include <stdlib.h>
#include <bmf/sdk/bmf_capi.h>
*/
import "C"
import (
	"encoding/json"
	"errors"
	"runtime"
	"unsafe"

	"github.com/babitmf/bmf-gosdk/hmp"
)

type Ctype_bmf_VideoFrame = C.bmf_VideoFrame

type VideoFrame struct {
	p   C.bmf_VideoFrame
	own bool
}

func deleteVideoFrame(o *VideoFrame) {
	if o != nil && o.own && o.p != nil {
		C.bmf_vf_free(o.p)
		o.p = nil
	}
}

func WrapVideoFrame(p C.bmf_VideoFrame, own bool) *VideoFrame {
	o := &VideoFrame{p: p, own: own}
	runtime.SetFinalizer(o, deleteVideoFrame)
	return o
}

func NewVideoFrameFromFrame(frame *hmp.Frame) (*VideoFrame, error) {
	p := C.bmf_vf_from_frame((C.hmp_Frame)(frame.Pointer()))
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapVideoFrame(p, true), nil
}

func NewVideoFrameAsFrame(width, height int, pix_info *hmp.PixelInfo, device string) (*VideoFrame, error) {
	dstr := C.CString(device)
	defer C.free(unsafe.Pointer(dstr))
	dptr := (*C.char)(unsafe.Pointer(dstr))

	p := C.bmf_vf_make_frame(C.int(width), C.int(height), C.hmp_PixelInfo(pix_info.Pointer()), dptr)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapVideoFrame(p, true), nil
}

func (self *VideoFrame) Free() {
	runtime.SetFinalizer(self, nil)
	deleteVideoFrame(self)
}

func (self *VideoFrame) Pointer() C.bmf_VideoFrame {
	return self.p
}

func (self *VideoFrame) Defined() bool {
	return bool(C.bmf_vf_defined(self.p))
}

func (self *VideoFrame) Width() int {
	return int(C.bmf_vf_width(self.p))
}

func (self *VideoFrame) Height() int {
	return int(C.bmf_vf_height(self.p))
}

func (self *VideoFrame) Dtype() hmp.ScalarType {
	return hmp.ScalarType(C.bmf_vf_dtype(self.p))
}

func (self *VideoFrame) Frame() (*hmp.Frame, error) {
	p := C.bmf_vf_frame(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return hmp.WrapFrame(hmp.Ctype_hmp_Frame(p), false), nil
}

func (self *VideoFrame) DeviceType() hmp.DeviceType {
	return hmp.DeviceType(C.bmf_vf_device_type(self.p))
}

func (self *VideoFrame) DeviceIndex() int {
	return int(C.bmf_vf_device_index(self.p))
}

func (self *VideoFrame) CopyFrom(from *VideoFrame) {
	C.bmf_vf_copy_from(self.p, from.p)
}

func (self *VideoFrame) ToDevice(device string, non_blocking bool) (*VideoFrame, error) {
	dstr := C.CString(device)
	defer C.free(unsafe.Pointer(dstr))
	dptr := (*C.char)(unsafe.Pointer(dstr))

	p := C.bmf_vf_to_device(self.p, dptr, C.bool(non_blocking))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapVideoFrame(p, true), nil
}

func (self *VideoFrame) CopyProps(from *VideoFrame) {
	C.bmf_vf_copy_props(self.p, from.p)
}

func (self *VideoFrame) PrivateMerge(from *VideoFrame) {
	C.bmf_vf_private_merge(self.p, from.p)
}

func (self *VideoFrame) PrivateGet(key OpaqueDataKey, data interface{}) error {
	switch key {
	case kJsonParam:
		ptr := C.bmf_vf_private_get_json_param(self.p)
		if ptr == nil {
			return errors.New(LastError())
		}
		cstr := C.bmf_json_param_dump(ptr)
		defer C.free(unsafe.Pointer(cstr))
		gstr := []byte(C.GoString(cstr))
		err := json.Unmarshal(gstr, &data)
		if err != nil {
			return err
		}
		return nil
	default:
		panic("Unknown opaque data key in PrivateGet")
	}
}

func (self *VideoFrame) PrivateAttach(key OpaqueDataKey, option interface{}) {
	switch key {
	case kJsonParam:
		opt, _ := json.Marshal(option)
		cs := C.CString(string(opt))
		defer C.free(unsafe.Pointer(cs))
		ptr := C.bmf_json_param_parse(cs)
		defer C.bmf_json_param_free(ptr)
		C.bmf_vf_private_attach_json_param(self.p, ptr)
	default:
		panic("Unknown opaque data key in PrivateAttach")
	}
}

func (self *VideoFrame) ReFormat(pix_info *hmp.PixelInfo) (*VideoFrame, error) {
	p := C.bmf_vf_reformat(self.p, C.hmp_PixelInfo(pix_info.Pointer()))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapVideoFrame(p, true), nil
}

func (self *VideoFrame) SetPts(pts int64) {
	C.bmf_vf_set_pts(self.p, C.int64_t(pts))
}

func (self *VideoFrame) Pts() int64 {
	return int64(C.bmf_vf_pts(self.p))
}

func (self *VideoFrame) SetTimeBase(num, den int) {
	C.bmf_vf_set_time_base(self.p, C.int(num), C.int(den))
}

func (self *VideoFrame) TimeBase() (num, den int) {
	n, d := C.int(0), C.int(0)
	C.bmf_vf_time_base(self.p, &n, &d)
	return int(n), int(d)
}

func (self *VideoFrame) Ready() bool {
	return bool(C.bmf_vf_ready(self.p))
}

func (self *VideoFrame) Record(use_current bool) {
	C.bmf_vf_record(self.p, C.bool(use_current))
}

func (self *VideoFrame) Synchronize() {
	C.bmf_vf_synchronize(self.p)
}
