package builder

/*
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include "connector_capi.h"
#include "bmf/sdk/bmf_capi.h"
*/
import "C"
import (
	"encoding/json"
	"errors"
	"runtime"
	"unsafe"

	"github.com/babitmf/bmf-gosdk/bmf"
)

type Ctype_bmf_BMFGraph = C.bmf_BMFGraph

type CBMFGraph struct {
	p   C.bmf_BMFGraph
	own bool
}

func deleteCBMFGraph(o *CBMFGraph) {
	if o != nil && o.own && o.p != nil {
		C.bmf_graph_free(o.p)
		o.p = nil
	}
}

func WrapCBMFGraph(p C.bmf_BMFGraph, own bool) *CBMFGraph {
	o := &CBMFGraph{p: p, own: own}
	runtime.SetFinalizer(o, deleteCBMFGraph)
	return o
}

func NewCBMFGraph(config interface{}, needMerge bool) (*CBMFGraph, error) {
	graphConf, _ := json.Marshal(config)
	cs := C.CString(string(graphConf))
	defer C.free(unsafe.Pointer(cs))
	p := C.bmf_make_graph(cs, C.bool(false), C.bool(needMerge))
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapCBMFGraph(p, true), nil
}

func (self *CBMFGraph) UID() uint32 {
	return uint32(C.bmf_graph_uid(self.p))
}

func (self *CBMFGraph) Start() error {
	rc := C.bmf_graph_start(self.p)
	if rc < 0 {
		return errors.New(LastError())
	}
	return nil
}

func (self *CBMFGraph) Close() error {
	rc := C.bmf_graph_close(self.p)
	if rc < 0 {
		return errors.New(LastError())
	}
	return nil
}

func (self *CBMFGraph) AddInputStreamPacket(streamName string, packet *bmf.Packet) error {
	cs := C.CString(streamName)
	defer C.free(unsafe.Pointer(cs))
	rc := C.bmf_graph_add_input_stream_packet(self.p, cs, C.bmf_Packet(packet.Pointer()), false)
	if rc < 0 {
		return errors.New(LastError())
	}
	return nil
}

func (self *CBMFGraph) PollOutputStreamPacket(streamName string) (*bmf.Packet, error) {
	cs := C.CString(streamName)
	defer C.free(unsafe.Pointer(cs))
	pkt := C.bmf_graph_poll_output_stream_packet(self.p, cs)
	if pkt == nil {
		return nil, errors.New(LastError())
	}
	return bmf.WrapPacket(bmf.Ctype_bmf_Packet(pkt), true), nil
}

func (self *CBMFGraph) Update(info interface{}) error {
	ud, _ := json.Marshal(info)
	cs := C.CString(string(ud))
	defer C.free(unsafe.Pointer(cs))
	rc := C.bmf_graph_update(self.p, cs, false)
	if rc < 0 {
		return errors.New(LastError())
	}
	return nil
}

func (self *CBMFGraph) ForceClose() error {
	rc := C.bmf_graph_force_close(self.p)
	if rc < 0 {
		return errors.New(LastError())
	}
	return nil
}

func (self *CBMFGraph) Status() (string, error) {
	cs := C.bmf_graph_status(self.p)
	if cs == nil {
		return "", errors.New(LastError())
	}
	defer C.free(unsafe.Pointer(cs))
	return C.GoString(cs), nil
}

type Ctype_bmf_BMFModule = C.bmf_BMFModule

type CBMFModule struct {
	p   C.bmf_BMFModule
	own bool
}

func deleteCBMFModule(o *CBMFModule) {
	if o != nil && o.own && o.p != nil {
		C.bmf_module_free(o.p)
		o.p = nil
	}
}

func WrapCBMFModule(p C.bmf_BMFModule, own bool) *CBMFModule {
	o := &CBMFModule{p: p, own: own}
	runtime.SetFinalizer(o, deleteCBMFModule)
	return o
}

func NewCBMFModule(moduleName string, option interface{}, moduleType BMFModuleType, modulePath, moduleEntry string) (*CBMFModule, error) {
	opt, _ := json.Marshal(option)
	copts := C.CString(string(opt))
	defer C.free(unsafe.Pointer(copts))

	cmoduleName, cmoduleType, cmodulePath, cmoduleEntry := C.CString(moduleName), C.CString(moduleTypeToString(moduleType)), C.CString(modulePath), C.CString(moduleEntry)
	defer C.free(unsafe.Pointer(cmoduleName))
	defer C.free(unsafe.Pointer(cmoduleType))
	defer C.free(unsafe.Pointer(cmodulePath))
	defer C.free(unsafe.Pointer(cmoduleEntry))

	p := C.bmf_make_module(cmoduleName, copts, cmoduleType, cmodulePath, cmoduleEntry)
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapCBMFModule(p, true), nil
}

func (self *CBMFModule) UID() uint32 {
	return uint32(C.bmf_module_uid(self.p))
}

func (self *CBMFModule) Process(task *bmf.Task) error {
	rc := C.bmf_module_process(self.p, C.bmf_Task(task.Pointer()))
	if rc < 0 {
		return errors.New(LastError())
	}
	return nil
}

func (self *CBMFModule) Init() error {
	rc := C.bmf_module_init(self.p)
	if rc < 0 {
		return errors.New(LastError())
	}
	return nil
}

func (self *CBMFModule) Close() error {
	rc := C.bmf_module_close(self.p)
	if rc < 0 {
		return errors.New(LastError())
	}
	return nil
}
