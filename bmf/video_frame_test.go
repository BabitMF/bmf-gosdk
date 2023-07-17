package bmf

import (
	"github.com/babitmf/bmf-gosdk/hmp"
	"testing"
)

func TestVideoFrame(t *testing.T) {
	// construct as image
	{
		vf_image, err := NewVideoFrameAsImage(1920, 1080, 3, hmp.NCHW, hmp.UInt8, "cpu", false)
		if err != nil || !vf_image.Defined() || !vf_image.IsImage() {
			t.Errorf("NewVideoFrameAsImage failed")
		} else {
			pix_info := hmp.NewPixelInfoV1(hmp.PF_YUV420P, hmp.CS_BT709, hmp.CR_MPEG)
			vf_frame, err1 := vf_image.ToFrame(pix_info)
			if err1 != nil || !vf_frame.Defined() || vf_frame.IsImage() {
				t.Errorf("VideoFrame.ToFrame failed with %v", err1)
			}

			if vf_image.Width() != 1920 || vf_image.Height() != 1080 || vf_image.Dtype() != hmp.UInt8 {
				t.Errorf("VideoFrame with invalid size or dtype")
			}

			image, _ := vf_image.Image()
			if image.Width() != vf_image.Width() || image.Height() != vf_image.Height() || image.Nchannels() != 3 {
				t.Errorf("VideoFrame with invalid Image")
			}

			frame, _ := vf_frame.Frame()
			if frame.Width() != vf_frame.Width() || frame.Height() != vf_frame.Height() || frame.Format() != hmp.PF_YUV420P {
				t.Errorf("VideoFrame with invalid Frame")
			}

			if vf_image.DeviceType() != hmp.CPU || vf_image.DeviceIndex() != 0 {
				t.Errorf("VideoFrame wiht invalid device info")
			}

			if hmp.DeviceCount(hmp.CUDA) > 0 {
				vf_image_cuda, err2 := vf_image.ToDevice("cuda", false)
				if err2 != nil || !vf_image_cuda.Defined() || vf_image_cuda.DeviceType() != hmp.CUDA {
					t.Errorf("VideoFrame ToDevice(cuda) failed")
				} else {
					vf_image_cpu, err3 := vf_image_cuda.ToDevice("cpu", false)
					if err3 != nil || !vf_image_cpu.Defined() || vf_image_cpu.DeviceType() != hmp.CPU {
						t.Errorf("VideoFrame ToDevice(cpu) failed")
					}
				}
			}

			{
				vf_image_f32, err4 := vf_image.ToDtype(hmp.Float32)
				if err4 != nil || vf_image_f32.Dtype() != hmp.Float32 {
					t.Errorf("VideoFrame ToDtype failed with error %v", err4)
				}

				// set & copy_props & get
				vf_image_f32.SetPts(1234567)
				vf_image_f32.SetTimeBase(12, 31)
				vf_image.CopyProps(vf_image_f32)
				num, den := vf_image.TimeBase()
				if vf_image.Pts() != 1234567 || num != 12 || den != 31 {
					t.Errorf("VideoFrame CopyProps or SetPts or SetTimebase failed")
				}
			}

			// private_get, private_attach and private_merge
			{
				jsonInfo := map[string]string{}
				jsonInfo["type"] = "video_frame"
				vf_image.PrivateAttach(kJsonParam, jsonInfo)
				jsonData := map[string]string{}
				err := vf_image.PrivateGet(kJsonParam, &jsonData)
				if err != nil {
					t.Errorf("VideoFrame private get failed")
				}
				v, ok := jsonData["type"]
				if !ok || v != "video_frame" {
					t.Errorf("VideoFrame get jsonparam data failed ")
				}

				vfMerge, _ := NewVideoFrameAsImage(1920, 1080, 3, hmp.NCHW, hmp.UInt8, "cpu", false)
				vfMerge.PrivateMerge(vf_image)
				jsonData2 := map[string]string{}
				vfMerge.PrivateGet(kJsonParam, &jsonData2)
				v2, ok2 := jsonData2["type"]
				if !ok2 || v2 != "video_frame" {
					t.Errorf("VideoFrame private merge failed ")
				}
			}

		}
	}

	// construct as frame
	{
		pix_info := hmp.NewPixelInfoV1(hmp.PF_YUV420P, hmp.CS_BT709, hmp.CR_MPEG)
		vf_frame, err := NewVideoFrameAsFrame(1920, 1080, pix_info, "cpu")
		if err != nil || !vf_frame.Defined() || vf_frame.IsImage() {
			t.Errorf("NewVideoFrameAsFrame failed")
		} else {
			vf_image, err1 := vf_frame.ToImage(hmp.NCHW, true)
			if err1 != nil || !vf_image.Defined() || !vf_image.IsImage() {
				t.Errorf("VideoFrame.ToImage failed with %v", err1)
			}
		}
	}

	// construct from frame
	{
		pix_info := hmp.NewPixelInfoV1(hmp.PF_YUV420P, hmp.CS_BT709, hmp.CR_MPEG)
		frame, _ := hmp.NewFrame(1920, 1080, pix_info, "cpu")
		vf, err := NewVideoFrameFromFrame(frame)
		defer vf.Free() //delete manually

		if err != nil || vf.IsImage() || vf.Width() != frame.Width() || vf.Height() != frame.Height() {
			t.Errorf("NewVideoFrameFromFrame failed")
		}
	}

	// construct from image
	{
		image, _ := hmp.NewImage(1920, 1080, 3, hmp.NCHW, hmp.Float32, "cpu", false)
		vf, err := NewVideoFrameFromImage(image)
		// delete by GC

		if err != nil || !vf.IsImage() || vf.Width() != image.Width() || vf.Height() != image.Height() {
			t.Errorf("NewVideoFrameFromImage failed")
		}
	}

	//Operation on Stream(Async)
	if hmp.DeviceCount(hmp.CUDA) > 0 {
		vf_image_cuda, _ := NewVideoFrameAsImage(8000, 4000, 3, hmp.NCHW, hmp.UInt8, "cuda", false)
		vf_image_cpu, _ := vf_image_cuda.ToDevice("cpu", true) //non_blocking = true
		current, _ := hmp.CurrentStream(hmp.CUDA)
		done := current.Query()
		vf_image_cpu.Record(true)
		done2 := current.Query()
		if vf_image_cpu.Ready() {
			t.Errorf("Expect copy is not done %v, %v", done, done2)
		} else {
			vf_image_cpu.Synchronize()
			if !vf_image_cpu.Ready() {
				t.Errorf("Expect copy is done")
			}
		}

	}

}
