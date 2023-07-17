package bmf

import (
	"errors"
	"testing"

	"github.com/babitmf/bmf-gosdk/hmp"
)

// A type safe wrapper of PassThrough module
func (self *ModuleFunctor) callI2O2(vf0 *VideoFrame, vf1 *VideoFrame) (*VideoFrame, *VideoFrame, error) {
	// Step 0 -> convert input to pkts
	ipkt0, _ := NewPacketFromVideoFrame(vf0)
	ipkt1, _ := NewPacketFromVideoFrame(vf1)

	// Step 1 -> Invoke module
	opkts, err := self.Call([]*Packet{ipkt0, ipkt1})
	if err != nil {
		return nil, nil, err
	}

	// Step 2 -> Convert output packets back to concrete types
	ovf0, _ := opkts[0].GetVideoFrame()
	ovf1, _ := opkts[1].GetVideoFrame()
	return ovf0, ovf1, nil
}

func (self *ModuleFunctor) callI0O0() error {
	opkts, err := self.Call([]*Packet{})
	if err != nil {
		return err
	}

	if opkts != nil {
		return errors.New("Expect nothing returned")
	}

	return nil
}

func TestModuleFunctor(t *testing.T) {
	// 2-inputs, 2-outputs
	{
		mf, err0 := NewModuleFunctor("PassThrough", "go", "../example/go_pass_through.so", "", nil, 2, 2)
		if mf == nil {
			t.Errorf("NewModuleFunctor failed, %v", err0)
		} else {
			ivf0, _ := NewVideoFrameAsImage(1920, 1080, 3, hmp.NCHW, hmp.UInt8, "cpu", false)
			ivf1, _ := NewVideoFrameAsImage(1280, 720, 3, hmp.NCHW, hmp.UInt8, "cpu", false)
			ovf0, ovf1, err1 := mf.callI2O2(ivf0, ivf1)
			if ovf0 == nil || ovf1 == nil {
				t.Errorf("CallPassThrough failed, %v", err1)
			}

			if ovf0.Width() != 1920 || ovf1.Width() != 1280 {
				t.Errorf("Invalid result from CallPassThrough")
			}
		}
	}

	{
		// 0-inputs, 0-outputs
		mf, err0 := NewModuleFunctor("PassThrough", "go", "../example/go_pass_through.so", "", nil, 0, 0)
		if mf == nil {
			t.Errorf("NewModuleFunctor failed, %v", err0)
		} else {
			err1 := mf.callI0O0()
			if err1 != nil {
				t.Errorf("%v", err1)
			}
		}
	}
}

func TestModuleFunctorProcessDone(t *testing.T) {
	{
		mf, err0 := NewModuleFunctor("PassThrough", "go", "../example/go_pass_through.so", "", nil, 2, 2)
		if mf == nil {
			t.Errorf("NewModuleFunctor failed, %v", err0)
		} else {
			ipkts := []*Packet{GenerateEofPacket(), GenerateEofPacket()}
			opkts, err := mf.Call(ipkts)
			if err != nil || opkts != nil {
				t.Errorf("Expect no output and no error")
			}
		}
	}
}
