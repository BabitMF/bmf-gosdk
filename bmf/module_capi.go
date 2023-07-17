package bmf

/*
#include <stdint.h>
#include <stdlib.h>
#include <bmf/sdk/bmf_capi.h>
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

type ModuleConstructor func(nodeId int32, option []byte) (Module, error)
type ModuleInfoRegister func(ModuleInfo)

var G struct {
	sync.Mutex

	//Module instances
	modules map[int32]interface{}
	alloced int32

	//Module constructors
	constructors map[string]interface{}

	//Module Info registers
	registers map[string]ModuleInfoRegister
}

//init globals
func init() {
	G.Lock()
	defer G.Unlock()

	G.constructors = make(map[string]interface{})
	G.modules = make(map[int32]interface{})
	G.alloced = 0
	G.registers = make(map[string]ModuleInfoRegister)
}

func RegisterModuleConstructor(cls string, constructor ModuleConstructor, register ModuleInfoRegister) {
	G.Lock()
	defer G.Unlock()
	G.constructors[cls] = constructor
	if register != nil {
		G.registers[cls] = register
	}
}

func GetModuleConstructor(cls string) (interface{}, bool) {
	G.Lock()
	defer G.Unlock()
	constructor, ok := G.constructors[cls]
	return constructor, ok
}

func getModuleInfoRegister(cls string) (interface{}, bool) {
	G.Lock()
	defer G.Unlock()
	register, ok := G.registers[cls]
	return register, ok
}

func RegisterModuleInstance(instance interface{}) int32 {
	G.Lock()
	defer G.Unlock()
	id := G.alloced
	G.alloced += 1
	G.modules[id] = instance
	return id
}

func UnregisterModuleInstance(id int32) {
	G.Lock()
	defer G.Unlock()

	if _, ok := G.modules[id]; ok {
		delete(G.modules, id)
	}
}

func GetModuleInstance(id int32) (interface{}, bool) {
	G.Lock()
	defer G.Unlock()
	module, ok := G.modules[id]
	return module, ok
}

func cerr2Error(cstr *C.char) error {
	if cstr != nil {
		defer C.free(unsafe.Pointer(cstr))
		return errors.New(C.GoString(cstr))
	}
	return nil
}

//export ModuleConstruct
func ModuleConstruct(cls *C.char, nodeId int32, optionJson *C.char) int32 {
	if constructor, ok := GetModuleConstructor(C.GoString(cls)); ok {
		module, err := constructor.(ModuleConstructor)(nodeId, []byte(C.GoString(optionJson)))
		if err != nil {
			return -1 //construct failed
		}

		return RegisterModuleInstance(module)
	}

	return -2 //no constructor found
}

//export GetModuleInfoRegister
func GetModuleInfoRegister(cls *C.char, cinfo C.bmf_ModuleInfo) bool {
	if register, ok := getModuleInfoRegister(C.GoString(cls)); ok {
		info := WrapModuleInfo(cinfo)
		register.(ModuleInfoRegister)(info)
		return true
	}
	return false
}

// for testing
func ModuleConstructGo(cls string, nodeId int32, optionJson string) int32 {
	ccls := C.CString(cls)
	defer C.free(unsafe.Pointer(ccls))
	coption := C.CString(optionJson)
	defer C.free(unsafe.Pointer(coption))

	return ModuleConstruct(ccls, nodeId, coption)
}

//export ModuleProcess
func ModuleProcess(id int32, task C.bmf_Task) *C.char {
	if module, ok := GetModuleInstance(id); ok {
		task := WrapTask(task, false)
		if err := module.(Module).Process(task); err != nil {
			return C.CString(err.Error())
		}
	} else {
		panic(fmt.Sprintf("No module instance found by id %d", id))
	}

	return nil
}

func ModuleProcessGo(id int32, task *Task) error {
	errstr := ModuleProcess(id, task.p)
	if errstr != nil {
		defer C.free(unsafe.Pointer(errstr))
		return errors.New(C.GoString(errstr))
	}
	return nil
}

//export ModuleInit
func ModuleInit(id int32) *C.char {
	if module, ok := GetModuleInstance(id); ok {
		if err := module.(Module).Init(); err != nil {
			return C.CString(err.Error())
		}
	} else {
		panic(fmt.Sprintf("No module instance found by id %d", id))
	}
	return nil
}

func ModuleInitGo(id int32) error {
	return cerr2Error(ModuleInit(id))
}

//export ModuleReset
func ModuleReset(id int32) *C.char {
	if module, ok := GetModuleInstance(id); ok {
		if err := module.(Module).Reset(); err != nil {
			return C.CString(err.Error())
		}
	} else {
		panic(fmt.Sprintf("No module instance found by id %d", id))
	}
	return nil
}

func ModuleResetGo(id int32) error {
	return cerr2Error(ModuleReset(id))
}

//export ModuleClose
func ModuleClose(id int32) *C.char {
	if module, ok := GetModuleInstance(id); ok {
		UnregisterModuleInstance(id)
		if err := module.(Module).Close(); err != nil {
			return C.CString(err.Error())
		}
	}
	return nil
}

func ModuleCloseGo(id int32) error {
	return cerr2Error(ModuleClose(id))
}

//export ModuleGetInfo
func ModuleGetInfo(id int32) *C.char {
	if module, ok := GetModuleInstance(id); ok {
		info, err := module.(Module).GetModuleInfo()
		if err != nil {
			fmt.Printf("GoModule.GetModuleInfo failed for id %d", id)
			return nil
		}

		jsonStr, _ := json.Marshal(info)
		return C.CString(string(jsonStr))
	} else {
		panic(fmt.Sprintf("No module instance found by id %d", id))
	}
}

func ModuleGetInfoGo(id int32) interface{} {
	cstr := ModuleGetInfo(id)
	if cstr != nil {
		defer C.free(unsafe.Pointer(cstr))

		str := []byte(C.GoString(cstr))
		info := map[string]string{}
		json.Unmarshal(str, &info)
		return &info
	}

	return nil
}

//export ModuleNeedHungryCheck
func ModuleNeedHungryCheck(id int32, istreamId int32) bool {
	if module, ok := GetModuleInstance(id); ok {
		ret, err := module.(Module).NeedHungryCheck(istreamId)
		if err != nil {
			panic(err)
		}
		return ret
	} else {
		panic(fmt.Sprintf("No module instance found by id %d", id))
	}
}

//export ModuleIsHungry
func ModuleIsHungry(id int32, istreamId int32) bool {
	if module, ok := GetModuleInstance(id); ok {
		ret, err := module.(Module).IsHungry(istreamId)
		if err != nil {
			panic(err)
		}
		return ret
	} else {
		panic(fmt.Sprintf("No module instance found by id %d", id))
	}
}

//export ModuleIsInfinity
func ModuleIsInfinity(id int32) bool {
	if module, ok := GetModuleInstance(id); ok {
		ret, err := module.(Module).IsInfinity()
		if err != nil {
			panic(err)
		}
		return ret
	} else {
		panic(fmt.Sprintf("No module instance found by id %d", id))
	}
}

//export BmfSdkVersion
func BmfSdkVersion() *C.char {
	return C.CString(SdkVersion())
}

//export FreeCString
func FreeCString(p *C.char) {
	if (p != nil) {
		C.free(unsafe.Pointer(p))
	}
}
