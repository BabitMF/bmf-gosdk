package hmp

type ScalarType int

const (
	kUInt8   ScalarType = 0
	kInt8    ScalarType = 1
	kUInt16  ScalarType = 2
	kInt16   ScalarType = 3
	kInt32   ScalarType = 4
	kInt64   ScalarType = 5
	kFloat32 ScalarType = 6
	kFloat64 ScalarType = 7
	kHalf    ScalarType = 8

	UInt8   ScalarType = kUInt8
	Int8    ScalarType = kInt8
	UInt16  ScalarType = kUInt16
	Int16   ScalarType = kInt16
	Int32   ScalarType = kInt32
	Int64   ScalarType = kInt64
	Float32 ScalarType = kFloat32
	Float64 ScalarType = kFloat64
	Half    ScalarType = kHalf
)
