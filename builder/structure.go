package builder

import "fmt"

func modeToString(m BMFGraphMode) string {
	switch m {
	case Normal:
		return "normal"
	case Server:
		return "Server"
	case Generator:
		return "Generator"
	case SubGraph:
		return "subgraph"
	case Update:
		return "Update"
	default:
		panic("Unknown graph mode.")
	}
}

func inputManagerToString(im BMFInputManager) string {
	switch im {
	case Immediate:
		return "immediate"
	case Default:
		return "default"
	case FrameSync:
		return "framesync"
	case ClockSync:
		return "clocksync"
	case server:
		return "server"
	default:
		panic("Unknown input stream manager.")
	}
}

type BMFUpdator struct {
	nodes    []*BMFUpdatorNode
	relatedG *BMFGraph
}

func (bu *BMFUpdator) ToInfo() *BMFUpdatorInfo {
	return &BMFUpdatorInfo{Nodes: batchedUpdatorNodeToInfo(bu.nodes)}
}

type BMFUpdatorInfo struct {
	Nodes []*BMFUpdatorNodeInfo `json:"nodes"`
}

type BMFUpdatorNode struct {
	id           int
	action       string
	alias        string
	moduleInfo   *BMFModuleMeta
	metaInfo     *BMFNodeMeta
	option       interface{}
	inputManager BMFInputManager
	scheduler    int
	inStream     []*BMFUpdatorStream
	outStream    []*BMFUpdatorStream
}

func batchedUpdatorNodeToInfo(nds []*BMFUpdatorNode) []*BMFUpdatorNodeInfo {
	var ret []*BMFUpdatorNodeInfo
	for _, nd := range nds {
		ret = append(ret, nd.toInfo())
	}
	return ret
}

func (bun *BMFUpdatorNode) toInfo() *BMFUpdatorNodeInfo {
	return &BMFUpdatorNodeInfo{
		Id:            bun.id,
		Action:        bun.action,
		Alias:         bun.alias,
		ModuleInfo:    bun.moduleInfo.toInfo(),
		MetaInfo:      bun.metaInfo.toInfo(),
		Option:        bun.option,
		InputManager:  inputManagerToString(bun.inputManager),
		Scheduler:     bun.scheduler,
		InputStreams:  batchedUpdatorStreamToInfo(bun.inStream),
		OutputStreams: batchedUpdatorStreamToInfo(bun.outStream),
	}
}

type BMFUpdatorNodeInfo struct {
	Id            int                     `json:"id"`
	Action        string                  `json:"action"`
	Alias         string                  `json:"alias"`
	ModuleInfo    *BMFModuleMetaInfo      `json:"module_info"`
	MetaInfo      *BMFNodeMetaInfo        `json:"meta_info"`
	Option        interface{}             `json:"option"`
	InputManager  string                  `json:"input_manager"`
	Scheduler     int                     `json:"scheduler"`
	InputStreams  []*BMFUpdatorStreamInfo `json:"input_streams"`
	OutputStreams []*BMFUpdatorStreamInfo `json:"output_streams"`
}

type BMFUpdatorStream struct {
	identifier string
	alias      string
}

func batchedUpdatorStreamToInfo(ss []*BMFUpdatorStream) []*BMFUpdatorStreamInfo {
	var ret []*BMFUpdatorStreamInfo
	for _, s := range ss {
		ret = append(ret, s.toInfo())
	}
	return ret
}

func (bus *BMFUpdatorStream) toInfo() *BMFUpdatorStreamInfo {
	return &BMFUpdatorStreamInfo{
		Identifier: bus.identifier,
		Alias:      bus.alias,
	}
}

type BMFUpdatorStreamInfo struct {
	Identifier string `json:"identifier"`
	Alias      string `json:"alias"`
}

type BMFGraph struct {
	// Config info
	mode          BMFGraphMode
	inputStreams  []*BMFStream
	outputStreams []*BMFStream
	nodes         []*BMFNode
	inProxy       *BMFNode
	outProxy      *BMFNode
	option        interface{}

	// External info for runtime checking
	existedStreamAlias map[string]*BMFStream
	existedNodeAlias   map[string]*BMFNode

	// Running Instance
	instance *CBMFGraph
}

