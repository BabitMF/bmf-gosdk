package builder

import (
	"errors"
	"fmt"
)

func (n *BMFNode) Stream(id interface{}) *BMFStream {
	switch id.(type) {
		case int, int8, int16, int32, int64:
			sid := id.(int)
			if sid < 0 {
				panic(fmt.Sprintf("stream id cannot be negative. (requesting id = %d)", sid))
			}
			for len(n.outputStreams) <= sid {
				n.outputStreams = append(n.outputStreams, &BMFStream{
					node:   n,
					name:   fmt.Sprintf("%s_%d_%d", n.moduleInfo.name, n.id, len(n.outputStreams)),
					notify: "",
					alias:  "",
				})
			}
			return n.outputStreams[sid]
		case string:
			sid := id.(string)
			n.outputStreams = append(n.outputStreams, &BMFStream{
				node:   n,
				name:   sid,
				notify: "",
				alias:  "",
			})
			return n.outputStreams[0]
	}
	panic("stream indexing only supports integer and string")
}

func (n *BMFNode) GiveStreamNotify(id int, notify string) error {
	if id < 0 {
		return errors.New(fmt.Sprintf("stream id cannot be negative. (requesting id = %d)", id))
	}
	if id >= len(n.outputStreams) {
		return errors.New(fmt.Sprintf("requesting unexisted stream. (id = %d)", id))
	}
	if s, ok := n.existedStreamNotify[notify]; ok {
		if s == n.outputStreams[id] {
			return errors.New(fmt.Sprintf("cannot give notify to the same stream more than once. (notify = %s)", notify))
		}
		return errors.New(fmt.Sprintf("stream notify duplicated. (notify = %s)", notify))
	}
	if n.graph.CheckAliasExistence(notify) {
		return errors.New(fmt.Sprintf("stream notify duplicated with existing stream or node alias. (nofify = %s)", notify))
	}
	n.outputStreams[id].notify = notify
	n.existedStreamNotify[notify] = n.outputStreams[id]
	return nil
}

func (n *BMFNode) GiveStreamAlias(id int, alias string) error {
	if id < 0 {
		return errors.New(fmt.Sprintf("stream id cannot be negative. (requesting id = %d)", id))
	}
	if id >= len(n.outputStreams) {
		return errors.New(fmt.Sprintf("requesting unexisted stream. (id = %d)", id))
	}
	return n.graph.GiveStreamAlias(alias, n.outputStreams[id])
}

func (n *BMFNode) SetInputManager(im BMFInputManager) error {
	if n.graph.isServer() {
		return errors.New("cannot modify input manager under server mode")
	}
	if im == server {
		return errors.New("cannot set {server} input manager in a non-server-mode graph")
	}
	n.inputManager = im
	return nil
}

func (n *BMFNode) SetScheduler(schedulerId int) {
	n.scheduler = schedulerId
}

func (n *BMFNode) SetPreModule(premodule *CBMFModule) error {
	if n.metaInfo.preModuleId != -1 {
		return errors.New("cannot assign premodule to node more than once")
	}
	uid := premodule.UID()
	n.metaInfo.preModuleId = int(uid)
	return nil
}

func (n *BMFNode) Module(inputs []*BMFStream, moduleName string, moduleType BMFModuleType, modulePath string, moduleEntry string,
	option interface{}, preModule *CBMFModule) *BMFNode {
	return n.Stream(0).Module(inputs, moduleName, moduleType, modulePath, moduleEntry, option, preModule)
}

func (n *BMFNode) PyModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string, moduleEntry string,
	preModule *CBMFModule) *BMFNode {
	return n.Stream(0).PyModule(inputs, moduleName, option, modulePath, moduleEntry, preModule)
}

func (n *BMFNode) CppModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string, moduleEntry string,
	preModule *CBMFModule) *BMFNode {
	return n.Stream(0).CppModule(inputs, moduleName, option, modulePath, moduleEntry, preModule)
}

func (n *BMFNode) GoModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string, moduleEntry string,
	preModule *CBMFModule) *BMFNode {
	return n.Stream(0).GoModule(inputs, moduleName, option, modulePath, moduleEntry, preModule)
}
