package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Vector4 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
	W float32 `json:"w"`
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

func WriteInt32StringNull(data string, buf *bytes.Buffer) int {
	nulled := 1
	if len(data) == 0 {
		nulled = 0
	}

	err := binary.Write(buf, binary.LittleEndian, int32(len(data)+nulled)) // Null termination

	if err != nil {
		panic(err)
	}

	if nulled == 1 {
		_, err = buf.WriteString(data)

		if err != nil {
			panic(err)
		}

		err = buf.WriteByte(0x00) // Null termination

		if err != nil {
			panic(err)
		}
	}

	return len(data) + nulled
}

func RoWInt32StringNull(data []byte, target *string, buf *bytes.Buffer) int {
	if target != nil {
		if buf == nil && data != nil {
			str, length := Int32StringNull(data)
			*target = str
			return length
		} else if buf != nil && data == nil {
			return WriteInt32StringNull(*target, buf)
		}
	}

	panic("Invalid State!")
}

func RoWInt32(data []byte, target *int32, buf *bytes.Buffer) {
	if target != nil {
		if buf == nil && data != nil {
			*target = Int32(data)
			return
		} else if buf != nil && data == nil {
			err := binary.Write(buf, binary.LittleEndian, *target)

			if err != nil {
				panic(err)
			}

			return
		}
	}

	panic("Invalid State!")
}

func RoWInt8(data []byte, target *int8, buf *bytes.Buffer) {
	if target != nil {
		if buf == nil && data != nil {
			*target = int8(data[0])
			return
		} else if buf != nil && data == nil {
			err := binary.Write(buf, binary.LittleEndian, *target)

			if err != nil {
				panic(err)
			}

			return
		}
	}

	panic("Invalid State!")
}

func RoWInt64(data []byte, target *int64, buf *bytes.Buffer) {
	if target != nil {
		if buf == nil && data != nil {
			*target = Int64(data)
			return
		} else if buf != nil && data == nil {
			err := binary.Write(buf, binary.LittleEndian, *target)

			if err != nil {
				panic(err)
			}

			return
		}
	}

	panic("Invalid State!")
}

func RoWFloat32(data []byte, target *float32, buf *bytes.Buffer) {
	if target != nil {
		if buf == nil && data != nil {
			*target = Float32(data)
			return
		} else if buf != nil && data == nil {
			err := binary.Write(buf, binary.LittleEndian, *target)

			if err != nil {
				panic(err)
			}

			return
		}
	}

	panic("Invalid State!")
}

func RoWB(data *byte, target *byte, buf *bytes.Buffer) {
	if target != nil {
		if buf == nil && data != nil {
			*target = *data
			return
		} else if buf != nil && data == nil {
			err := buf.WriteByte(*target)

			if err != nil {
				panic(err)
			}

			return
		}
	}

	panic("Invalid State!")
}

func RoWBytes(data []byte, target *[]byte, buf *bytes.Buffer) {
	if target != nil {
		if buf == nil && data != nil {
			*target = data
			return
		} else if buf != nil && data == nil {
			_, err := buf.Write(*target)

			if err != nil {
				panic(err)
			}

			return
		}
	}

	panic("Invalid State!")
}

func RoWVec4(data []byte, target *Vector4, buf *bytes.Buffer) {
	if target != nil {
		if buf == nil && data != nil {
			*target = Vec4(data)
			return
		} else if buf != nil && data == nil {
			err := binary.Write(buf, binary.LittleEndian, target.X)

			if err != nil {
				panic(err)
			}

			err = binary.Write(buf, binary.LittleEndian, target.Y)

			if err != nil {
				panic(err)
			}

			err = binary.Write(buf, binary.LittleEndian, target.Z)

			if err != nil {
				panic(err)
			}

			err = binary.Write(buf, binary.LittleEndian, target.W)

			if err != nil {
				panic(err)
			}

			return
		}
	}

	panic("Invalid State!")
}

func RoWVec3(data []byte, target *Vector3, buf *bytes.Buffer) {
	if target != nil {
		if buf == nil && data != nil {
			*target = Vec3(data)
			return
		} else if buf != nil && data == nil {
			err := binary.Write(buf, binary.LittleEndian, target.X)

			if err != nil {
				panic(err)
			}

			err = binary.Write(buf, binary.LittleEndian, target.Y)

			if err != nil {
				panic(err)
			}

			err = binary.Write(buf, binary.LittleEndian, target.Z)

			if err != nil {
				panic(err)
			}

			return
		}
	}

	panic("Invalid State!")
}

func Float32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b[:4]))
}

func HexDump(data []byte) string {
	result := ""

	perRow := 32
	rows := int(math.Ceil(float64(len(data)) / float64(perRow)))

	rowWidth := perRow * 5
	if len(data) < perRow {
		rowWidth = len(data) * 5
	}

	for i := 0; i < rows; i++ {
		hexSide := ""
		charSide := ""
		for k := 0; k < perRow && k < len(data[i*perRow:]); k++ {
			hexSide += fmt.Sprintf("%#-4x", data[i*perRow+k]) + " "
			charSide += fmt.Sprintf("%s", safeChar(data[i*perRow+k]))
		}
		result += fmt.Sprintf("%-#6x: %-"+strconv.Itoa(rowWidth)+"s%s\n", i*perRow, hexSide, charSide)
	}

	return result
}

func safeChar(char byte) string {
	if char <= 0x1F {
		return "."
	}

	return string(char)
}

type RawHolder struct {
	Data []byte
}

func (raw *RawHolder) From(from int) []byte {
	if raw.Data == nil {
		return nil
	}

	return raw.Data[from:]
}

func (raw *RawHolder) FromTo(from int, to int) []byte {
	if raw.Data == nil {
		return nil
	}

	return raw.Data[from:to]
}

func (raw *RawHolder) At(at int) *byte {
	if raw.Data == nil {
		return nil
	}

	return &raw.Data[at]
}

func (raw *RawHolder) FromNew(from int) RawHolder {
	if raw.Data == nil {
		return RawHolder{}
	}

	return RawHolder{
		Data: raw.Data[from:],
	}
}

func (raw *RawHolder) IsNil() bool {
	return raw.Data == nil
}

func (raw *RawHolder) Length() int {
	if raw.Data == nil {
		return 0
	}

	return len(raw.Data)
}

func (raw *RawHolder) FromToNew(from int, to int) RawHolder {
	if raw.Data == nil {
		return RawHolder{}
	}

	return RawHolder{
		Data: raw.Data[from:to],
	}
}
