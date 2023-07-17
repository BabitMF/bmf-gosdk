package hmp

/*
#include <stdint.h>
#include <stdlib.h>
#include "hmp_capi.h"
*/
import "C"
import "fmt"

type DeviceType int

const (
	kCPU  DeviceType = 0
	kCUDA DeviceType = 1

	CPU  DeviceType = kCPU
	CUDA DeviceType = kCUDA
)

func DeviceCount(d DeviceType) int {
	return int(C.hmp_device_count(C.int(d)))
}

func Device(d DeviceType, index int) string {
	if d == kCUDA {
		return fmt.Sprintf("cuda:%d", index)
	}
	return "cpu"
}
