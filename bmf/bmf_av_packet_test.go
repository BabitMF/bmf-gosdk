package bmf

import (
	"testing"
)

func TestBMFAVPacket(t *testing.T) {
	// construct bmf_av_packet
	{
		bmf_av_pkt, err := NewBMFAVPacket(1, 1)
		defer bmf_av_pkt.Free()
		if err != nil || !bmf_av_pkt.Defined() {
			t.Errorf("create BMFAVPacket failed")
		} else {
			if bmf_av_pkt.Nbytes() != 1 {
				t.Errorf("BMFAVPacket with invalid nbytes")
			}

			bmf_av_pkt.SetPts(1000)
			if bmf_av_pkt.Pts() != 1000 {
				t.Errorf("BMFAVPacket with invalid pts")
			}
			bmf_av_pkt.SetTimeBase(1, 100)
			num, den := bmf_av_pkt.TimeBase()
			if num != 1 || den != 100 {
				t.Errorf("BMFAVPacket with invalid time base")
			}

			// construct bmf_av_packet from tensor
			data, _ := bmf_av_pkt.Data()
			bmf_av_pkt2, err2 := NewBMFAVPacketFromData(data)
			defer bmf_av_pkt2.Free()
			if err2 != nil || !bmf_av_pkt2.Defined() {
				t.Errorf("create NewBMFAVPacketFromData failed")
			} else {
				// copy props
				num2, den2 := bmf_av_pkt2.TimeBase()
				if bmf_av_pkt2.Pts() != 0 || num2 != -1 || den2 != -1 {
					t.Errorf("BMFAVPacket construct from data has incorrect pts or timebase")
				}
				bmf_av_pkt2.CopyProps(bmf_av_pkt)
				num2, den2 = bmf_av_pkt2.TimeBase()
				if bmf_av_pkt2.Pts() != 1000 || num2 != 1 || den2 != 100 {
					t.Errorf("BMFAVPacket CopyProps or SetPts or SetTimebase failed")
				}
			}

			// private_get, private_attach and private_merge
			{
				jsonInfo := map[string]string{}
				jsonInfo["type"] = "bmf_av_packet"
				bmf_av_pkt.PrivateAttach(kJsonParam, jsonInfo)
				jsonData := map[string]string{}
				err := bmf_av_pkt.PrivateGet(kJsonParam, &jsonData)
				if err != nil {
					t.Errorf("BMFAVPacket private get failed")
				}
				v, ok := jsonData["type"]
				if !ok || v != "bmf_av_packet" {
					t.Errorf("BMFAVPacket get jsonparam data failed ")
				}

				pktMerge, _ := NewBMFAVPacket(1, 1)
				pktMerge.PrivateMerge(bmf_av_pkt)
				jsonData2 := map[string]string{}
				pktMerge.PrivateGet(kJsonParam, &jsonData2)
				v2, ok2 := jsonData2["type"]
				if !ok2 || v2 != "bmf_av_packet" {
					t.Errorf("BMFAVPacket private merge failed ")
				}
			}
		}
	}
}
