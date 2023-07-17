package builder

type BMFGraphMode int
type BMFInputManager int
type BMFModuleType int

type BMFCallbackType func([]byte) []byte

const (
	Normal BMFGraphMode = iota
	Server
	Generator
	SubGraph
	Update
)

const (
	Immediate BMFInputManager = iota
	Default
	server
	FrameSync
	ClockSync
)

const (
	Python BMFModuleType = iota
	Cpp
	Go
)

func moduleTypeToString(t BMFModuleType) string {
	switch t {
	case Python:
		return "python"
	case Cpp:
		return "c++"
	case Go:
		return "go"
	default:
		panic("Unknown module type.")
	}
}

func NewBMFGraph(mode BMFGraphMode, option interface{}) *BMFGraph {
	return &BMFGraph{
		mode:               mode,
		inputStreams:       []*BMFStream{},
		outputStreams:      []*BMFStream{},
		nodes:              []*BMFNode{},
		inProxy:            nil,
		outProxy:           nil,
		option:             option,
		existedStreamAlias: map[string]*BMFStream{},
		existedNodeAlias:   map[string]*BMFNode{},
	}
}

func NewBMFUpdator(baseGraph *BMFGraph) *BMFUpdator {
	return &BMFUpdator{
		nodes:    []*BMFUpdatorNode{},
		relatedG: baseGraph,
	}
}
