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
)

const (
	UNSET     int64 = -1
	BMF_PAUSE int64 = 9223372036854775802
	EOF       int64 = 9223372036854775804
	EOS       int64 = 9223372036854775805
	INF_SRC   int64 = 9223372036854775806
	DONE      int64 = 9223372036854775807
)

type Ctype_bmf_TypeInfo = C.bmf_TypeInfo

type TypeInfo struct {
	p C.bmf_TypeInfo
}

func WrapTypeInfo(p C.bmf_TypeInfo) *TypeInfo {
	o := &TypeInfo{p: p}
	return o
}

func (self *TypeInfo) Name() string {
	return C.GoString(C.bmf_type_info_name(self.p))
}

func (self *TypeInfo) Index() uint64 {
	return uint64(C.bmf_type_info_index(self.p))
}

type Ctype_bmf_Packet = C.bmf_Packet

type Packet struct {
	p   C.bmf_Packet
	own bool
}

func deletePacket(o *Packet) {
	if o != nil && o.own && o.p != nil {
		C.bmf_packet_free(o.p)
		o.p = nil
	}
}

func WrapPacket(p C.bmf_Packet, own bool) *Packet {
	o := &Packet{p: p, own: own}
	runtime.SetFinalizer(o, deletePacket)
	return o
}

func GenerateEosPacket() *Packet {
	p := C.bmf_packet_generate_eos_packet()
	return WrapPacket(p, true)
}

func GenerateEofPacket() *Packet {
	p := C.bmf_packet_generate_eof_packet()
	return WrapPacket(p, true)
}

func (self *Packet) Free() {
	runtime.SetFinalizer(self, nil)
	deletePacket(self)
}

func (self *Packet) Timestamp() int64 {
	return int64(C.bmf_packet_timestamp(self.p))
}

func (self *Packet) SetTimestamp(timestamp int64) {
	C.bmf_packet_set_timestamp(self.p, C.int64_t(timestamp))
}

func (self *Packet) Defined() bool {
	return int(C.bmf_packet_defined(self.p)) != 0
}

func (self *Packet) Pointer() Ctype_bmf_Packet {
	return self.p
}

// Packet with VideoFrame
func NewPacketFromVideoFrame(vf *VideoFrame) (*Packet, error) {
	p := C.bmf_packet_from_videoframe(vf.p)
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapPacket(p, true), nil
}

func (self *Packet) GetVideoFrame() (*VideoFrame, error) {
	p := C.bmf_packet_get_videoframe(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapVideoFrame(p, true), nil
}

func (self *Packet) IsVideoFrame() bool {
	return int(C.bmf_packet_is_videoframe(self.p)) != 0
}

// Packet with AudioFrame
func NewPacketFromAudioFrame(af *AudioFrame) (*Packet, error) {
	p := C.bmf_packet_from_audioframe(af.p)
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapPacket(p, true), nil
}

func (self *Packet) GetAudioFrame() (*AudioFrame, error) {
	p := C.bmf_packet_get_audioframe(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapAudioFrame(p, true), nil
}

func (self *Packet) IsAudioFrame() bool {
	return int(C.bmf_packet_is_audioframe(self.p)) != 0
}

// Packet with BMFAVPacket
func NewPacketFromBMFAVPacket(avPkt *BMFAVPacket) (*Packet, error) {
	p := C.bmf_packet_from_bmfavpacket(avPkt.p)
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapPacket(p, true), nil
}

func (self *Packet) GetBMFAVPacket() (*BMFAVPacket, error) {
	p := C.bmf_packet_get_bmfavpacket(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapBMFAVPacket(p, true), nil
}

func (self *Packet) IsBMFAVPacket() bool {
	return int(C.bmf_packet_is_bmfavpacket(self.p)) != 0
}

// Packet with JsonParam
func NewPacketFromJsonParam(data interface{}) (*Packet, error) {
	str, _ := json.Marshal(data)
	cstr := C.CString(string(str))
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.bmf_json_param_parse(cstr)
	defer C.bmf_json_param_free(ptr)

	p := C.bmf_packet_from_json_param(ptr)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapPacket(p, true), nil
}

func (self *Packet) GetJsonParam(data interface{}) error {
	ptr := C.bmf_packet_get_json_param(self.p)
	if ptr == nil {
		return errors.New(LastError())
	}

	cstr := C.bmf_json_param_dump(ptr)
	defer C.free(unsafe.Pointer(cstr))
	gstr := []byte(C.GoString(cstr))
	err := json.Unmarshal(gstr, &data)
	return err
}

func (self *Packet) IsJsonParam() bool {
	return int(C.bmf_packet_is_json_param(self.p)) != 0
}
