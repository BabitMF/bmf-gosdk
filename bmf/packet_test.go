package bmf

import (
	"testing"

	"github.com/babitmf/bmf-gosdk/hmp"
)

func TestPacket(t *testing.T) {
	//eos & eof packet
	{
		eos_pkt := GenerateEosPacket()
		if eos_pkt.Timestamp() != EOS {
			t.Errorf("Incorrect timestamp of EOS packet")
		}

		eof_pkt := GenerateEofPacket()
		if eof_pkt.Timestamp() != EOF {
			t.Errorf("Incorrect timestamp of EOF packet")
		}
	}

	// VideoFrame
	{
		vf, err := NewVideoFrameAsImage(1920, 1080, 3, hmp.NCHW, hmp.UInt8, "cpu", false)
		if err != nil || !vf.Defined() || !vf.IsImage() {
			t.Fatalf("NewVideoFrameAsImage failed")
		}

		pkt, err := NewPacketFromVideoFrame(vf)
		if err != nil || !pkt.Defined() {
			t.Fatalf("NewPacketFromVideoFrame failed")
		}

		if !pkt.IsVideoFrame() {
			t.Fatalf("Packet video frame check failed")
		}

		vf1, err1 := pkt.GetVideoFrame()
		if err1 != nil || vf1.Width() != 1920 || vf1.Height() != 1080 {
			t.Errorf("Get video frame from Packet failed")
		}

		pkt.SetTimestamp(0xdeadbeef)
		if pkt.Timestamp() != 0xdeadbeef {
			t.Errorf("Set/Get timestamp to/from packet failed")
		}
	}

	// AudioFrame
	{
		af, err := NewAudioFrame(2, 1, true, 0)
		if err != nil || !af.Defined() {
			t.Errorf("create NewAudioFrame failed")
		}

		pkt, err := NewPacketFromAudioFrame(af)
		if err != nil || !pkt.Defined() {
			t.Fatalf("NewPacketFromAudioFrame failed")
		}

		if !pkt.IsAudioFrame() {
			t.Fatalf("Packet audio frame check failed")
		}

		af1, err1 := pkt.GetAudioFrame()
		if err1 != nil || af1.Dtype() != hmp.UInt8 || af1.Nsamples() != 2 {
			t.Errorf("Get audio frame from Packet failed")
		}
	}

	// BMFAVPacket
	{
		bmfAvPkt, err := NewBMFAVPacket(1, 1)
		if err != nil || !bmfAvPkt.Defined() {
			t.Errorf("create BMFAVPacket failed")
		}

		pkt, err := NewPacketFromBMFAVPacket(bmfAvPkt)
		if err != nil || !pkt.Defined() {
			t.Fatalf("NewPacketFromBMFAVPacket failed")
		}

		if !pkt.IsBMFAVPacket() {
			t.Fatalf("Packet bmf_av_packet check failed")
		}

		bmfAvPkt1, err1 := pkt.GetBMFAVPacket()
		if err1 != nil || bmfAvPkt1.Nbytes() != 1 {
			t.Errorf("Get bmf_av_packet from Packet failed")
		}
	}

	// JsonParam
	{
        data := map[string]int{"apple": 5, "lettuce": 7}
		pkt, err0 := NewPacketFromJsonParam(&data)
		if err0 != nil {
			t.Errorf("New packet with struct failed, %v", err0)
		} else {
            if !pkt.IsJsonParam(){
                t.Errorf("Expect struct data in packet")
            } else {
			    var p_data map[string]int
                err1 := pkt.GetJsonParam(&p_data)
                if err1 != nil{
                    t.Errorf("Parse struct from packet failed, %v", err1)
                }

                if p_data["apple"] != data["apple"] || p_data["lettuce"] != data["lettuce"] {
                    t.Errorf("Invalid value parsed from packet")
                }
            }
		}

	}

}
