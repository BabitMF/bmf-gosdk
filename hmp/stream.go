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
)

type Ctype_hmp_Stream = C.hmp_Stream
type Ctype_hmp_StreamGuard = C.hmp_StreamGuard

type Stream struct {
	p   C.hmp_Stream
	own bool
}

type StreamGuard struct {
	p   C.hmp_StreamGuard
	own bool
}

func deleteStream(o *Stream) {
	if o != nil && o.own && o.p != nil {
		C.hmp_stream_free(o.p)
		o.p = nil
	}
}

func WrapStream(p C.hmp_Stream, own bool) *Stream {
	o := &Stream{p: p, own: own}
	runtime.SetFinalizer(o, deleteStream)
	return o
}

func NewStream(device_type DeviceType, flags uint64) (*Stream, error) {
	p := C.hmp_stream_create(C.int(device_type), C.uint64_t(flags))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapStream(p, true), nil
}

func SetCurrentStream(stream *Stream) {
	C.hmp_stream_set_current(stream.p)
}

func CurrentStream(device_type DeviceType) (*Stream, error) {
	p := C.hmp_stream_current(C.int(device_type))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapStream(p, true), nil
}

func (self *Stream) Free() {
	runtime.SetFinalizer(self, nil)
	deleteStream(self)
}

func (self *Stream) Pointer() C.hmp_Stream {
	return self.p
}

func (self *Stream) Query() bool {
	return bool(C.hmp_stream_query(self.p))
}

func (self *Stream) Synchronize() {
	C.hmp_stream_synchronize(self.p)
}

func (self *Stream) Handle() uint64 {
	return uint64(C.hmp_stream_handle(self.p))
}

func (self *Stream) DeviceType() DeviceType {
	return DeviceType(C.hmp_stream_device_type(self.p))
}

func (self *Stream) DeviceIndex() int {
	return int(C.hmp_stream_device_index(self.p))
}

func deleteStreamGuard(o *StreamGuard) {
	if o != nil && o.own && o.p != nil {
		C.hmp_stream_guard_free(o.p)
		o.p = nil
	}
}

func WrapStreamGuard(p C.hmp_StreamGuard, own bool) *StreamGuard {
	o := &StreamGuard{p: p, own: own}
	runtime.SetFinalizer(o, deleteStreamGuard)
	return o
}

func NewStreamGuard(stream *Stream) (*StreamGuard, error) {
	p := C.hmp_stream_guard_create(stream.p)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapStreamGuard(p, true), nil
}

func (self *StreamGuard) Free() {
	runtime.SetFinalizer(self, nil)
	deleteStreamGuard(self)
}

func (self *StreamGuard) Pointer() C.hmp_StreamGuard {
	return self.p
}
