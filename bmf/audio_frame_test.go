package bmf

import (
	"github.com/babitmf/bmf-gosdk/hmp"
	"testing"
)

func TestAudioFrame(t *testing.T) {
	// construct audio frame
	{
		af_frame, err := NewAudioFrame(2, 1, true, 0)
		defer af_frame.Free()
		if err != nil || !af_frame.Defined() {
			t.Errorf("create NewAudioFrame failed")
		} else {
			if af_frame.Dtype() != hmp.UInt8 {
				t.Errorf("AudioFrame with invalid dtype")
			}
			if af_frame.Planer() != true {
				t.Errorf("AudioFrame with invalid planer")
			}
			if af_frame.Nsamples() != 2 {
				t.Errorf("AudioFrame with invalid nsamples")
			}
			if af_frame.Nchannels() != 1 {
				t.Errorf("AudioFrame with invalid nchannels")
			}
			if af_frame.SampleRate() != 1 {
				t.Errorf("AudioFrame with invalid sample rate")
			}
			if af_frame.Nplanes() != 1 {
				t.Errorf("AudioFrame with invalid plane number")
			}
			af_frame.SetPts(1000)
			if af_frame.Pts() != 1000 {
				t.Errorf("AudioFrame with invalid pts")
			}
			af_frame.SetTimeBase(1, 100)
			num, den := af_frame.TimeBase()
			if num != 1 || den != 100 {
				t.Errorf("AudioFrame with invalid time base")
			}

			// construct audio frame from tensorList
			data, _ := af_frame.Planes()
			size := af_frame.Nplanes()
			af_frame2, err2 := NewAudioFrameFromData(data, size, 1, true)
			defer af_frame2.Free()
			if err2 != nil || !af_frame2.Defined() {
				t.Errorf("create NewAudioFrameFromData failed")
			} else {
				// copy props
				num2, den2 := af_frame2.TimeBase()
				if af_frame2.Pts() != 0 || num2 != -1 || den2 != -1 {
					t.Errorf("AudioFrame construct from data has incorrect pts or timebase")
				}
				af_frame2.CopyProps(af_frame)
				num2, den2 = af_frame2.TimeBase()
				if af_frame2.Pts() != 1000 || num2 != 1 || den2 != 100 {
					t.Errorf("AudioFrame CopyProps or SetPts or SetTimebase failed")
				}
			}

			// get exact tensor
			tensor, err3 := af_frame.Plane(0)
			defer tensor.Free()
			if err3 != nil || !af_frame2.Defined() || !tensor.Defined() {
				t.Errorf("cannot get exact plane from audio frame")
			}

			// private_get, private_attach and private_merge
			{
				jsonInfo := map[string]string{}
				jsonInfo["type"] = "audio_frame"
				af_frame.PrivateAttach(kJsonParam, jsonInfo)
				jsonData := map[string]string{}
				err := af_frame.PrivateGet(kJsonParam, &jsonData)
				if err != nil {
					t.Errorf("AudioFrame private get failed")
				}
				v, ok := jsonData["type"]
				if !ok || v != "audio_frame" {
					t.Errorf("AudioFrame get jsonparam data failed ")
				}

				afFrameMerge, _ := NewAudioFrame(2, 1, true, 0)
				afFrameMerge.PrivateMerge(af_frame)
				jsonData2 := map[string]string{}
				afFrameMerge.PrivateGet(kJsonParam, &jsonData2)
				v2, ok2 := jsonData2["type"]
				if !ok2 || v2 != "audio_frame" {
					t.Errorf("AudioFrame private merge failed ")
				}
			}
		}
	}
}
