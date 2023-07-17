package builder

import (
	"github.com/babitmf/bmf-gosdk/bmf"
	"errors"
	"fmt"
)

func (s *BMFStream) setAlias(alias string) error {
	if len(s.alias) != 0 {
		return errors.New(fmt.Sprintf("One stream must be given alias only once. (had {%s}, giving {%s})", s.alias, alias))
	}

	var idx int
	for idx = 0; idx < len(s.node.outputStreams); idx++ {
		if s == s.node.outputStreams[idx] {
			break
		}
	}

	return s.node.GiveStreamAlias(idx, alias)
}

func (s *BMFStream) setNotify(notify string) error {
	if len(s.notify) != 0 {
		return errors.New(fmt.Sprintf("One stream must be given notify only once. (had {%s}, giving {%s})", s.notify, notify))
	}
	var idx int
	for idx = 0; idx < len(s.node.outputStreams); idx++ {
		if s == s.node.outputStreams[idx] {
			break
		}
	}

	return s.node.GiveStreamNotify(idx, notify)
}

func (s *BMFStream) Module(inputs []*BMFStream, moduleName string, moduleType BMFModuleType, modulePath string, moduleEntry string,
	option interface{}, preModule *CBMFModule) *BMFNode {
	if inputs == nil {
		inputs = []*BMFStream{s}
	} else {
		inputs = append([]*BMFStream{s}, inputs...)
	}
	return s.node.graph.Module(inputs, moduleName, moduleType, modulePath, moduleEntry, option, preModule)
}

func (s *BMFStream) PyModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string,
	moduleEntry string, preModule *CBMFModule) *BMFNode {
	if inputs == nil {
		inputs = []*BMFStream{s}
	} else {
		inputs = append([]*BMFStream{s}, inputs...)
	}
	return s.node.graph.PyModule(inputs, moduleName, option, modulePath, moduleEntry, preModule)
}

func (s *BMFStream) CppModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string,
	moduleEntry string, preModule *CBMFModule) *BMFNode {
	if inputs == nil {
		inputs = []*BMFStream{s}
	} else {
		inputs = append([]*BMFStream{s}, inputs...)
	}
	return s.node.graph.CppModule(inputs, moduleName, option, modulePath, moduleEntry, preModule)
}

func (s *BMFStream) GoModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string,
	moduleEntry string, preModule *CBMFModule) *BMFNode {
	if inputs == nil {
		inputs = []*BMFStream{s}
	} else {
		inputs = append([]*BMFStream{s}, inputs...)
	}
	return s.node.graph.GoModule(inputs, moduleName, option, modulePath, moduleEntry, preModule)
}

func (s *BMFStream) Start() func() (pkt *bmf.Packet) {
	identifier := s.ToInfo().Identifier
	return s.node.graph.Start(true, identifier)
}
