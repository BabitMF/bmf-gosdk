package hmp

import (
	"testing"
)

func TestStream(t *testing.T) {
	if DeviceCount(kCUDA) == 0 {
		t.Skip("Skip Stream test as no cuda device found")
	} else {
		stream, err := NewStream(kCUDA, 0)
		if err != nil {
			t.Errorf("NewStream failed")
		}

		stream.Synchronize()

		{
			guard, _ := NewStreamGuard(stream)
			defer guard.Free()

			current, _ := CurrentStream(kCUDA)
			if current.Handle() != stream.Handle() {
				t.Errorf("Stream Guard failed")
			}

			data, _ := Empty([]int64{1 << 24}, kFloat32, "cuda:0", false)
			data.ToDevice("cpu", true)
			if stream.Query() {
				t.Errorf("Expect stream is busy")
			} else {
				stream.Synchronize()
				if !stream.Query() {
					t.Errorf("Expect stream is done")
				}
			}
		}

	}

}
