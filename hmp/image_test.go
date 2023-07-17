package hmp

import (
	"testing"
)

func TestFrame(t *testing.T) {
	pix_info := NewPixelInfoV1(PF_YUV420P, CS_BT709, CR_MPEG)

	// construct and attributes
	{
		frame0, err0 := NewFrame(1920, 1080, pix_info, "cpu")
		defer frame0.Free()

		if err0 != nil {
			t.Errorf("NewFrame failed")
		} else {
			if !frame0.Defined() {
				t.Errorf("Invalid Frame")
			}

			if frame0.PixInfo().Format() != PF_YUV420P {
				t.Errorf("Invalid PixelInfo")
			}

			if frame0.Format() != PF_YUV420P {
				t.Errorf("Invalid PixelFormat")
			}

			if frame0.Width() != 1920 {
				t.Errorf("Invalid Width")
			}

			if frame0.Height() != 1080 {
				t.Errorf("Invalid Height")
			}

			if frame0.Dtype() != kUInt8 {
				t.Errorf("Invalid dtype")
			}

			if frame0.DeviceType() != kCPU {
				t.Errorf("Invalid device type")
			}

			if frame0.DeviceIndex() != 0 {
				t.Errorf("Invalid device index")
			}

			if frame0.Nplanes() != 3 {
				t.Errorf("Invalid nplanes")
			}

			Y := frame0.Plane(0)
			U := frame0.Plane(1)
			if Y.Size(0) != 1080 || Y.Size(1) != 1920 {
				t.Errorf("Invalid Y plane")
			}
			if U.Size(0) != 540 || U.Size(1) != 960 {
				t.Errorf("Invalid U plane")
			}

			if DeviceCount(kCUDA) > 0 {
				frame1, _ := frame0.ToDevice("cuda:0", false)
				if frame1.DeviceType() != kCUDA {
					t.Errorf("Frame.ToDevice failed")
				}

				// ensure no coredump
				frame1.CopyFrom(frame0)
			}

			_, err1 := frame0.Clone()
			if err1 != nil {
				t.Errorf("Frame.Clone failed")
			}

			frame1, err2 := frame0.Crop(200, 300, 400, 500)
			if err2 != nil {
				t.Errorf("Frame.Crop failed")
			} else {
				if frame1.Width() != 400 || frame1.Height() != 500 {
					t.Errorf("Frame.Crop failed")
				}
			}

			image, err := frame0.ToImage(kNHWC)
			if err != nil {
				t.Errorf("ToImage failed")
			} else {
				if image.Width() != 1920 || image.Height() != 1080 {
					t.Errorf("ToImage failed")
				}
			}

			str := frame0.String()
			if len(str) == 0 {
				t.Errorf("Frame.String failed")
			}
		}

	}

	//
	{
		Y, _ := Empty([]int64{1080, 1920}, kUInt8, "cpu", false)
		U, _ := Empty([]int64{540, 960}, kUInt8, "cpu", false)
		V, _ := Empty([]int64{540, 960}, kUInt8, "cpu", false)
		data := []*Tensor{Y, U, V}

		frame0, err0 := NewFrameFromData(data, pix_info)
		if err0 != nil {
			t.Errorf("NewFrameFromData failed")
		} else {
			if frame0.Width() != 1920 || frame0.Height() != 1080 {
				t.Errorf("NewFrameFromData with invalid size")
			}
		}

		// invalid size
		_, err1 := NewFrameFromDataV1(data, 1111, 1111, pix_info)
		if err1 == nil {
			t.Errorf("NewFrameFromDataV1 failed")
		}

		frame2, err2 := NewFrameFromDataV1(data, 1920, 1080, pix_info)
		if err2 != nil {
			t.Errorf("NewFrameFromDataV1 failed")
		} else {
			if frame2.Width() != 1920 || frame2.Height() != 1080 {
				t.Errorf("NewFrameFromDataV1: with invalid size")
			}

			if frame2.Format() != PF_YUV420P {
				t.Errorf("NewFrameFromDataV1: invalid format")
			}
		}
	}

}

func TestImage(t *testing.T) {
	{
		image, err := NewImage(1920, 1080, 3, kNCHW, kFloat32, "cpu", false)
		defer image.Free()

		if err != nil || !image.Defined() {
			t.Errorf("Image: Invalid defined")
		}

		if image.Format() != kNCHW {
			t.Errorf("Image: Invalid format")
		}

		cm := NewColorModel(CS_BT709, CR_MPEG, CP_BT709, CTC_BT709)
		image.SetColorModel(cm)
		if image.ColorModel().Space() != CS_BT709 {
			t.Errorf("Image: Invalid ColorModel")
		}

		if image.Wdim() != 2 || image.Hdim() != 1 || image.Cdim() != 0 {
			t.Errorf("Image: Invalid channel format")
		}

		if image.Width() != 1920 || image.Height() != 1080 || image.Nchannels() != 3 {
			t.Errorf("Image: Invalid size")
		}

		if image.Dtype() != kFloat32 {
			t.Errorf("Image: Invalid dtype")
		}

		if image.DeviceType() != kCPU || image.DeviceIndex() != 0 {
			t.Errorf("Image: Invalid device info")
		}

		im_data := image.Data()
		if im_data.Size(0) != 3 || im_data.Size(1) != 1080 || im_data.Size(2) != 1920 {
			t.Errorf("Image: Invalid data")
		}

		if DeviceCount(kCUDA) > 0 {
			image1, err1 := image.ToDevice(Device(kCUDA, 0), false)
			if err1 != nil || image1.DeviceType() != kCUDA {
				t.Errorf("Image: to device failed")
			}
		}

		{
			image1, err1 := image.ToDtype(kUInt8)
			if err1 != nil || image1.Dtype() != kUInt8 {
				t.Errorf("Image: to dtype failed")
			}

			//
			image.CopyFrom(image1)
		}

		{
			image1, err1 := image.Clone()
			if err1 != nil || !image1.Defined() {
				t.Errorf("Image: clone failed")
			}
		}

		{
			image1, err1 := image.Crop(200, 300, 400, 500)
			if err1 != nil || image1.Width() != 400 || image1.Height() != 500 {
				t.Errorf("Image: crop failed")
			}
		}

		{
			image1, err1 := image.Select(1)
			if err1 != nil || image1.Nchannels() != 1 {
				t.Errorf("Image: select failed")
			}
		}

		str := image.String()
		if len(str) == 0 {
			t.Errorf("Image: String failed")
		}
	}

	{
		data, _ := Empty([]int64{1080, 1920, 4}, kUInt8, "cpu", false)
		image, err := NewImageFromData(data, kNHWC)
		if err != nil || !image.Defined() || image.Nchannels() != 4 {
			t.Errorf("Image: NewImageFromData failed")
		}
	}

	{
		cm := NewColorModel(CS_BT709, CR_MPEG, CP_BT709, CTC_BT709)
		data, _ := Empty([]int64{1080, 1920, 4}, kUInt8, "cpu", false)
		image, err := NewImageFromDataV1(data, kNHWC, cm)
		if err != nil || !image.Defined() || image.ColorModel().Space() != CS_BT709 {
			t.Errorf("Image: NewImageFromDataV1 failed")
		}
	}

}
