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
	"reflect"
	"github.com/babitmf/bmf-gosdk/hmp"
)

type Ctype_bmf_BMFAVPacket = C.bmf_BMFAVPacket

type BMFAVPacket struct {
	p   C.bmf_BMFAVPacket
	own bool
}

func deleteBMFAVPacket(o *BMFAVPacket) {
	if o != nil && o.own && o.p != nil {
		C.bmf_pkt_free(o.p)
		o.p = nil
	}
}

func WrapBMFAVPacket(p C.bmf_BMFAVPacket, own bool) *BMFAVPacket {
	o := &BMFAVPacket{p: p, own: own}
	runtime.SetFinalizer(o, deleteBMFAVPacket)
	return o
}

func NewBMFAVPacketFromData(data *hmp.Tensor) (*BMFAVPacket, error) {
	p := C.bmf_pkt_make_from_data(C.hmp_Tensor(data.Pointer()))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapBMFAVPacket(p, true), nil
}

func NewBMFAVPacket(size int, dtype int) (*BMFAVPacket, error) {
	p := C.bmf_pkt_make(C.int(size), C.int(dtype))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapBMFAVPacket(p, true), nil
}

func (self *BMFAVPacket) Free() {
	runtime.SetFinalizer(self, nil)
	deleteBMFAVPacket(self)
}

func (self *BMFAVPacket) Pointer() C.bmf_BMFAVPacket {
	return self.p
}

func (self *BMFAVPacket) Defined() bool {
	return bool(C.bmf_pkt_defined(self.p))
}

func (self *BMFAVPacket) Data() (*hmp.Tensor, error) {
	p := C.bmf_pkt_data(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return hmp.WrapTensor(hmp.Ctype_hmp_Tensor(p), true), nil
}

func (self *BMFAVPacket) DataPtr() ([]byte, error) {
	p := C.bmf_pkt_data_ptr(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	siz := self.Nbytes();
	h := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(p)),
		Len:  siz,
		Cap:  siz,
	}

	return *((*[]byte)(unsafe.Pointer(h))), nil
}

func (self *BMFAVPacket) Nbytes() int {
	return int(C.bmf_pkt_nbytes(self.p))
}

func (self *BMFAVPacket) CopyProps(from *BMFAVPacket) {
	C.bmf_pkt_copy_props(self.p, from.p)
}

func (self *BMFAVPacket) PrivateMerge(from *BMFAVPacket) {
	C.bmf_pkt_private_merge(self.p, from.p)
}

func (self *BMFAVPacket) PrivateGet(key OpaqueDataKey, data interface{}) error {
	switch key {
	case kJsonParam:
		ptr := C.bmf_pkt_private_get_json_param(self.p)
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

func (self *BMFAVPacket) PrivateAttach(key OpaqueDataKey, option interface{}) {
	switch key {
	case kJsonParam:
		opt, _ := json.Marshal(option)
		cs := C.CString(string(opt))
		defer C.free(unsafe.Pointer(cs))
		ptr := C.bmf_json_param_parse(cs)
		defer C.bmf_json_param_free(ptr)
		C.bmf_pkt_private_attach_json_param(self.p, ptr)
	default:
		panic("Unknown opaque data key in PrivateAttach")
	}
}

func (self *BMFAVPacket) SetPts(pts int64) {
	C.bmf_pkt_set_pts(self.p, C.int64_t(pts))
}

func (self *BMFAVPacket) Pts() int64 {
	return int64(C.bmf_pkt_pts(self.p))
}

func (self *BMFAVPacket) SetTimeBase(num, den int) {
	C.bmf_pkt_set_time_base(self.p, C.int(num), C.int(den))
}

func (self *BMFAVPacket) TimeBase() (num, den int) {
	n, d := C.int(0), C.int(0)
	C.bmf_pkt_time_base(self.p, &n, &d)
	return int(n), int(d)
}

func (self *BMFAVPacket) GetOffset() (int64, error) {
	ci := C.int64_t(0)
	ci = C.bmf_pkt_offset(self.p)
	return int64(ci), nil
}

func (self *BMFAVPacket) GetWhence() (int64, error) {
	ci := C.int64_t(0)
	ci = C.bmf_pkt_whence(self.p)
	return int64(ci), nil
}
