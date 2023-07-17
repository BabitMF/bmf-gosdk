package hmp

import (
	"testing"
)

func compareArray(a []int64, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestTensorFactory(t *testing.T) {
	var devices = []DeviceType{
		kCPU,
		kCUDA,
	}

	for _, device_type := range devices {
		if DeviceCount(device_type) == 0 {
			continue
		}

		device := Device(device_type, 0)
		t.Logf("Test Tensor on device %s", device)

		// Empty
		{
			shape := []int64{2, 3, 4}
			tensor, err := Empty(shape, kFloat64, device, true)
			defer tensor.Free() //free manually

			if err != nil {
				t.Errorf("Tensor.Empty failed, %s", err.Error())
			} else {
				if !compareArray(tensor.Shape(), shape) {
					t.Errorf("Invalid shape")
				}

				if tensor.DeviceType() != device_type {
					t.Errorf("Invalid device type")
				}

				if tensor.DeviceIndex() != 0 {
					t.Errorf("Invalid device index")
				}

				if tensor.Dtype() != kFloat64 {
					t.Errorf("Invalid dtype")
				}

				if tensor.Nitems() != 24 {
					t.Errorf("Invalid nitems")
				}

				if tensor.Nbytes() != tensor.Nitems()*tensor.Itemsize() {
					t.Errorf("Invalid nbytes")
				}

				if !tensor.IsContiguous() {
					t.Errorf("Invalid is_contiguous")
				}

				if tensor.Data() == 0 {
					t.Errorf("Invalid data")
				}

				if tensor.Dim() != int64(len(shape)) {
					t.Errorf("Invalid dim")
				}

				if !tensor.Defined() {
					t.Errorf("Invalid defined")
				}

				str := tensor.String()
				if len(str) == 0 {
					t.Errorf("Invalid string")
				}

				strides := []int64{12, 4, 1}
				if !compareArray(tensor.Strides(), strides) {
					t.Errorf("Invalid strides")
				}

			}
		}

		//Arange & Fill
		{
			tensor, _ := Arange(1, 48, 2, kFloat64, device, false)
			//tensor is managed by GC

			value := NewFloatScalar(1.2)
			tensor.Fill(value)
		}

		// Shape ops
		{
			tensor, _ := Arange(0, 24, 1, kFloat64, device, false)

			// reshape
			tensor_reshaped, err := tensor.Reshape([]int64{2, 3, 4})
			if err != nil {
				t.Errorf("Tensor reshape failed")
			}
			if !compareArray(tensor_reshaped.Shape(), []int64{2, 3, 4}) {
				t.Errorf("Invalid shape")
			}

			// view & slice
			tensor_view, err1 := tensor.Reshape([]int64{2, 3, 4})
			if err1 != nil {
				t.Errorf("Tensor view failed")
			}
			if !compareArray(tensor_view.Shape(), []int64{2, 3, 4}) {
				t.Errorf("Invalid shape")
			}

			tensor_slice, err2 := tensor_reshaped.Slice(-2, 0, 3, 2)
			if err2 != nil {
				t.Errorf("Tensor slice failed")
			}
			if !compareArray(tensor_slice.Shape(), []int64{2, 2, 4}) {
				t.Errorf("Invalid slice shape")
			}

			if tensor_slice.IsContiguous() {
				t.Errorf("Invalid is_contiguous")
			}

			_, err3 := tensor_slice.View([]int64{16})
			if err3 == nil {
				t.Errorf("Tensor view Non-contiguous failed")
			}

			// select
			tensor_select, err4 := tensor_reshaped.Select(-2, 1)
			if err4 != nil {
				t.Errorf("Tensor select failed")
			}
			if !compareArray(tensor_select.Shape(), []int64{2, 4}) {
				t.Errorf("Invalid select shape")
			}

			// unsqueeze
			tensor_unsqueeze, err5 := tensor_reshaped.Unsqueeze(0)
			if err5 != nil {
				t.Errorf("Tensor unsqueeze failed")
			}
			if !compareArray(tensor_unsqueeze.Shape(), []int64{1, 2, 3, 4}) {
				t.Errorf("Invalid unsqueeze shape")
			}

			// squeeze
			tensor_squeeze, err6 := tensor_unsqueeze.Squeeze(0)
			if err6 != nil {
				t.Errorf("Tensor squeeze failed")
			}
			if !compareArray(tensor_squeeze.Shape(), []int64{2, 3, 4}) {
				t.Errorf("Invalid squeeze shape")
			}

			// permute
			tensor_permute, err7 := tensor_reshaped.Permute([]int64{2, 1, 0})
			if err7 != nil {
				t.Errorf("Tensor permute failed")
			}
			if !compareArray(tensor_permute.Shape(), []int64{4, 3, 2}) {
				t.Errorf("Invalid permute shape")
			}
		}

		// alias & clone & to & copy
		{
			shape := []int64{2, 3, 4}
			tensor, _ := Empty(shape, kFloat64, device, true)

			// alias
			tensor_alias, err0 := tensor.Alias()
			if err0 != nil {
				t.Errorf("Tensor alias failed")
			}
			if tensor_alias.Data() != tensor.Data() {
				t.Errorf("Invalid alias data")
			}

			tensor_alias.Reshape([]int64{-1})
			if !compareArray(tensor.Shape(), shape) {
				t.Errorf("Tensor alias failed")
			}

			//
			tensor_clone, err1 := tensor.Clone()
			if err1 != nil {
				t.Errorf("Tensor clone failed")
			}
			if tensor_clone.Data() == tensor.Data() {
				t.Errorf("Invalid clone data")
			}

			{
				tensor2, err2 := tensor.ToDevice("cpu", false)
				if err2 != nil {
					t.Errorf("Tensor to device failed")
				}

				tensor.CopyFrom(tensor2)
			}

			{
				_, err2 := tensor.ToDtype(kFloat64)
				if err2 != nil {
					t.Errorf("Tensor to dtype failed")
				}
			}

		}

	}

}
