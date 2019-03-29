package util

import (
	"encoding/binary"
	"math"
)

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

type Vector4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

func Int16(b []byte) int16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return int16(b[0]) | int16(b[1])<<8
}

func Int32(b []byte) int32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}

func Int64(b []byte) int64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 |
		int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56
}

func Vec3(b []byte) Vector3 {
	return Vector3{
		X: Float32(b),
		Y: Float32(b[4:]),
		Z: Float32(b[8:]),
	}
}

func Vec4(b []byte) Vector4 {
	return Vector4{
		X: Float32(b),
		Y: Float32(b[4:]),
		Z: Float32(b[8:]),
		W: Float32(b[12:]),
	}
}

func Int32StringNull(b []byte) (string, int) {
	strLength := int(Int32(b))
	if strLength == 0 {
		return "", strLength
	}
	return string(b[4 : 4+strLength-1 /* Null termination */]), strLength
}

func Float32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b[:4]))
}
