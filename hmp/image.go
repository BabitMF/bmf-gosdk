package hmp

/*
#include <stdint.h>
#include <stdlib.h>
#include "hmp_capi.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type Ctype_hmp_Frame = C.hmp_Frame

type Frame struct {
	p   C.hmp_Frame
	own bool
}

func deleteFrame(o *Frame) {
	if o != nil && o.own && o.p != nil {
		C.hmp_frame_free(o.p)
		o.p = nil
	}
}

func WrapFrame(p Ctype_hmp_Frame, own bool) *Frame {
	o := &Frame{p: p, own: own}
	runtime.SetFinalizer(o, deleteFrame)
	return o
}

/////////////////// Frame /////////////////

func NewFrame(width int, height int, pix_info *PixelInfo, device string) (*Frame, error) {
	dstr := C.CString(device)
	defer C.free(unsafe.Pointer(dstr))
	dptr := (*C.char)(unsafe.Pointer(dstr))
	p := C.hmp_frame_make(C.int(width), C.int(height), pix_info.p, dptr)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapFrame(p, true), nil
}

func NewFrameFromData(data []*Tensor, pix_info *PixelInfo) (*Frame, error) {
	// FIX: cgo argument has Go pointer to Go pointer
	sz := len(data)
	dptr := (*C.hmp_Tensor)(C.malloc(C.size_t(sz * 8)))
	defer C.free(unsafe.Pointer(dptr))
	darr := (*[1 << 31]C.hmp_Tensor)(unsafe.Pointer(dptr))[:sz:sz]
	for i := 0; i < sz; i++ {
		darr[i] = data[i].p
	}

	p := C.hmp_frame_from_data(dptr, C.int(sz), pix_info.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapFrame(p, true), nil
}

func NewFrameFromDataV1(data []*Tensor, width int, height int, pix_info *PixelInfo) (*Frame, error) {
	sz := len(data)
	dptr := (*C.hmp_Tensor)(C.malloc(C.size_t(sz * 8)))
	defer C.free(unsafe.Pointer(dptr))
	darr := (*[1 << 31]C.hmp_Tensor)(unsafe.Pointer(dptr))[:sz:sz]
	for i := 0; i < sz; i++ {
		darr[i] = data[i].p
	}

	p := C.hmp_frame_from_data_v1(dptr, C.int(sz), C.int(width), C.int(height), pix_info.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapFrame(p, true), nil
}


func (self *Frame) Free() {
	runtime.SetFinalizer(self, nil)
	deleteFrame(self)
}

func (self *Frame) Pointer() C.hmp_Frame {
	return self.p
}

func (self *Frame) Defined() bool {
	return bool(C.hmp_frame_defined(self.p))
}

func (self *Frame) PixInfo() *PixelInfo {
	p := C.hmp_frame_pix_info(self.p)
	return WrapPixelInfo(p, false)
}

func (self *Frame) Format() PixelFormat {
	return PixelFormat(C.hmp_frame_format(self.p))
}

func (self *Frame) Width() int {
	return int(C.hmp_frame_width(self.p))
}

func (self *Frame) Height() int {
	return int(C.hmp_frame_height(self.p))
}

func (self *Frame) Dtype() ScalarType {
	return ScalarType(C.hmp_frame_dtype(self.p))
}

func (self *Frame) DeviceType() DeviceType {
	return DeviceType(C.hmp_frame_device_type(self.p))
}

func (self *Frame) DeviceIndex() int {
	return int(C.hmp_frame_device_index(self.p))
}

func (self *Frame) Nplanes() int {
	return int(C.hmp_frame_nplanes(self.p))
}

func (self *Frame) Plane(plane int) *Tensor {
	p := C.hmp_frame_plane(self.p, C.int64_t(plane))
	return WrapTensor(p, false)
}

func (self *Frame) ToDevice(device string, non_blocking bool) (*Frame, error) {
	dstr := C.CString(device)
	defer C.free(unsafe.Pointer(dstr))
	dptr := (*C.char)(unsafe.Pointer(dstr))

	p := C.hmp_frame_to_device(self.p, dptr, C.bool(non_blocking))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapFrame(p, true), nil
}

func (self *Frame) CopyFrom(from *Frame) {
	C.hmp_frame_copy_from(self.p, from.p)
}

func (self *Frame) Clone() (*Frame, error) {
	p := C.hmp_frame_clone(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapFrame(p, true), nil
}

func (self *Frame) Crop(left int, top int, width int, height int) (*Frame, error) {
	p := C.hmp_frame_crop(self.p, C.int(left), C.int(top), C.int(width), C.int(height))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapFrame(p, true), nil
}

func (self *Frame) ReFormat(pix_info *PixelInfo) (*Frame, error) {
	p := C.hmp_frame_reformat(self.p, pix_info.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapFrame(p, true), nil
}

func (self *Frame) String() string {
	var c_str_size = int(0)
	c_str := C.hmp_frame_stringfy(self.p, (*C.int)(unsafe.Pointer(&c_str_size)))
	return string((*[1 << 31]byte)(unsafe.Pointer(c_str))[:c_str_size:c_str_size])
}