type BMFGraphInfo struct {
	Mode          string           `json:"mode"`
	InputStreams  []*BMFStreamInfo `json:"input_streams"`
	OutputStreams []*BMFStreamInfo `json:"output_streams"`
	Nodes         []*BMFNodeInfo   `json:"nodes"`
	Option        interface{}      `json:"option"`
}

func (g *BMFGraph) ToInfo() *BMFGraphInfo {
	return &BMFGraphInfo{
		Mode:          modeToString(g.mode),
		InputStreams:  batchedStreamToInfo(g.inputStreams),
		OutputStreams: batchedStreamToInfo(g.outputStreams),
		Nodes:         batchedNodeToInfo(g.nodes),
		Option:        g.option,
	}
}

type BMFNode struct {
	// Pointer to graph
	graph *BMFGraph

	// Config info
	id            int
	alias         string
	moduleInfo    *BMFModuleMeta
	inputStreams  []*BMFStream
	outputStreams []*BMFStream
	metaInfo      *BMFNodeMeta
	option        interface{}
	inputManager  BMFInputManager
	scheduler     int

	// External info for runtime checking
	existedStreamNotify map[string]*BMFStream
}

type BMFNodeInfo struct {
	Id            int                `json:"id"`
	Alias         string             `json:"alias"`
	ModuleInfo    *BMFModuleMetaInfo `json:"module_info"`
	InputStreams  []*BMFStreamInfo   `json:"input_streams"`
	OutputStreams []*BMFStreamInfo   `json:"output_streams"`
	MetaInfo      *BMFNodeMetaInfo   `json:"meta_info"`
	Option        interface{}        `json:"option"`
	InputManager  string             `json:"input_manager"`
	Scheduler     int                `json:"scheduler"`
}

func batchedNodeToInfo(nds []*BMFNode) []*BMFNodeInfo {
	var ret []*BMFNodeInfo
	for _, nd := range nds {
		ret = append(ret, nd.ToInfo())
	}
	return ret
}

func (n *BMFNode) ToInfo() *BMFNodeInfo {
	return &BMFNodeInfo{
		Id:            n.id,
		ModuleInfo:    n.moduleInfo.toInfo(),
		InputStreams:  batchedStreamToInfo(n.inputStreams),
		OutputStreams: batchedStreamToInfo(n.outputStreams),
		MetaInfo:      n.metaInfo.toInfo(),
		Option:        n.option,
		InputManager:  inputManagerToString(n.inputManager),
		Scheduler:     n.scheduler,
	}
}

type BMFNodeMeta struct {
	preModuleId      int
	callbackBindings map[int64]uint32
}

type BMFNodeMetaInfo struct {
	PreModuleId      int              `json:"premodule_id"`
	CallbackBindings map[int64]uint32 `json:"callback_bindings"`
}

func (m *BMFNodeMeta) toInfo() *BMFNodeMetaInfo {
	return &BMFNodeMetaInfo{
		PreModuleId:      m.preModuleId,
		CallbackBindings: m.callbackBindings,
	}
}

type BMFModuleMeta struct {
	name     string
	language BMFModuleType
	path     string
	entry    string
}

type BMFModuleMetaInfo struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Path  string `json:"path"`
	Entry string `json:"entry"`
}

func batchedStreamToInfo(ss []*BMFStream) []*BMFStreamInfo {
	var ret []*BMFStreamInfo
	for _, s := range ss {
		ret = append(ret, s.ToInfo())
	}
	return ret
}

func (m *BMFModuleMeta) toInfo() *BMFModuleMetaInfo {
	return &BMFModuleMetaInfo{
		Name:  m.name,
		Type:  moduleTypeToString(m.language),
		Path:  m.path,
		Entry: m.entry,
	}
}

type BMFStream struct {
	// Pointer to Node
	node *BMFNode

	// Config info
	name   string
	notify string
	alias  string
}

type BMFStreamInfo struct {
	Identifier string `json:"identifier"`
	Alias      string `json:"alias"`
}

func (s *BMFStream) ToInfo() *BMFStreamInfo {
	if len(s.notify) == 0 {
		return &BMFStreamInfo{
			Identifier: s.name,
			Alias:      s.alias,
		}
	}
	return &BMFStreamInfo{
		Identifier: fmt.Sprintf("%s:%s", s.notify, s.name),
		Alias:      s.alias,
	}
}
