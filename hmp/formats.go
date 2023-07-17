package hmp

/*
#include <stdint.h>
#include <stdlib.h>
#include "hmp_capi.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type ColorPrimaries int
type ColorTransferCharacteristic int
type ColorSpace int
type ColorRange int
type PixelFormat int
type ChannelFormat int

const (
	CP_RESERVED0   ColorPrimaries = 0
	CP_BT709       ColorPrimaries = 1
	CP_UNSPECIFIED ColorPrimaries = 2
	CP_RESERVED    ColorPrimaries = 3
	CP_BT470M      ColorPrimaries = 4

	CP_BT470BG      ColorPrimaries = 5
	CP_SMPTE170M    ColorPrimaries = 6
	CP_SMPTE240M    ColorPrimaries = 7
	CP_FILM         ColorPrimaries = 8
	CP_BT2020       ColorPrimaries = 9
	CP_SMPTE428     ColorPrimaries = 10
	CP_SMPTEST428_1 ColorPrimaries = CP_SMPTE428
	CP_SMPTE431     ColorPrimaries = 11
	CP_SMPTE432     ColorPrimaries = 12
	CP_EBU3213      ColorPrimaries = 22
	CP_JEDEC_P22    ColorPrimaries = CP_EBU3213
)

const (
	CTC_RESERVED0    ColorTransferCharacteristic = 0
	CTC_BT709        ColorTransferCharacteristic = 1
	CTC_UNSPECIFIED  ColorTransferCharacteristic = 2
	CTC_RESERVED     ColorTransferCharacteristic = 3
	CTC_GAMMA22      ColorTransferCharacteristic = 4
	CTC_GAMMA28      ColorTransferCharacteristic = 5
	CTC_SMPTE170M    ColorTransferCharacteristic = 6
	CTC_SMPTE240M    ColorTransferCharacteristic = 7
	CTC_LINEAR       ColorTransferCharacteristic = 8
	CTC_LOG          ColorTransferCharacteristic = 9
	CTC_LOG_SQRT     ColorTransferCharacteristic = 10
	CTC_IEC61966_2_4 ColorTransferCharacteristic = 11
	CTC_BT1361_ECG   ColorTransferCharacteristic = 12
	CTC_IEC61966_2_1 ColorTransferCharacteristic = 13
	CTC_BT2020_10    ColorTransferCharacteristic = 14
	CTC_BT2020_12    ColorTransferCharacteristic = 15
	CTC_SMPTE2084    ColorTransferCharacteristic = 16
	CTC_SMPTEST2084  ColorTransferCharacteristic = CTC_SMPTE2084
	CTC_SMPTE428     ColorTransferCharacteristic = 17
	CTC_SMPTEST428_1 ColorTransferCharacteristic = CTC_SMPTE428
	CTC_ARIB_STD_B67 ColorTransferCharacteristic = 18
)

const (
	CS_RGB                ColorSpace = 0
	CS_BT709              ColorSpace = 1
	CS_UNSPECIFIED        ColorSpace = 2
	CS_RESERVED           ColorSpace = 3
	CS_FCC                ColorSpace = 4
	CS_BT470BG            ColorSpace = 5
	CS_SMPTE170M          ColorSpace = 6
	CS_SMPTE240M          ColorSpace = 7
	CS_YCGCO              ColorSpace = 8
	CS_YCOCG              ColorSpace = CS_YCGCO
	CS_BT2020_NCL         ColorSpace = 9
	CS_BT2020_CL          ColorSpace = 10
	CS_SMPTE2085          ColorSpace = 11
	CS_CHROMA_DERIVED_NCL ColorSpace = 12
	CS_CHROMA_DERIVED_CL  ColorSpace = 13
	CS_ICTCP              ColorSpace = 14
)

const (
	CR_UNSPECIFIED ColorRange = 0
	CR_MPEG        ColorRange = 1
	CR_JPEG        ColorRange = 2
)

const (
	PF_NONE    PixelFormat = -1
	PF_YUV420P PixelFormat = 0
	PF_YUV422P PixelFormat = 4
	PF_YUV444P PixelFormat = 5
	PF_NV12    PixelFormat = 23
	PF_NV21    PixelFormat = 24

	PF_GRAY8  PixelFormat = 8
	PF_RGB24  PixelFormat = 2
	PF_RGBA32 PixelFormat = 26

	PF_GRAY16 PixelFormat = 30
	PF_RGB48  PixelFormat = 35
	PF_RGBA64 PixelFormat = 107
)

const (
	kNCHW ChannelFormat = 0
	kNHWC ChannelFormat = 1

	NCHW ChannelFormat = kNCHW
	NHWC ChannelFormat = kNHWC
)

type Ctype_hmp_ColorModel = C.hmp_ColorModel
type Ctype_hmp_PixelInfo = C.hmp_PixelInfo

type ColorModel struct {
	p   C.hmp_ColorModel
	own bool
}

func deleteColorModel(o *ColorModel) {
	if o != nil && o.own && o.p != nil {
		C.hmp_color_model_free(o.p)
		o.p = nil
	}
}

func WrapColorModel(p C.hmp_ColorModel, own bool) *ColorModel {
	o := &ColorModel{p: p, own: own}
	runtime.SetFinalizer(o, deleteColorModel)
	return o
}

func NewColorModel(cs ColorSpace, cr ColorRange, cp ColorPrimaries,
	ctc ColorTransferCharacteristic) *ColorModel {
	cm := C.hmp_color_model(C.int(cs), C.int(cr), C.int(cp), C.int(ctc))
	return WrapColorModel(cm, true)
}

func (self *ColorModel) Free() {
	runtime.SetFinalizer(self, nil)
	deleteColorModel(self)
}

func (self *ColorModel) Pointer() C.hmp_ColorModel {
	return self.p
}

func (self *ColorModel) Space() ColorSpace {
	ret := C.hmp_color_model_space(self.p)
	return ColorSpace(ret)
}

func (self *ColorModel) Range() ColorRange {
	ret := C.hmp_color_model_range(self.p)
	return ColorRange(ret)
}

func (self *ColorModel) Primaries() ColorPrimaries {
	ret := C.hmp_color_model_primaries(self.p)
	return ColorPrimaries(ret)
}

func (self *ColorModel) TransferCharacteristic() ColorTransferCharacteristic {
	ret := C.hmp_color_model_ctc(self.p)
	return ColorTransferCharacteristic(ret)
}

type PixelInfo struct {
	p   C.hmp_PixelInfo
	own bool
}

func deletePixelInfo(o *PixelInfo) {
	if o != nil && o.own && o.p != nil {
		C.hmp_pixel_info_free(o.p)
		o.p = nil
	}
}

func WrapPixelInfo(p C.hmp_PixelInfo, own bool) *PixelInfo {
	o := &PixelInfo{p: p, own: own}
	runtime.SetFinalizer(o, deletePixelInfo)
	return o
}

func NewPixelInfo(format PixelFormat, cm *ColorModel) *PixelInfo {
	p := C.hmp_pixel_info(C.int(format), cm.p)
	return WrapPixelInfo(p, true)
}

func NewPixelInfoV1(format PixelFormat, cs ColorSpace, cr ColorRange) *PixelInfo {
	p := C.hmp_pixel_info_v1(C.int(format), C.int(cs), C.int(cr))
	return WrapPixelInfo(p, true)
}

func NewPixelInfoV2(format PixelFormat, cp ColorPrimaries, ctc ColorTransferCharacteristic) *PixelInfo {
	p := C.hmp_pixel_info_v2(C.int(format), C.int(cp), C.int(ctc))
	return WrapPixelInfo(p, true)
}

func (self *PixelInfo) Free() {
	runtime.SetFinalizer(self, nil)
	deletePixelInfo(self)
}

func (self *PixelInfo) Pointer() C.hmp_PixelInfo {
	return self.p
}

func (self *PixelInfo) Format() PixelFormat {
	ret := C.hmp_pixel_info_format(self.p)
	return PixelFormat(ret)
}

func (self *PixelInfo) Space() ColorSpace {
	ret := C.hmp_pixel_info_space(self.p)
	return ColorSpace(ret)
}

func (self *PixelInfo) Range() ColorRange {
	ret := C.hmp_pixel_info_range(self.p)
	return ColorRange(ret)
}

func (self *PixelInfo) Primaries() ColorPrimaries {
	ret := C.hmp_pixel_info_primaries(self.p)
	return ColorPrimaries(ret)
}

func (self *PixelInfo) TransferCharacteristic() ColorTransferCharacteristic {
	ret := C.hmp_pixel_info_ctc(self.p)
	return ColorTransferCharacteristic(ret)
}

func (self *PixelInfo) InferSpace() ColorSpace {
	ret := C.hmp_pixel_info_infer_space(self.p)
	return ColorSpace(ret)
}

func (self *PixelInfo) ColorModel() *ColorModel {
	ret := C.hmp_pixel_info_color_model(self.p)
	return WrapColorModel(ret, false)
}

func (self *PixelInfo) IsRgbx() bool {
	ret := C.hmp_pixel_info_is_rgbx(self.p)
	return bool(ret)
}

func (self *PixelInfo) String() string {
	var c_str_size = int(0)
	c_str := C.hmp_pixel_info_stringfy(self.p, (*C.int)(unsafe.Pointer(&c_str_size)))
	return string((*[1 << 31]byte)(unsafe.Pointer(c_str))[:c_str_size:c_str_size])
}
