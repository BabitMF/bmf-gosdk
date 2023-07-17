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

type Ctype_hmp_Tensor = C.hmp_Tensor

type Tensor struct {
	p   C.hmp_Tensor
	own bool
}

func LastError() string {
	cstr := C.GoString(C.hmp_last_error())
	return cstr
}

func deleteTensor(o *Tensor) {
	if o != nil && o.own && o.p != nil {
		C.hmp_tensor_free(o.p)
		o.p = nil
	}
}

func WrapTensor(p C.hmp_Tensor, own bool) *Tensor {
	o := &Tensor{p: p, own: own}
	runtime.SetFinalizer(o, deleteTensor)
	return o
}

func Empty(shape []int64, dtype ScalarType, device string, pinned_memory bool) (*Tensor, error) {
	sptr := (*C.int64_t)(unsafe.Pointer(&shape[0]))
	dstr := C.CString(device)
	defer C.free(unsafe.Pointer(dstr))
	dptr := (*C.char)(unsafe.Pointer(dstr))

	tensor := C.hmp_tensor_empty(sptr, C.int(len(shape)),
		C.int(dtype), dptr, C.bool(pinned_memory))
	if tensor == nil {
		return nil, errors.New(LastError())
	}

	return WrapTensor(tensor, true), nil
}

func Arange(start int64, end int64, step int64, dtype ScalarType, device string, pinned_memory bool) (*Tensor, error) {
	dstr := C.CString(device)
	defer C.free(unsafe.Pointer(dstr))
	dptr := (*C.char)(unsafe.Pointer(dstr))

	tensor := C.hmp_tensor_arange(C.int64_t(start), C.int64_t(end), C.int64_t(step),
		C.int(dtype), dptr, C.bool(pinned_memory))
	if tensor == nil {
		return nil, errors.New(LastError())
	}

	return WrapTensor(tensor, true), nil
}

func (self *Tensor) Free() {
	runtime.SetFinalizer(self, nil) //remove finalizer
	deleteTensor(self)
}

func (self *Tensor) Pointer() C.hmp_Tensor {
	return self.p
}

func (self *Tensor) String() string {
	var c_str_size = int(0)
	c_str := C.hmp_tensor_stringfy(self.p, (*C.int)(unsafe.Pointer(&c_str_size)))
	return string((*[1 << 31]byte)(unsafe.Pointer(c_str))[:c_str_size:c_str_size])
}

func (self *Tensor) Fill(value *Scalar) {
	C.hmp_tensor_fill(self.p, value.p)
}

func (self *Tensor) Defined() bool {
	return bool(C.hmp_tensor_defined(self.p))
}

func (self *Tensor) Dim() int64 {
	return int64(C.hmp_tensor_dim(self.p))
}

func (self *Tensor) Size(dim int64) int64 {
	return int64(C.hmp_tensor_size(self.p, C.int64_t(dim)))
}

func (self *Tensor) Shape() []int64 {
	var ret []int64
	for i := int64(0); i < self.Dim(); i++ {
		ret = append(ret, self.Size(i))
	}
	return ret
}

func (self *Tensor) Stride(dim int64) int64 {
	return int64(C.hmp_tensor_stride(self.p, C.int64_t(dim)))
}

func (self *Tensor) Strides() []int64 {
	var ret []int64
	for i := int64(0); i < self.Dim(); i++ {
		ret = append(ret, self.Stride(i))
	}
	return ret
}

func (self *Tensor) Nitems() int64 {
	return int64(C.hmp_tensor_nitems(self.p))
}

func (self *Tensor) Nbytes() int64 {
	return int64(C.hmp_tensor_nbytes(self.p))
}

func (self *Tensor) Itemsize() int64 {
	return int64(C.hmp_tensor_itemsize(self.p))
}

func (self *Tensor) Dtype() ScalarType {
	return ScalarType(C.hmp_tensor_dtype(self.p))
}

func (self *Tensor) IsContiguous() bool {
	return bool(C.hmp_tensor_is_contiguous(self.p))
}

func (self *Tensor) DeviceType() DeviceType {
	return DeviceType(C.hmp_tensor_device_type(self.p))
}

func (self *Tensor) DeviceIndex() int {
	return int(C.hmp_tensor_device_index(self.p))
}

func (self *Tensor) Data() uintptr {
	return uintptr(C.hmp_tensor_data(self.p))
}

func (self *Tensor) Clone() (*Tensor, error) {
	p := C.hmp_tensor_clone(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapTensor(p, true), nil
}

func (self *Tensor) Alias() (*Tensor, error) {
	p := C.hmp_tensor_alias(self.p)
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapTensor(p, true), nil
}

func (self *Tensor) View(shape []int64) (*Tensor, error) {
	sptr := (*C.int64_t)(unsafe.Pointer(&shape[0]))
	p := C.hmp_tensor_view(self.p, sptr, C.int(len(shape)))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapTensor(p, true), nil
}

func (self *Tensor) Reshape(shape []int64) (*Tensor, error) {
	sptr := (*C.int64_t)(unsafe.Pointer(&shape[0]))
	p := C.hmp_tensor_reshape(self.p, sptr, C.int(len(shape)))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapTensor(p, true), nil
}

func (self *Tensor) Slice(dim int64, start int64, end int64, step int64) (*Tensor, error) {
	p := C.hmp_tensor_slice(self.p, C.int64_t(dim), C.int64_t(start),
		C.int64_t(end), C.int64_t(step))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapTensor(p, true), nil
}

func (self *Tensor) Select(dim int64, index int64) (*Tensor, error) {
	p := C.hmp_tensor_select(self.p, C.int64_t(dim), C.int64_t(index))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapTensor(p, true), nil
}

func (self *Tensor) Permute(dims []int64) (*Tensor, error) {
	sptr := (*C.int64_t)(unsafe.Pointer(&dims[0]))
	p := C.hmp_tensor_permute(self.p, sptr, C.int(len(dims)))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapTensor(p, true), nil
}

func (self *Tensor) Squeeze(dim int64) (*Tensor, error) {
	p := C.hmp_tensor_squeeze(self.p, C.int64_t(dim))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapTensor(p, true), nil
}

func (self *Tensor) Unsqueeze(dim int64) (*Tensor, error) {
	p := C.hmp_tensor_unsqueeze(self.p, C.int64_t(dim))
	if p == nil {
		return nil, errors.New(LastError())
	}
	return WrapTensor(p, true), nil
}

func (self *Tensor) ToDevice(device string, non_blocking bool) (*Tensor, error) {
	dstr := C.CString(device)
	defer C.free(unsafe.Pointer(dstr))
	dptr := (*C.char)(unsafe.Pointer(dstr))

	p := C.hmp_tensor_to_device(self.p, dptr, C.bool(non_blocking))
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapTensor(p, true), nil
}

func (self *Tensor) ToDtype(dtype ScalarType) (*Tensor, error) {
	p := C.hmp_tensor_to_dtype(self.p, C.int(dtype))
	if p == nil {
		return nil, errors.New(LastError())
	}

	return WrapTensor(p, true), nil
}

func (self *Tensor) CopyFrom(from *Tensor) {
	C.hmp_tensor_copy_from(self.p, from.p)
}
