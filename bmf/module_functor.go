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

type Ctype_bmf_ModuleFunctor = C.bmf_ModuleFunctor

type ModuleFunctor struct {
	p   C.bmf_ModuleFunctor
	own bool
}

func deleteModuleFunctor(o *ModuleFunctor) {
	if o != nil && o.own && o.p != nil {
		C.bmf_module_functor_free(o.p)
		o.p = nil
	}
}

func WrapModuleFunctor(p C.bmf_ModuleFunctor, own bool) *ModuleFunctor {
	o := &ModuleFunctor{p: p, own: own}
	runtime.SetFinalizer(o, deleteModuleFunctor)
	return o
}

func NewModuleFunctor(name, tp, path, entry string, option interface{}, ninputs, noutputs int32) (*ModuleFunctor, error) {
	c_name := C.CString(name)
	c_type := C.CString(tp)
	c_path := C.CString(path)
	c_entry := C.CString(entry)
	defer C.free(unsafe.Pointer(c_name))
	defer C.free(unsafe.Pointer(c_type))
	defer C.free(unsafe.Pointer(c_path))
	defer C.free(unsafe.Pointer(c_entry))

	option_str, err := json.Marshal(option)
	if err != nil {
		return nil, err
	}
	c_option_str := C.CString(string(option_str))
	defer C.free(unsafe.Pointer(c_option_str))

	p := C.bmf_module_functor_make(c_name, c_type, c_path, c_entry, c_option_str,
		C.int(ninputs), C.int(noutputs), -1)
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapModuleFunctor(p, true), nil
}

func NewModuleFunctorBuiltin(name string, option interface{}, ninputs, noutputs int32) (*ModuleFunctor, error) {
	return NewModuleFunctor(name, "", "", "", option, ninputs, noutputs)
}

func (self *ModuleFunctor) Free() {
	runtime.SetFinalizer(self, nil)
	deleteModuleFunctor(self)
}

func (self *ModuleFunctor) Call(inputs []*Packet) ([]*Packet, error) {
	isz := len(inputs)
	dptr := (*C.bmf_Packet)(C.malloc(C.size_t(isz * 8)))
	defer C.free(unsafe.Pointer(dptr))
	darr := (*[1 << 31]C.bmf_Packet)(unsafe.Pointer(dptr))[:isz:isz]
	for i := 0; i < isz; i++ {
		darr[i] = inputs[i].Pointer()
	}

	osz := C.int(0)
	is_done := C.bool(false)
	obuf := C.bmf_module_functor_call(self.p, dptr, C.int(isz), &osz, &is_done)
	var opkts []*Packet

	if bool(is_done) {
		return nil, nil
	}

	if int(osz) == 0 {
		return opkts, nil
	}

	if obuf == nil {
		return nil, errors.New(LastError())
	}

	defer C.free(unsafe.Pointer(obuf))
	oarr := (*[1 << 31]C.bmf_Packet)(unsafe.Pointer(obuf))[:osz:osz]
	for i := 0; i < int(osz); i++ {
		opkts = append(opkts, WrapPacket(oarr[i], true))
	}

	return opkts, nil
}

func (self *ModuleFunctor) Execute(inputs []*Packet, cleanup bool) (bool, error) {
	isz := len(inputs)
	dptr := (*C.bmf_Packet)(C.malloc(C.size_t(isz * 8)))
	defer C.free(unsafe.Pointer(dptr))
	darr := (*[1 << 31]C.bmf_Packet)(unsafe.Pointer(dptr))[:isz:isz]
	for i := 0; i < isz; i++ {
		darr[i] = inputs[i].Pointer()
	}

	is_done := C.bool(false)
	rc := C.bmf_module_functor_execute(self.p, dptr, C.int(isz), C.bool(cleanup), &is_done)
	if rc != 0 {
		return bool(is_done), errors.New(LastError())
	}

	return bool(is_done), nil
}

func (self *ModuleFunctor) Fetch(index int) ([]*Packet, error) {
	is_done := C.bool(false)
	osz := C.int(0)
	obuf := C.bmf_module_functor_fetch(self.p, C.int(index), &osz, &is_done)
	var opkts []*Packet

	if bool(is_done) {
		return nil, nil
	}

	if int(osz) == 0 {
		return opkts, nil
	}

	if obuf == nil {
		return nil, errors.New(LastError())
	}

	defer C.free(unsafe.Pointer(obuf))
	oarr := (*[1 << 31]C.bmf_Packet)(unsafe.Pointer(obuf))[:osz:osz]
	for i := 0; i < int(osz); i++ {
		opkts = append(opkts, WrapPacket(oarr[i], true))
	}

	return opkts, nil
}
