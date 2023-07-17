package custom

/*
#include <stdint.h>
#include <stdlib.h>
#include <libavformat/avformat.h>

*/
import "C"

type AVFrame struct {
	p *C.AVFrame
}

func (self *AVFrame) Free() {
	p := self.p
	C.av_frame_free(&p)
}

func NewAVFrameWithVideo(width, height, format int) *AVFrame {
	p := C.av_frame_alloc()
	if p == nil {
		return nil
	}
	p.width = C.int(width)
	p.height = C.int(height)
	p.format = C.int(format)

	rc := C.av_frame_get_buffer(p, 32)
	if rc != 0 {
		C.av_frame_free(&p)
		return nil
	}

	return &AVFrame{p: p}
}
