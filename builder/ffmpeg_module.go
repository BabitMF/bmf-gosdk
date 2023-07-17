package builder

func (g *BMFGraph) Decode(decodePara interface{}, controlStream *BMFStream) *BMFNode {
	var input []*BMFStream
	if controlStream != nil {
		input = append(input, controlStream)
	}
	return g.Module(input, "c_ffmpeg_decoder", Cpp, "", "", decodePara, nil)
}

func (n *BMFNode) Decode(decodePara interface{}) *BMFNode {
	return n.graph.Decode(decodePara, n.Stream(0))
}

func (s *BMFStream) Decode(decodePara interface{}) *BMFNode {
	return s.node.graph.Decode(decodePara, s)
}

func (g *BMFGraph) Encode(videoStream *BMFStream, audioStream *BMFStream, encoderPara interface{}) *BMFNode {
	var input []*BMFStream
	if videoStream == nil {
		input = append(input, &BMFStream{
			node:   nil,
			name:   "EncoderPlaceHolder_Video",
			notify: "",
			alias:  "",
		})
	} else {
		input = append(input, videoStream)
	}
	if audioStream == nil {
		input = append(input, &BMFStream{
			node:   nil,
			name:   "EncoderPlaceHolder_Audio",
			notify: "",
			alias:  "",
		})
	} else {
		input = append(input, audioStream)
	}
	return g.Module(input, "c_ffmpeg_encoder", Cpp, "", "", encoderPara, nil)
}

func (n *BMFNode) Encode(audioStream *BMFStream, encodePara interface{}) *BMFNode {
	return n.graph.Encode(n.Stream(0), audioStream, encodePara)
}

func (s *BMFStream) Encode(audioStream *BMFStream, encodePara interface{}) *BMFNode {
	return s.node.graph.Encode(s, audioStream, encodePara)
}

func (g *BMFGraph) FFmpegFilter(inputStreams []*BMFStream, filterName string, filterPara interface{}) *BMFNode {
	opt := make(map[string]string)
	opt["name"] = filterName
	if filterPara != nil {
		opt["para"] = dumpFilterOption(filterPara)
	}
	return g.Module(inputStreams, "c_ffmpeg_filter", Cpp, "", "", opt, nil)
}

func (n *BMFNode) FFmpegFilter(inputStreams []*BMFStream, filterName string, filterPara interface{}) *BMFNode {
	if inputStreams == nil {
		inputStreams = []*BMFStream{n.Stream(0)}
	} else {
		inputStreams = append([]*BMFStream{n.Stream(0)}, inputStreams...)
	}
	return n.graph.FFmpegFilter(inputStreams, filterName, filterPara)
}

func (s *BMFStream) FFmpegFilter(inputStream []*BMFStream, filterName string, filterPara interface{}) *BMFNode {
	if inputStream == nil {
		inputStream = []*BMFStream{s}
	} else {
		inputStream = append([]*BMFStream{s}, inputStream...)
	}
	return s.node.graph.FFmpegFilter(inputStream, filterName, filterPara)
}

func (g *BMFGraph) Vflip(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "vflip", filterPara)
}

func (n *BMFNode) Vflip(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "vflip", filterPara)
}

func (s *BMFStream) Vflip(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "vflip", filterPara)
}

func (g *BMFGraph) Scale(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "scale", filterPara)
}

func (n *BMFNode) Scale(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "scale", filterPara)
}

func (s *BMFStream) Scale(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "scale", filterPara)
}

func (g *BMFGraph) Setsar(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "setsar", filterPara)
}

func (n *BMFNode) Setsar(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "setsar", filterPara)
}

func (s *BMFStream) Setsar(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "setsar", filterPara)
}

func (g *BMFGraph) Pad(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "pad", filterPara)
}

func (n *BMFNode) Pad(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "pad", filterPara)
}

func (s *BMFStream) Pad(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "pad", filterPara)
}

func (g *BMFGraph) Trim(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "trim", filterPara)
}

func (n *BMFNode) Trim(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "trim", filterPara)
}

func (s *BMFStream) Trim(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "trim", filterPara)
}

func (g *BMFGraph) Setpts(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "setpts", filterPara)
}

func (n *BMFNode) Setpts(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "setpts", filterPara)
}

func (s *BMFStream) Setpts(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "setpts", filterPara)
}

func (g *BMFGraph) Loop(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "loop", filterPara)
}

func (n *BMFNode) Loop(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "loop", filterPara)
}

func (s *BMFStream) Loop(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "loop", filterPara)
}

func (g *BMFGraph) Split(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "split", filterPara)
}

func (n *BMFNode) Split(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "split", filterPara)
}

func (s *BMFStream) Split(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "split", filterPara)
}

func (g *BMFGraph) Adelay(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "adelay", filterPara)
}

func (n *BMFNode) Adelay(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "adelay", filterPara)
}

func (s *BMFStream) Adelay(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "adelay", filterPara)
}

func (g *BMFGraph) Atrim(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "atrim", filterPara)
}

func (n *BMFNode) Atrim(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "atrim", filterPara)
}

func (s *BMFStream) Atrim(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "atrim", filterPara)
}

func (g *BMFGraph) Afade(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "afade", filterPara)
}

func (n *BMFNode) Afade(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "afade", filterPara)
}

func (s *BMFStream) Afade(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "afade", filterPara)
}

func (g *BMFGraph) Asetpts(inputStream *BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "asetpts", filterPara)
}

func (n *BMFNode) Asetpts(filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(nil, "asetpts", filterPara)
}

func (s *BMFStream) Asetpts(filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(nil, "asetpts", filterPara)
}

func (g *BMFGraph) Amix(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter(inputStream, "amix", filterPara)
}

func (n *BMFNode) Amix(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(inputStream, "amix", filterPara)
}

func (s *BMFStream) Amix(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(inputStream, "amix", filterPara)
}

func (g *BMFGraph) Overlay(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter(inputStream, "overlay", filterPara)
}

func (n *BMFNode) Overlay(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(inputStream, "overlay", filterPara)
}

func (s *BMFStream) Overlay(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(inputStream, "overlay", filterPara)
}

func (g *BMFGraph) Concat(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return g.FFmpegFilter(inputStream, "concat", filterPara)
}

func (n *BMFNode) Concat(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return n.FFmpegFilter(inputStream, "concat", filterPara)
}

func (s *BMFStream) Concat(inputStream []*BMFStream, filterPara interface{}) *BMFNode {
	return s.FFmpegFilter(inputStream, "concat", filterPara)
}

func (g *BMFGraph) Fps(inputStream *BMFStream, fps int) *BMFNode {
	return g.FFmpegFilter([]*BMFStream{inputStream}, "fps", map[string]int{"fps": fps})
}

func (n *BMFNode) Fps(fps int) *BMFNode {
	return n.FFmpegFilter(nil, "fps", map[string]int{"fps": fps})
}

func (s *BMFStream) Fps(fps int) *BMFNode {
	return s.FFmpegFilter(nil, "fps", map[string]int{"fps": fps})
}
