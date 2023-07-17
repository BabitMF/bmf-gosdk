package bmf

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/babitmf/bmf-gosdk/hmp"
)

type MockModuleOption struct {
	Info int32
}

type MockModule struct {
	option MockModuleOption
	nodeId int32

	initError        error
	resetError       error
	closeError       error
	moduleInfoResult map[string]string
}

func (self *MockModule) Process(task *Task) error {
	iids := task.GetInputStreamIds()
	oids := task.GetOutputStreamIds()

	gotEof := false
	for i, iid := range iids {
		for pkt, err := task.PopPacketFromInputQueue(iid); err == nil; {
			defer pkt.Free()
			if ok := task.FillOutputPacket(oids[i], pkt); !ok {
				return errors.New("Fill output queue failed")
			}

			if pkt.Timestamp() == EOF {
				gotEof = true
			}

			pkt, err = task.PopPacketFromInputQueue(iid)
		}
	}

	if gotEof {
		task.SetTimestamp(DONE)
	}

	return nil
}

func (self *MockModule) Init() error {
	return self.initError
}

func (self *MockModule) Reset() error {
	return self.resetError
}

func (self *MockModule) Close() error {
	return self.closeError
}

func (self *MockModule) GetModuleInfo() (interface{}, error) {
	return self.moduleInfoResult, nil
}

func (self *MockModule) NeedHungryCheck(istreamId int32) (bool, error) {
	return true, nil
}

func (self *MockModule) IsHungry(istreamId int32) (bool, error) {
	return true, nil
}

func (self *MockModule) IsInfinity() (bool, error) {
	return true, nil
}

func NewMockModule(nodeId int32, option []byte) (Module, error) {
	m := &MockModule{}
	err := json.Unmarshal(option, &m.option)
	if err != nil {
		return nil, err
	}
	m.nodeId = nodeId

	return m, nil
}

func TestModule(t *testing.T) {
	RegisterModuleConstructor("MockModule", NewMockModule)

	sdkVersion := SdkVersion()
	if len(sdkVersion) == 0 {
		t.Errorf("Invalid BMF SDK Version")
	}

	//normal construction
	{
		option := MockModuleOption{42}
		optionStr, _ := json.Marshal(option)

		id := ModuleConstructGo("MockModule", 1, string(optionStr))
		if id < 0 {
			t.Errorf("Expect construct MockModule success")
		}

		if m, ok := GetModuleInstance(id); ok {
			if m.(Module) == nil {
				t.Errorf("convert to Module failed")
			}

			module := m.(*MockModule)
			if module.nodeId != 1 {
				t.Errorf("Expect nodeId == 42, got %d", module.nodeId)
			}

			if module.option.Info != 42 {
				t.Errorf("Invalid module option value, expect 42, got %d", module.option.Info)
			}

			//Init
			{
				module.initError = errors.New("initError")
				err := ModuleInitGo(id)
				if err == nil || err.Error() != "initError" {
					t.Errorf("Expect init failed")
				}

				module.initError = nil
				err = ModuleInitGo(id)
				if err != nil {
					t.Errorf("Expect init success")
				}
			}

			//Reset
			{
				module.resetError = errors.New("resetError")
				err := ModuleResetGo(id)
				if err == nil || err.Error() != "resetError" {
					t.Errorf("Expect reset failed")
				}

				module.resetError = nil
				err = ModuleResetGo(id)
				if err != nil {
					t.Errorf("Expect reset success")
				}
			}

			//
			{
				if ModuleNeedHungryCheck(id, 0) != true {
					t.Errorf("Expect NeedHungryCheck == true")
				}

				if ModuleIsHungry(id, 0) != true {
					t.Errorf("Expect IsHungry == true")
				}

				if ModuleIsInfinity(id) != true {
					t.Errorf("Expect IsInfinity == true")
				}
			}

			//ModuleInfo
			{
				module.moduleInfoResult = map[string]string{
					"name": "MockModule",
				}

				info := ModuleGetInfoGo(id).(*map[string]string)
				if info == nil || (*info)["name"] != "MockModule" {
					t.Errorf("Invalid result of GetModuleInfo")
				}

			}

			// process
			{
				task, _ := NewTask(module.nodeId, []int32{1, 2}, []int32{1, 2})
				vf, _ := NewVideoFrameAsImage(1920, 1080, 3, hmp.NCHW, hmp.UInt8, "cpu", false)
				for i := 0; i < 2; i++ {
					pkt, _ := NewPacketFromVideoFrame(vf)
					task.FillInputPacket(1, pkt)
					task.FillInputPacket(2, pkt)
				}
				task.FillInputPacket(1, GenerateEofPacket())
				task.FillInputPacket(2, GenerateEofPacket())

				err := ModuleProcessGo(id, task)
				if err != nil {
					t.Errorf("Module process failed")
				} else {
					if task.Timestamp() != DONE {
						t.Errorf("Expect get DONE task")
					}
				}
			}

			//Close
			{
				module.closeError = errors.New("closeError")
				err := ModuleCloseGo(id)
				if err == nil || err.Error() != "closeError" {
					t.Errorf("Expect close failed")
				}
			}

		} else {
			t.Errorf("No module instance found for id %d", id)
		}
	}

	//construct with invalid option
	{
		option := "Just a string"
		id := ModuleConstructGo("MockModule", 1, option)
		if id != -1 {
			t.Errorf("Expect construct MockModule failed with errno -1")
		}
	}

	//construct with unkonwn cls
	{
		option := MockModuleOption{42}
		optionStr, _ := json.Marshal(option)

		id := ModuleConstructGo("UknownModule", 1, string(optionStr))
		if id != -2 {
			t.Errorf("Expect construct MockModule failed with errno -2")
		}
	}

}
