package save

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"satisfactory-tool/util"
	"strconv"
	"strings"
)

func ProcessType(typeName string, data util.RawHolder, target *ReadOrWritable, buf *bytes.Buffer, inArray bool, depth int) (int, int) {
	switch typeName {
	case "NameProperty":
		fallthrough
	case "StrProperty":
		return ProcessStringProperty(data, target, buf)
	case "IntProperty":
		return ProcessIntProperty(data, target, buf, inArray)
	case "Int8Property":
		return ProcessInt8Property(data, target, buf, inArray)
	case "MapProperty":
		return ProcessMapProperty(data, target, buf, depth)
	case "ArrayProperty":
		return ProcessArrayProperty(data, target, buf, depth)
	case "ObjectProperty":
		return ProcessObjectProperty(data, target, buf, inArray)
	case "FloatProperty":
		return ProcessFloatProperty(data, target, buf)
	case "BoolProperty":
		return ProcessBoolProperty(data, target, buf)
	case "EnumProperty":
		return ProcessEnumProperty(data, target, buf)
	case "ByteProperty":
		return ProcessByteProperty(data, target, buf, inArray)
	case "StructProperty":
		return ProcessStructProperty(data, target, buf, nil, depth)
	case "TextProperty":
		return ProcessTextProperty(data, target, buf)
	}

	logrus.Panic("Don't know how to process: " + typeName)

	panic(1) // Logrus will panic for us
}

func ProcessProperty(data util.RawHolder, target *Property, buf *bytes.Buffer, depth int) int {
	padding := 0

	padding += util.RoWInt32StringNull(data.From(padding), &target.Name, buf)
	padding += 4

	if target.Name == "None" || target.Name == "" {
		return padding
	}

	padding += util.RoWInt32StringNull(data.From(padding), &target.Type, buf)
	padding += 4

	logrus.Debug(strings.Repeat(" ", depth), target.Type, " ", target.Name)

	util.RoWInt32(data.From(padding), &target.Size, buf)
	padding += 4

	util.RoWInt32(data.From(padding), &target.Index, buf)
	padding += 4

	valuePadding, actualValueSize := ProcessType(target.Type, data.FromNew(padding), &target.Value, buf, false, depth+1)
	padding += valuePadding

	if actualValueSize != int(target.Size) {
		logrus.Errorf("%s%s read %d bytes, expected %d (%d)\n", strings.Repeat(" ", depth), target.Type, actualValueSize, target.Size, int(target.Size)-actualValueSize)
	}

	return padding
}

func ProcessIntProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, inArray bool) (int, int) {
	// TODO Merge Reading/Writing
	if buf == nil {
		if inArray {
			var temp int32
			util.RoWInt32(data.From(0), &temp, buf)
			*target.(*ReadOrWritable) = temp
			return 4, 4
		}

		var tempInt int32
		util.RoWInt32(data.From(1), &tempInt, buf)
		*target.(*ReadOrWritable) = tempInt
		return 5, 4
	} else {
		var targetInt *int32

		switch v := (*target.(*ReadOrWritable)).(type) {
		case int:
			tempInt := int32((*target.(*ReadOrWritable)).(int))
			targetInt = &tempInt
		case float64:
			tempInt := int32((*target.(*ReadOrWritable)).(float64))
			targetInt = &tempInt
		case *int32:
			targetInt = (*target.(*ReadOrWritable)).(*int32)
		default:
			panic(fmt.Sprintf("Unknown integer type: %T\n", v))
		}

		if inArray {
			util.RoWInt32(data.From(0), targetInt, buf)
			return 4, 4
		}

		// Skip byte
		var skip byte = 0x00
		util.RoWB(data.At(0), &skip, buf)

		util.RoWInt32(data.From(1), targetInt, buf)
		return 5, 4
	}
}

