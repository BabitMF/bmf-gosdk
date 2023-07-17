package builder

import (
	"github.com/babitmf/bmf-gosdk/bmf"
	"errors"
	"fmt"
)

func (g *BMFGraph) newNode(inputs []*BMFStream, module *BMFModuleMeta, meta *BMFNodeMeta, option interface{}) *BMFNode {
	nid := len(g.nodes)
	g.nodes = append(g.nodes, nil)
	g.nodes[nid] = &BMFNode{
		graph:               g,
		id:                  nid,
		alias:               "",
		moduleInfo:          module,
		inputStreams:        inputs,
		outputStreams:       []*BMFStream{},
		metaInfo:            meta,
		option:              option,
		inputManager:        Immediate,
		scheduler:           0,
		existedStreamNotify: map[string]*BMFStream{},
	}
	if g.mode == Server {
		g.nodes[nid].inputManager = server
	}

	return g.nodes[nid]
}

func (g *BMFGraph) inPxy() *BMFNode {
	if g.inProxy == nil {
		g.inProxy = &BMFNode{
			graph: g,
			id:    -1,
			alias: "",
			moduleInfo: &BMFModuleMeta{
				name: "BMFGraphInProxy",
			},
			inputStreams:        []*BMFStream{},
			outputStreams:       []*BMFStream{},
			metaInfo:            &BMFNodeMeta{},
			option:              map[string]interface{}{},
			inputManager:        0,
			scheduler:           0,
			existedStreamNotify: map[string]*BMFStream{},
		}
	}
	return g.inProxy
}

func (g *BMFGraph) outPxy() *BMFNode {
	if g.outProxy == nil {
		g.outProxy = &BMFNode{
			graph: g,
			id:    -1,
			alias: "",
			moduleInfo: &BMFModuleMeta{
				name: "BMFGraphOutProxy",
			},
			inputStreams:        []*BMFStream{},
			outputStreams:       []*BMFStream{},
			metaInfo:            &BMFNodeMeta{},
			option:              map[string]interface{}{},
			inputManager:        0,
			scheduler:           0,
			existedStreamNotify: map[string]*BMFStream{},
		}
	}
	return g.outProxy
}

func (g *BMFGraph) isServer() bool {
	return g.mode == Server
}

func (g *BMFGraph) PyModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string, moduleEntry string,
	preModule *CBMFModule) *BMFNode {
	return g.Module(inputs, moduleName, Python, modulePath, moduleEntry, option, preModule)
}

func (g *BMFGraph) CppModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string, moduleEntry string,
	preModule *CBMFModule) *BMFNode {
	return g.Module(inputs, moduleName, Cpp, modulePath, moduleEntry, option, preModule)
}

func (g *BMFGraph) GoModule(inputs []*BMFStream, moduleName string, option interface{}, modulePath string, moduleEntry string,
	preModule *CBMFModule) *BMFNode {
	return g.Module(inputs, moduleName, Go, modulePath, moduleEntry, option, preModule)
}

func (g *BMFGraph) Module(inputs []*BMFStream, moduleName string, moduleType BMFModuleType, modulePath string, moduleEntry string,
	option interface{}, preModule *CBMFModule) *BMFNode {
	if inputs == nil {
		inputs = []*BMFStream{}
	}
	for i, s := range inputs {
		if s == nil {
			inputs[i] = &BMFStream{
				node:   nil,
				name:   fmt.Sprintf("InputStreamPlaceHolder_%d", i),
				notify: "",
				alias:  "",
			}
		}
	}
	preId := -1
	if preModule != nil {
		uid := preModule.UID()
		preId = int(uid)
	}

	return g.newNode(inputs, &BMFModuleMeta{
		name:     moduleName,
		language: moduleType,
		path:     modulePath,
		entry:    moduleEntry,
	}, &BMFNodeMeta{
		preModuleId:      preId,
		callbackBindings: make(map[int64]uint32),
	}, option)
}

func (g *BMFGraph) InStream(id interface{}) *BMFStream {
	s := g.inPxy().Stream(id)
	g.inputStreams = g.inPxy().outputStreams
	return s
}

func (g *BMFGraph) outStream(id interface{}) *BMFStream {
	s := g.outPxy().Stream(id)
	g.outputStreams = g.outPxy().outputStreams
	return s
}

func (g *BMFGraph) CheckAliasExistence(alias string) bool {
	_, ok := g.existedStreamAlias[alias]
	_, ok2 := g.existedNodeAlias[alias]
	return ok || ok2
}

func (g *BMFGraph) CheckStreamNotifyExistence(notify string) bool {
	for _, nd := range g.nodes {
		if _, ok := nd.existedStreamNotify[notify]; ok {
			return true
		}
	}
	return false
}

func (g *BMFGraph) GiveStreamAlias(alias string, stream *BMFStream) error {
	if g.CheckAliasExistence(alias) {
		if g.existedStreamAlias[alias] == stream {
			return errors.New(fmt.Sprintf("cannot give alias to the same stream more than once. (alias = %s)", alias))
		}
		return errors.New(fmt.Sprintf("stream alias duplicated with existing stream or node alias. (alias = %s)", alias))
	}
	if g.CheckStreamNotifyExistence(alias) {
		return errors.New(fmt.Sprintf("stream alias duplicated with existing stream notify. (alias = %s)", alias))
	}
	stream.alias = alias
	g.existedStreamAlias[alias] = stream
	return nil
}

func (g *BMFGraph) GiveNodeAlias(alias string, node *BMFNode) error {
	if g.CheckAliasExistence(alias) {
		if g.existedNodeAlias[alias] == node {
			return errors.New(fmt.Sprintf("cannot give alias to the same node more than once. (alias = %s)", alias))
		}
		return errors.New(fmt.Sprintf("node alias duplicated with existing stream or node alias. (alias = %s)", alias))
	}
	node.alias = alias
	g.existedNodeAlias[alias] = node
	return nil
}

func (g *BMFGraph) Instance() (*CBMFGraph, error) {
	if g.instance == nil {
		return nil, errors.New("graph must be run first")
	}
	return g.instance, nil
}

func (g *BMFGraph) Instantiate(needMerge bool) *CBMFGraph {
	if g.instance != nil {
		panic("Graph cannot be instantiated more than once")
	}
	g.instance, _ = NewCBMFGraph(g.ToInfo(), needMerge)
	return g.instance
}

func (g *BMFGraph) Run(needMerge bool) *CBMFGraph {
	if g.instance == nil {
		g.instance, _ = NewCBMFGraph(g.ToInfo(), needMerge)
	}
	g.instance.Start()
	return g.instance
}

func (g *BMFGraph) Start(needMerge bool, streamName string) func() (pkt *bmf.Packet) {
	g.outStream(streamName)
	if g.instance == nil {
		g.instance, _ = NewCBMFGraph(g.ToInfo(), needMerge)
	}
	g.instance.Start()
	return GeneratorClosure(g.instance, streamName)
}

func (g *BMFGraph) Close(){
	g.instance.Close()
}

func GeneratorClosure(graph *CBMFGraph, streamName string) func() (pkt *bmf.Packet) {
	g := graph
	return func() (pkt *bmf.Packet) {
		for true {
			pkt, _ = g.PollOutputStreamPacket(streamName)
			timestamp := pkt.Timestamp()
			if timestamp != -1 {
				break
			}
		}
		return
	}
}
