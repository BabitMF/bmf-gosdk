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

type Ctype_bmf_AudioFrame = C.bmf_AudioFrame

type AudioFrame struct {
	p   C.bmf_AudioFrame
	own bool
}

func deleteAudioFrame(o *AudioFrame) {
	if o != nil && o.own && o.p != nil {
		C.bmf_af_free(o.p)
		o.p = nil
	}
}

func WrapAudioFrame(p C.bmf_AudioFrame, own bool) *AudioFrame {
	o := &AudioFrame{p: p, own: own}
	runtime.SetFinalizer(o, deleteAudioFrame)
	return o
}

func NewAudioFrameFromData(data []*hmp.Tensor, size int, layout int, planer bool) (*AudioFrame, error) {
	sz := len(data)
	dptr := (*C.hmp_Tensor)(C.malloc(C.size_t(sz * 8)))
	defer C.free(unsafe.Pointer(dptr))
	darr := (*[1 << 31]C.hmp_Tensor)(unsafe.Pointer(dptr))[:sz:sz]
	for i := 0; i < sz; i++ {
		darr[i] = (C.hmp_Tensor)(data[i].Pointer())
	}
	p := C.bmf_af_make_from_data(dptr, C.int(size), C.uint64_t(layout), C.bool(planer))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapAudioFrame(p, true), nil
}

func NewAudioFrame(samples int, layout int, planer bool, dtype int) (*AudioFrame, error) {
	p := C.bmf_af_make(C.int(samples), C.uint64_t(layout), C.bool(planer), C.int(dtype))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapAudioFrame(p, true), nil
}

func (self *AudioFrame) Free() {
	runtime.SetFinalizer(self, nil)
	deleteAudioFrame(self)
}

func (self *AudioFrame) Pointer() C.bmf_AudioFrame {
	return self.p
}

func (self *AudioFrame) Defined() bool {
	return bool(C.bmf_af_defined(self.p))
}

func (self *AudioFrame) Dtype() hmp.ScalarType {
	return hmp.ScalarType(C.bmf_af_dtype(self.p))
}

func (self *AudioFrame) Planer() bool {
	return bool(C.bmf_af_planer(self.p))
}

func (self *AudioFrame) Nsamples() int {
	return int(C.bmf_af_nsamples(self.p))
}

func (self *AudioFrame) Nchannels() int {
	return int(C.bmf_af_nchannels(self.p))
}

func (self *AudioFrame) SetSampleRate(ar float32) {
	C.bmf_af_set_sample_rate(self.p, C.float(ar))
}

func (self *AudioFrame) SampleRate() float32 {
	return float32(C.bmf_af_sample_rate(self.p))
}

func (self *AudioFrame) Planes() ([]*hmp.Tensor, error) {
	sz := self.Nplanes()
	data := make([]*hmp.Tensor, sz)
	dptr := (*C.hmp_Tensor)(C.malloc(C.size_t(sz * 8)))
	defer C.free(unsafe.Pointer(dptr))
	darr := (*[1 << 31]C.hmp_Tensor)(unsafe.Pointer(dptr))[:sz:sz]
	C.bmf_af_planes(self.p, dptr)
	for i := 0; i < sz; i++ {
		data[i] = hmp.WrapTensor(hmp.Ctype_hmp_Tensor(darr[i]), true)
	}
	return data, nil
}

func (self *AudioFrame) Nplanes() int {
	return int(C.bmf_af_nplanes(self.p))
}

func (self *AudioFrame) Plane(i int) (*hmp.Tensor, error) {
	p := C.bmf_af_plane(self.p, C.int(i))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return hmp.WrapTensor(hmp.Ctype_hmp_Tensor(p), true), nil
}

func (self *AudioFrame) CopyProps(from *AudioFrame) {
	C.bmf_af_copy_props(self.p, from.p)
}

func (self *AudioFrame) PrivateMerge(from *AudioFrame) {
	C.bmf_af_private_merge(self.p, from.p)
}

func (self *AudioFrame) PrivateGet(key OpaqueDataKey, data interface{}) error {
	switch key {
	case kJsonParam:
		ptr := C.bmf_af_private_get_json_param(self.p)
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

func (self *AudioFrame) PrivateAttach(key OpaqueDataKey, option interface{}) {
	switch key {
	case kJsonParam:
		opt, _ := json.Marshal(option)
		cs := C.CString(string(opt))
		defer C.free(unsafe.Pointer(cs))
		ptr := C.bmf_json_param_parse(cs)
		defer C.bmf_json_param_free(ptr)
		C.bmf_af_private_attach_json_param(self.p, ptr)
	default:
		panic("Unknown opaque data key in PrivateAttach")
	}
}

func (self *AudioFrame) SetPts(pts int64) {
	C.bmf_af_set_pts(self.p, C.int64_t(pts))
}

func (self *AudioFrame) Pts() int64 {
	return int64(C.bmf_af_pts(self.p))
}

func (self *AudioFrame) SetTimeBase(num, den int) {
	C.bmf_af_set_time_base(self.p, C.int(num), C.int(den))
}

func (self *AudioFrame) TimeBase() (num, den int) {
	n, d := C.int(0), C.int(0)
	C.bmf_af_time_base(self.p, &n, &d)
	return int(n), int(d)
}
