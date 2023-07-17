package custom

/*
#cgo CFLAGS: -I../cpp
#cgo LDFLAGS:  -L../cpp/build  -lavutil -lavformat -lswscale -lcustom_ffmpeg
#include <stdint.h>
#include <stdlib.h>
#include <custom_ffmpeg.h>
*/
import "C"
import (
	"errors"

	"github.com/babitmf/bmf-gosdk/bmf"
)

func BfmVfFromAVFrame(avf *C.AVFrame) (*bmf.VideoFrame, error) {
	p := C.bmf_vf_from_avframe(avf)
	if p == nil {
		return nil, errors.New(bmf.LastError())
	}

	return bmf.WrapVideoFrame(bmf.Ctype_bmf_VideoFrame(p), true), nil
}

func BfmVfToAVFrame(vf *bmf.VideoFrame) (*C.AVFrame, error) {
	p := C.bmf_vf_to_avframe(C.bmf_VideoFrame(vf.Pointer()))
	if p == nil {
		return nil, errors.New(bmf.LastError())
	}

	return p, nil
}

func BfmAfFromAVFrame(aaf *C.AVFrame) (*bmf.AudioFrame, error) {
	p := C.bmf_af_from_avframe(aaf)
	if p == nil {
		return nil, errors.New(bmf.LastError())
	}

	return bmf.WrapAudioFrame(bmf.Ctype_bmf_AudioFrame(p), true), nil
}

func BfmAfToAVFrame(af *bmf.AudioFrame) (*C.AVFrame, error) {
	p := C.bmf_af_to_avframe(C.bmf_AudioFrame(af.Pointer()))
	if p == nil {
		return nil, errors.New(bmf.LastError())
	}

	return p, nil
}
