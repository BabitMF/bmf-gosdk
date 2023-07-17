package custom

import (
	"testing"
)

func TestVideoFrameConvert(t *testing.T) {
	{
		// new AVFrame
		af := NewAVFrameWithVideo(1920, 1080, 0)
		if af == nil {
			t.Fatalf("New AVFrameWithVideo failed")
		}

		// convert AVFrame to VideoFrame
		vf, err := BfmVfFromAVFrame(af.p)
		if err != nil {
			t.Fatalf("BmfVfFromAVFrame failed, %v", err)
		}

		if vf.Width() != 1920 || vf.Height() != 1080 {
			t.Fatalf("Invalid image size %d %d", vf.Width(), vf.Height())
		}

		// convert VideoFrame to AVFrame
		af_copy, err1 := BfmVfToAVFrame(vf)
		if err1 != nil {
			t.Fatalf("BmfVfToAVFrame failed, %v", err1)
		}
		if int(af_copy.width) != 1920 || int(af_copy.height) != 1080 {
			t.Fatalf("Invalid AVFrame size, %d %d", int(af_copy.width), int(af_copy.height))
		}

	}
}