func ProcessInt8Property(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, inArray bool) (int, int) {
	var targetInt *int8

	switch v := (*target.(*ReadOrWritable)).(type) {
	case int:
		tempInt := int8((*target.(*ReadOrWritable)).(int))
		targetInt = &tempInt
	case float64:
		tempInt := int8((*target.(*ReadOrWritable)).(float64))
		targetInt = &tempInt
	case *int32:
		tempInt := int8((*target.(*ReadOrWritable)).(int32))
		targetInt = &tempInt
	case *int8:
		targetInt = (*target.(*ReadOrWritable)).(*int8)
	case string:
		// No way to preserve reference to an incorrect type
		tempInt, err := strconv.Atoi((*target.(*ReadOrWritable)).(string))

		if err != nil {
			logrus.Panic(err)
		}

		tempInt32 := int8(tempInt)
		targetInt = &tempInt32
	default:
		panic(fmt.Sprintf("Unknown integer type: %T\n", v))
	}

	if inArray {
		util.RoWInt8(data.From(0), targetInt, buf)
		return 4, 4
	}

	// Skip byte
	var skip byte = 0x00
	util.RoWB(data.At(0), &skip, buf)

	util.RoWInt8(data.From(1), targetInt, buf)
	return 5, 4
}

func ProcessBoolProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	// TODO Merge reading/writing
	if buf == nil {
		*target.(*ReadOrWritable) = *data.At(0) > 0
		return 2, 0
	} else {
		var targetBool = (*target.(*ReadOrWritable)).(*bool)
		if *targetBool {
			buf.Write([]byte{0x01, 0x00})
		} else {
			buf.Write([]byte{0x00, 0x00})
		}
		return 2, 0
	}
}

func ProcessByteProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, inArray bool) (int, int) {
	var targetByte ByteProperty

	if buf == nil {
		*target.(*ReadOrWritable) = &targetByte
	} else {
		targetByte = *(*target.(*ReadOrWritable)).(*ByteProperty)
	}

	if !inArray {
		padding := 0

		padding += util.RoWInt32StringNull(data.From(padding), &targetByte.EnumType, buf)
		padding += 4

		// Skip byte
		var skip byte = 0x00
		util.RoWB(data.At(0), &skip, buf)
		padding += 1

		if targetByte.EnumType == "None" {
			util.RoWB(data.At(padding), &targetByte.Byte, buf)
			return padding + 1, 1
		} else {
			strLength := util.RoWInt32StringNull(data.From(padding), &targetByte.EnumName, buf)
			padding += 4 + strLength
			return padding, 4 + strLength
		}
	}

	util.RoWB(data.At(0), &targetByte.Byte, buf)
	return 1, 1
}

func ProcessEnumProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	var targetEnum EnumProperty

	if buf == nil {
		targetEnum = EnumProperty{}
		*target.(*ReadOrWritable) = &targetEnum
	} else {
		targetEnum = *(*target.(*ReadOrWritable)).(*EnumProperty)
	}

	padding := 0

	padding += util.RoWInt32StringNull(data.From(padding), &targetEnum.Type, buf)
	padding += 4

	// Skip byte
	var skip byte = 0x00
	util.RoWB(data.At(0), &skip, buf)
	padding += 1

	enumNameLength := util.RoWInt32StringNull(data.From(padding), &targetEnum.Name, buf)
	padding += 4 + enumNameLength

	return padding, enumNameLength + 4
}

func ProcessFloatProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	// TODO Merge reading/writing
	if buf == nil {
		var temp float32
		util.RoWFloat32(data.From(1), &temp, buf)
		*target.(*ReadOrWritable) = temp
		return 5, 4
	} else {
		// Skip byte
		var skip byte = 0x00
		util.RoWB(data.At(0), &skip, buf)

		var targetFloat = (*target.(*ReadOrWritable)).(*float32)
		util.RoWFloat32(data.From(1), targetFloat, buf)

		return 5, 4
	}
}

func ProcessStringProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	// TODO Merge reading/writing
	if buf == nil {
		padding := 1
		var temp string
		padding += util.RoWInt32StringNull(data.From(1), &temp, buf)
		padding += 4
		*target.(*ReadOrWritable) = temp
		return padding, padding - 1
	} else {
		padding := 0

		// Skip byte
		var skip byte = 0x00
		util.RoWB(data.At(0), &skip, buf)
		padding += 1

		var targetString = (*target.(*ReadOrWritable)).(*string)
		padding += util.RoWInt32StringNull(data.From(padding), targetString, buf)
		padding += 4

		return padding, padding - 1
	}
}

func ProcessTextProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	var targetText TextProperty

	if buf == nil {
		targetText = TextProperty{}
		*target.(*ReadOrWritable) = &targetText
	} else {
		targetText = *(*target.(*ReadOrWritable)).(*TextProperty)
	}

	padding := 0

	var skip byte = 0x00
	util.RoWB(data.At(0), &skip, buf)
	padding += 1

	util.RoWBytes(data.FromTo(padding, padding+13), &targetText.Magic, buf)
	padding += 13

	padding += util.RoWInt32StringNull(data.From(padding), &targetText.String, buf)
	padding += 4

	return padding, padding - 1
}

func ProcessObjectProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, inArray bool) (int, int) {
	var targetObject ObjectProperty

	if buf == nil {
		targetObject = ObjectProperty{}
		*target.(*ReadOrWritable) = &targetObject
	} else {
		targetObject = *(*target.(*ReadOrWritable)).(*ObjectProperty)
	}

	padding := 0
	overhead := 0

	if !inArray {
		// Skip byte
		var skip byte = 0x00
		util.RoWB(data.At(0), &skip, buf)
		overhead += 1
	}

	padding += util.RoWInt32StringNull(data.From(padding+overhead), &targetObject.World, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding+overhead), &targetObject.Class, buf)
	padding += 4

	return padding + overhead, padding
}

func ProcessArrayProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, depth int) (int, int) {
	var targetArray ArrayProperty

	if buf == nil {
		targetArray = ArrayProperty{}
		*target.(*ReadOrWritable) = &targetArray
	} else {
		targetArray = *(*target.(*ReadOrWritable)).(*ArrayProperty)
	}

	padding := 0

	padding += util.RoWInt32StringNull(data.From(padding), &targetArray.Type, buf)
	padding += 4

	// Skip byte
	var skip byte = 0x00
	util.RoWB(data.At(padding), &skip, buf)
	padding += 1

	var valueCount = int32(len(targetArray.Values))
	util.RoWInt32(data.From(padding), &valueCount, buf)
	padding += 4

	beforePadding := padding

	logrus.Debug(strings.Repeat(" ", depth), "[", targetArray.Type, "] x ", valueCount)

	if targetArray.Type == "StructProperty" {
		padding += util.RoWInt32StringNull(data.From(padding), &targetArray.StructName, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetArray.StructType, buf)
		padding += 4

		util.RoWInt32(data.From(padding), &targetArray.StructSize, buf)
		padding += 4

		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+4), &targetArray.Magic1, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetArray.StructClassType, buf)
		padding += 4

		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+17), &targetArray.Magic2, buf)
		padding += 17
	}

	if buf == nil {
		targetArray.Values = make([]ReadOrWritable, valueCount)
	}

	valueSizes := padding - beforePadding

	for i := 0; i < int(valueCount); i++ {
		if targetArray.Type == "StructProperty" {
			padded, _ := ProcessStructProperty(data.FromNew(padding), &targetArray.Values[i], buf, &targetArray.StructClassType, depth+1)
			padding += padded
			valueSizes += padded
		} else {
			padded, _ := ProcessType(targetArray.Type, data.FromNew(padding), &targetArray.Values[i], buf, true, depth+1)
			padding += padded
			valueSizes += padded
		}
	}

	return padding, valueSizes + 4
}

func ProcessStructProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, arrayTypeName *string, depth int) (int, int) {
	var targetStruct StructProperty

	if buf == nil {
		targetStruct = StructProperty{}
		*target.(*ReadOrWritable) = &targetStruct
	} else {
		targetStruct = *(*target.(*ReadOrWritable)).(*StructProperty)
	}

	padding := 0

	var typeName string

	if arrayTypeName == nil {
		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Type, buf)
		padding += 4
		typeName = targetStruct.Type

		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+17), &targetStruct.Magic, buf)
		padding += 17
	} else {
		typeName = *arrayTypeName
		targetStruct.Type = typeName
	}

	logrus.Debug(strings.Repeat(" ", depth), typeName)

	beforePadding := padding

	// TODO Merge reading/writing
	switch typeName {
	case "Vector":
		fallthrough
	case "Rotator":
		var tempCast util.Vector3

		if buf == nil {
			tempCast = util.Vector3{}
			targetStruct.Value = &tempCast
		} else {
			tempCast = targetStruct.Value.(util.Vector3)
		}

		util.RoWVec3(data.From(padding), &tempCast, buf)
		padding += 12

		return padding, padding - beforePadding
	case "Color":
		if buf == nil {
			targetStruct.Value = Color{
				R: *data.At(padding),
				G: *data.At(padding + 1),
				B: *data.At(padding + 2),
				A: *data.At(padding + 3),
			}
		} else {
			tempCast := targetStruct.Value.(Color)
			buf.WriteByte(tempCast.R)
			buf.WriteByte(tempCast.G)
			buf.WriteByte(tempCast.B)
			buf.WriteByte(tempCast.A)
		}

		return padding + 4, 4
	case "LinearColor":
		if buf == nil {
			var vec4 util.Vector4
			util.RoWVec4(data.From(padding), &vec4, buf)
			padding += 16

			targetStruct.Value = LinearColor{
				R: vec4.X,
				G: vec4.Y,
				B: vec4.Z,
				A: vec4.W,
			}
		} else {
			tempCast := targetStruct.Value.(LinearColor)

			vec4 := util.Vector4{
				X: tempCast.R,
				Y: tempCast.G,
				Z: tempCast.B,
				W: tempCast.A,
			}

			util.RoWVec4(data.From(padding), &vec4, buf)

			padding += 16
		}

		return padding, padding - beforePadding
	case "Quat":
		var tempCast util.Vector4

		if buf == nil {
			tempCast = util.Vector4{}
			targetStruct.Value = &tempCast
		} else {
			tempCast = targetStruct.Value.(util.Vector4)
		}

		util.RoWVec4(data.From(padding), &tempCast, buf)
		padding += 16

		return padding, padding - beforePadding
	case "Box":
		var tempCast Box

		if buf == nil {
			tempCast = Box{}
			targetStruct.Value = &tempCast
		} else {
			tempCast = targetStruct.Value.(Box)
		}

		util.RoWVec3(data.From(padding), &tempCast.Min, buf)
		padding += 12

		util.RoWVec3(data.From(padding), &tempCast.Max, buf)
		padding += 12

		util.RoWB(data.At(padding), &tempCast.Valid, buf)
		padding += 1

		return padding, padding - beforePadding
	case "InventoryItem":
		var tempCast InventoryItem

		if buf == nil {
			tempCast = InventoryItem{}
			targetStruct.Value = &tempCast
		} else {
			tempCast = targetStruct.Value.(InventoryItem)
		}

		padding += util.RoWInt32StringNull(data.From(padding), &tempCast.Magic, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &tempCast.ItemName, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &tempCast.LevelName, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &tempCast.PathName, buf)
		padding += 4

		return padding, padding - beforePadding
	case "RailroadTrackPosition":
		var tempCast RailroadTrackPosition

		if buf == nil {
			tempCast = RailroadTrackPosition{}
			targetStruct.Value = &tempCast
		} else {
			tempCast = targetStruct.Value.(RailroadTrackPosition)
		}

		padding += util.RoWInt32StringNull(data.From(padding), &tempCast.World, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &tempCast.EntityType, buf)
		padding += 4

		util.RoWFloat32(data.From(padding), &tempCast.Offset, buf)
		padding += 4

		util.RoWFloat32(data.From(padding), &tempCast.Forward, buf)
		padding += 4

		return padding, padding - beforePadding
	case "SplitterSortRule":
		fallthrough
	case "SchematicCost":
		fallthrough
	case "ResearchTime":
		fallthrough
	case "FeetOffset":
		fallthrough
	case "TimeTableStop":
		fallthrough
	case "RemovedInstance":
		fallthrough
	case "SpawnData":
		fallthrough
	case "MessageData":
		fallthrough
	case "ItemFoundData":
		fallthrough
	case "CompletedResearch":
		fallthrough
	case "ResearchCost":
		fallthrough
	case "PhaseCost":
		fallthrough
	case "ItemAmount":
		fallthrough
	case "SplinePointData":
		fallthrough
	case "InventoryStack":
		fallthrough
	case "RemovedInstanceArray":
		fallthrough
	case "Transform":
		if buf == nil {
			var values []Property
			valueSize := RoWToNone(data.FromNew(padding), &values, buf, depth+1)
			padding += valueSize

			targetStruct.Value = GenericStruct{
				Values: values,
			}

			return padding, valueSize
		} else {
			tempCast := targetStruct.Value.(GenericStruct)

			valueSize := RoWToNone(data.FromNew(padding), &tempCast.Values, buf, depth+1)
			padding += valueSize

			return padding, valueSize
		}
	}

	logrus.Panicf("Unknown struct: %s\n", typeName)

	panic(1) // Logrus will panic for us
}

func ProcessMapProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, depth int) (int, int) {
	var targetMap MapProperty

	if buf == nil {
		targetMap = MapProperty{}
		*target.(*ReadOrWritable) = &targetMap
	} else {
		targetMap = *(*target.(*ReadOrWritable)).(*MapProperty)
	}

	padding := 0

	padding += util.RoWInt32StringNull(data.From(padding), &targetMap.KeyType, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &targetMap.ValueType, buf)
	padding += 4

	// Skip byte
	var skip byte = 0x00
	util.RoWB(data.At(padding), &skip, buf)
	padding += 1

	beforePadding := padding

	// TODO Unknown
	util.RoWBytes(data.FromTo(padding, padding+4), &targetMap.Magic, buf)
	padding += 4

	var pairCount = int32(len(targetMap.Values))
	util.RoWInt32(data.From(padding), &pairCount, buf)
	padding += 4

	if buf == nil {
		targetMap.Values = make([]MapEntry, pairCount)
	}

	// TODO Merge reading/writing
	if buf == nil {
		for i := 0; i < int(pairCount); i++ {
			var key ReadOrWritable
			keySize, _ := ProcessType(targetMap.KeyType, data.FromNew(padding), &key, buf, true, depth+1)
			padding += keySize

			var values []Property
			padding += RoWToNone(data.FromNew(padding), &values, buf, depth+1)

			targetMap.Values[i] = MapEntry{
				Key:   key,
				Value: values,
			}
		}
	} else {
		for _, entry := range targetMap.Values {
			y := InterfaceIfy(entry.Key).(ReadOrWritable)

			padded, _ := ProcessType(targetMap.KeyType, data.FromNew(padding), &y, buf, true, depth+1)
			padding += padded

			padding += RoWToNone(data.FromNew(padding), &entry.Value, buf, depth+1)
		}
	}

	return padding, padding - beforePadding
}

// Absolutely necessary hack
func InterfaceIfy(data interface{}) interface{} {
	return data
}
