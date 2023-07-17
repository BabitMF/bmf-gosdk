package bmf

/*
#include <stdint.h>
#include <stdlib.h>
#include <bmf/sdk/bmf_capi.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type Ctype_bmf_Task = C.bmf_Task

type Task struct {
	p   C.bmf_Task
	own bool
}

func deleteTask(o *Task) {
	if o != nil && o.own && o.p != nil {
		C.bmf_task_free(o.p)
		o.p = nil
	}
}

func WrapTask(p C.bmf_Task, own bool) *Task {
	o := &Task{p: p, own: own}
	runtime.SetFinalizer(o, deleteTask)
	return o
}

func NewTask(node_id int32, istream_ids []int32, ostream_ids []int32) (*Task, error) {
	var iids, oids *C.int
	if len(istream_ids) > 0 {
		iids = (*C.int)(unsafe.Pointer(&istream_ids[0]))
	} else {
		iids = nil
	}

	if len(ostream_ids) > 0 {
		oids = (*C.int)(unsafe.Pointer(&ostream_ids[0]))
	} else {
		oids = nil
	}

	p := C.bmf_task_make(C.int(node_id), iids, C.int(len(istream_ids)),
		oids, C.int(len(ostream_ids)))
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapTask(p, true), nil
}

func (self *Task) Free() {
	runtime.SetFinalizer(self, nil)
	deleteTask(self)
}

func (self *Task) FillInputPacket(stream_id int32, packet *Packet) bool {
	rc := C.bmf_task_fill_input_packet(self.p, C.int(stream_id), packet.p)
	return rc != 0
}

func (self *Task) FillOutputPacket(stream_id int32, packet *Packet) bool {
	rc := C.bmf_task_fill_output_packet(self.p, C.int(stream_id), packet.p)
	return rc != 0
}

func (self *Task) PopPacketFromOutQueue(stream_id int32) (*Packet, error) {
	p := C.bmf_task_pop_packet_from_out_queue(self.p, C.int(stream_id))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapPacket(p, true), nil
}

func (self *Task) PopPacketFromInputQueue(stream_id int32) (*Packet, error) {
	p := C.bmf_task_pop_packet_from_input_queue(self.p, C.int(stream_id))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapPacket(p, true), nil
}

func (self *Task) Timestamp() int64 {
	rc := C.bmf_task_timestamp(self.p)
	return int64(rc)
}

func (self *Task) SetTimestamp(timestamp int64) {
	C.bmf_task_set_timestamp(self.p, C.int64_t(timestamp))
}

func (self *Task) GetInputStreamIds() []int32 {
	n := C.bmf_task_get_input_stream_ids(self.p, nil)
	ids := make([]int32, n)
	if n > 0 {
		ids_ptr := (*C.int)(unsafe.Pointer(&ids[0]))
		C.bmf_task_get_input_stream_ids(self.p, ids_ptr)
	}
	return ids
}

func (self *Task) GetOutputStreamIds() []int32 {
	n := C.bmf_task_get_output_stream_ids(self.p, nil)
	ids := make([]int32, n)
	if n > 0 {
		ids_ptr := (*C.int)(unsafe.Pointer(&ids[0]))
		C.bmf_task_get_output_stream_ids(self.p, ids_ptr)
	}
	return ids
}

func (self *Task) Pointer() Ctype_bmf_Task {
	return self.p
}
