package hmp

import (
	"testing"
)

func TestColorModel(t *testing.T) {
	cm := NewColorModel(CS_BT709, CR_MPEG, CP_BT709, CTC_BT709)

	if cm.Space() != CS_BT709 {
		t.Errorf("Invalid ColorSpace")
	}

	if cm.Range() != CR_MPEG {
		t.Errorf("Invalid ColorRange")
	}

	if cm.Primaries() != CP_BT709 {
		t.Errorf("Invalid Primaries")
	}

	if cm.TransferCharacteristic() != CTC_BT709 {
		t.Errorf("Invalid CTC")
	}
}

func TestPixelInfo(t *testing.T) {
	pix_info0 := NewPixelInfo(PF_YUV420P, NewColorModel(CS_BT709, CR_MPEG, CP_BT709, CTC_BT709))

	if pix_info0.Format() != PF_YUV420P {
		t.Errorf("Invalid PixelFormat")
	}

	if pix_info0.Space() != CS_BT709 {
		t.Errorf("Invalid ColorSpace")
	}

	if pix_info0.Range() != CR_MPEG {
		t.Errorf("Invalid ColorRange")
	}

	//
	pix_info1 := NewPixelInfoV1(PF_YUV420P, CS_BT709, CR_MPEG)
	if pix_info1.Format() != PF_YUV420P {
		t.Errorf("Invalid PixelFormat")
	}

	if pix_info1.Space() != CS_BT709 {
		t.Errorf("Invalid ColorSpace")
	}

	if pix_info1.Range() != CR_MPEG {
		t.Errorf("Invalid ColorRange")
	}

	if pix_info1.Primaries() != CP_UNSPECIFIED {
		t.Errorf("Invalid ColorPrimaries")
	}

	if pix_info1.TransferCharacteristic() != CTC_UNSPECIFIED {
		t.Errorf("Invalid CTC")
	}

	//
	pix_info2 := NewPixelInfoV2(PF_YUV420P, CP_BT709, CTC_BT709)
	if pix_info2.Format() != PF_YUV420P {
		t.Errorf("Invalid PixelFormat")
	}

	if pix_info2.Space() != CS_UNSPECIFIED {
		t.Errorf("Invalid ColorSpace")
	}

	if pix_info2.Range() != CR_UNSPECIFIED {
		t.Errorf("Invalid ColorRange")
	}

	if pix_info2.Primaries() != CP_BT709 {
		t.Errorf("Invalid ColorPrimaries")
	}

	if pix_info2.TransferCharacteristic() != CTC_BT709 {
		t.Errorf("Invalid CTC")
	}

}
