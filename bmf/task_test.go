package bmf

import (
	"reflect"
	"testing"

	"github.com/babitmf/bmf-gosdk/hmp"
)

func TestTaskConstruct(t *testing.T) {
	{
		iids := []int32{1, 8, 3, 4}
		oids := []int32{3, 2, 1}

		task, err := NewTask(1, iids, oids)
		if err != nil {
			t.Errorf("Create task failed")
		}

		task.SetTimestamp(42)
		if task.Timestamp() != 42 {
			t.Errorf("Task SetTimestamp failed")
		}

		input_stream_ids := task.GetInputStreamIds()
		if len(input_stream_ids) != len(iids) || reflect.DeepEqual(input_stream_ids, iids) {
			t.Errorf("Task.GetInputStreamIds failed, %v", input_stream_ids)
		}

		output_stream_ids := task.GetOutputStreamIds()
		if len(output_stream_ids) != len(oids) || reflect.DeepEqual(output_stream_ids, oids) {
			t.Errorf("Task.GetOutputStreamIds failed, %v", output_stream_ids)
		}
	}
}

func TestTaskConstructZeroStreams(t *testing.T) {
	{
		iids := []int32{}
		oids := []int32{}

		task, err := NewTask(1, iids, oids)
		if err != nil {
			t.Errorf("Create task failed")
		}

		o_iids := task.GetInputStreamIds()
		o_oids := task.GetOutputStreamIds()
		if len(o_iids) != 0 || len(o_oids) != 0 {
			t.Errorf("Expect no stream in inputs and outputs")
		}
	}
}

func TestTaskFillPop(t *testing.T) {
	{
		iids := []int32{1, 4}
		oids := []int32{3, 2, 1}

		task, err := NewTask(1, iids, oids)
		if err != nil {
			t.Errorf("Create task failed")
		}

		//
		for i := 0; i < len(iids)*2; i++ {
			vf, err := NewVideoFrameAsImage(1920, 1080, 3, hmp.NCHW, hmp.UInt8, "cpu", false)
			if err != nil || !vf.Defined() || !vf.IsImage() {
				t.Fatalf("NewVideoFrameAsImage failed")
			}

			pkt, _ := NewPacketFromVideoFrame(vf)
			if err != nil {
				t.Fatalf("NewPaketFromVideoFrame failed")
			}
			pkt.SetTimestamp(int64(i))

			istream := iids[i%len(iids)]
			task.FillInputPacket(istream, pkt)
		}
		for i := 0; i < len(iids); i++ {
			pkt := GenerateEofPacket()
			task.FillInputPacket(iids[i], pkt)
		}

		//
		for i := 0; i < len(oids)*2; i++ {
			vf, err := NewVideoFrameAsImage(1280, 720, 3, hmp.NCHW, hmp.UInt8, "cpu", false)
			if err != nil || !vf.Defined() || !vf.IsImage() {
				t.Fatalf("NewVideoFrameAsImage failed")
			}

			pkt, _ := NewPacketFromVideoFrame(vf)
			if err != nil {
				t.Fatalf("NewPaketFromVideoFrame failed")
			}
			pkt.SetTimestamp(int64(i))

			ostream := oids[i%len(oids)]
			task.FillOutputPacket(ostream, pkt)
		}
		for i := 0; i < len(oids); i++ {
			task.FillOutputPacket(oids[i], GenerateEofPacket())
		}

		//check input queue
		for i := 0; i < len(iids)*2; i++ {
			istream := iids[i%len(iids)]
			pkt, err := task.PopPacketFromInputQueue(istream)
			if err != nil {
				t.Fatalf("PopPacketFromInputQueue failed")
			}

			vf, err := pkt.GetVideoFrame()
			if err != nil {
				t.Fatalf("GetVideFrame from packet failed")
			}

			if vf.Width() != 1920 || vf.Height() != 1080 {
				t.Errorf("Invalid VideoFrame from input queue")
			}

			if pkt.Timestamp() != int64(i) {
				t.Errorf("Invalid timestamp")
			}
		}
		for i := 0; i < len(iids); i++ {
			pkt, err := task.PopPacketFromInputQueue(iids[i])
			if err != nil || pkt.Timestamp() != EOF {
				t.Fatalf("Expect to get EOF packet %v", err)
			}
		}
		for i := 0; i < len(iids); i++ {
			istream := iids[i%len(iids)]
			_, err := task.PopPacketFromInputQueue(istream)
			if err == nil {
				t.Fatalf("Expect fail when pop from empty input queue")
			}
		}

		//check output queue
		for i := 0; i < len(oids)*2; i++ {
			ostream := oids[i%len(oids)]
			pkt, err := task.PopPacketFromOutQueue(ostream)
			if err != nil {
				t.Fatalf("PopPacketFromOutputQueue failed")
			}

			vf, err := pkt.GetVideoFrame()
			if err != nil {
				t.Fatalf("GetVideFrame from packet failed")
			}

			if vf.Width() != 1280 || vf.Height() != 720 {
				t.Errorf("Invalid VideoFrame from output queue")
			}

			if pkt.Timestamp() != int64(i) {
				t.Errorf("Invalid timestamp")
			}
		}
		for i := 0; i < len(oids); i++ {
			pkt, err := task.PopPacketFromOutQueue(oids[i])
			if err != nil || pkt.Timestamp() != EOF {
				t.Fatalf("Expect to get EOF packet %v", err)
			}
		}
		for i := 0; i < len(oids); i++ {
			ostream := oids[i%len(oids)]
			_, err := task.PopPacketFromOutQueue(ostream)
			if err == nil {
				t.Fatalf("Expect fail when pop from empty output queue")
			}
		}

	}
}
