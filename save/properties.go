package save

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"satisfactory-tool/util"
	"strconv"
	"strings"
)

func ResolveType(typeName string, data []byte, inArray bool, depth int) (ReadOrWritable, int, int) {
	switch typeName {
	case "StrProperty":
		return ParseStringProperty(data)
	case "TextProperty":
		return ParseTextProperty(data)
	case "NameProperty":
		return ParseStringProperty(data)
	case "BoolProperty":
		return ParseBoolProperty(data)
	case "ByteProperty":
		return ParseByteProperty(data)
	case "EnumProperty":
		return ParseEnumProperty(data)
	case "FloatProperty":
		return ParseFloatProperty(data)
	case "IntProperty":
		return ParseIntProperty(data, inArray)
	case "Int8Property":
		return ParseInt8Property(data, inArray)
	case "ObjectProperty":
		return ParseObjectProperty(data, inArray)
	case "ArrayProperty":
		return ParseArrayProperty(data, depth)
	case "MapProperty":
		return ParseMapProperty(data, depth)
	case "StructProperty":
		return ParseStructProperty(data, nil, depth)
	}

	logrus.Panic("Don't know how to process: " + typeName)

	panic(1) // Logrus will panic for us
}

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
		return ProcessByteProperty(data, target, buf)
	case "StructProperty":
		return ProcessStructProperty(data, target, buf, nil, depth)
	case "TextProperty":
		return ProcessTextProperty(data, target, buf)
	}

	logrus.Panic("Don't know how to process: " + typeName)

	panic(1) // Logrus will panic for us
}

func ParseProperty(data []byte, depth int) (Property, int) {
	padding := 0

	name, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	if name == "None" || name == "" {
		return Property{
			Name:  "None",
			Type:  "",
			Index: 0,
			Size:  0,
			Value: "",
		}, padding
	}

	typeName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	logrus.Debug(strings.Repeat(" ", depth), typeName, " ", name)

	valueSize := util.Int32(data[padding:])
	padding += 4

	keyIndex := util.Int32(data[padding:])
	padding += 4

	value, valuePadding, actualValueSize := ResolveType(typeName, data[padding:], false, depth+1)

	if actualValueSize != int(valueSize) {
		logrus.Errorf("%s%s read %d bytes, expected %d (%d)\n", strings.Repeat(" ", depth), typeName, actualValueSize, valueSize, int(valueSize)-actualValueSize)
	}

	return Property{
		Name:  name,
		Type:  typeName,
		Index: keyIndex,
		Size:  valueSize,
		Value: value,
	}, padding + valuePadding
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

func ParseIntProperty(data []byte, inArray bool) (int32, int, int) {
	if inArray {
		return util.Int32(data), 4, 4
	}

	return util.Int32(data[1:]), 5, 4
}

func ParseInt8Property(data []byte, inArray bool) (int8, int, int) {
	if inArray {
		return int8(data[0]), 1, 1
	}

	return int8(data[1]), 2, 1
}

func ProcessIntProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, inArray bool) (int, int) {
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
	case string:
		// No way to preserve reference to an incorrect type
		tempInt, err := strconv.Atoi((*target.(*ReadOrWritable)).(string))

		if err != nil {
			logrus.Panic(err)
		}

		tempInt32 := int32(tempInt)
		targetInt = &tempInt32
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

func ParseBoolProperty(data []byte) (bool, int, int) {
	return data[0] > 0, 2, 0
}

func ProcessBoolProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	var targetBool = (*target.(*ReadOrWritable)).(*bool)
	if buf == nil {
		*targetBool = *data.At(0) > 0
	} else {
		if *targetBool {
			buf.Write([]byte{0x01, 0x00})
		} else {
			buf.Write([]byte{0x00, 0x00})
		}
	}

	return 2, 0
}

func ParseByteProperty(data []byte) (ByteProperty, int, int) {
	padding := 0

	enumType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	if enumType == "None" {
		return ByteProperty{
			EnumType: enumType,
			Byte:     data[padding],
		}, padding + 1, 1
	} else {
		enumName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength
		return ByteProperty{
			EnumType: enumType,
			EnumName: &enumName,
		}, padding, 4 + strLength
	}
}

func ProcessByteProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	var targetByte = (*target.(*ReadOrWritable)).(*ByteProperty)

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
		strLength := util.RoWInt32StringNull(data.From(padding), targetByte.EnumName, buf)
		padding += 4 + strLength
		return padding, 4 + strLength
	}
}

func ParseEnumProperty(data []byte) (EnumProperty, int, int) {
	padding := 0

	enumProperty := EnumProperty{}

	var strLength, enumNameLength int
	enumProperty.Type, strLength = util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	enumProperty.Name, enumNameLength = util.Int32StringNull(data[padding:])
	padding += 4 + enumNameLength

	return enumProperty, padding, enumNameLength + 4
}

func ProcessEnumProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	var targetEnum = (*target.(*ReadOrWritable)).(*EnumProperty)

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

func ParseFloatProperty(data []byte) (float32, int, int) {
	return util.Float32(data[1:]), 5, 4
}

func ProcessFloatProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	// Skip byte
	var skip byte = 0x00
	util.RoWB(data.At(0), &skip, buf)

	var targetFloat = (*target.(*ReadOrWritable)).(*float32)
	util.RoWFloat32(data.From(1), targetFloat, buf)

	return 5, 4
}

func ParseStringProperty(data []byte) (string, int, int) {
	str, strLength := util.Int32StringNull(data[1:])
	return str, 1 + 4 + strLength, 4 + strLength
}

func ProcessStringProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
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

func ParseTextProperty(data []byte) (TextProperty, int, int) {
	padding := 1

	textProperty := TextProperty{}

	// TODO Unknown
	textProperty.Magic = data[padding : padding+13]
	padding += 13

	var strLength int
	textProperty.String, strLength = util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	return textProperty, padding, padding - 1
}

func ProcessTextProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer) (int, int) {
	var targetText = (*target.(*ReadOrWritable)).(*TextProperty)

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

func ParseObjectProperty(data []byte, inArray bool) (ObjectProperty, int, int) {
	padding := 0
	overhead := 0

	if !inArray {
		overhead += 1
	}

	objectProperty := ObjectProperty{}

	var strLength int
	objectProperty.World, strLength = util.Int32StringNull(data[padding+overhead:])
	padding += 4 + strLength

	objectProperty.Class, strLength = util.Int32StringNull(data[padding+overhead:])
	padding += 4 + strLength

	return objectProperty, padding + overhead, padding
}

func ProcessObjectProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, inArray bool) (int, int) {
	var targetObject = (*target.(*ReadOrWritable)).(*ObjectProperty)

	padding := 0
	overhead := 0

	if !inArray {
		// Skip byte
		var skip byte = 0x00
		util.RoWB(data.At(0), &skip, buf)
		overhead += 1
	}

	padding += util.RoWInt32StringNull(data.From(padding), &targetObject.World, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &targetObject.Class, buf)
	padding += 4

	return padding + overhead, padding
}

func ParseArrayProperty(data []byte, depth int) (ArrayProperty, int, int) {
	padding := 0

	typeName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	valueCount := int(util.Int32(data[padding:]))
	padding += 4

	var structName, structType, structClassType string
	var structSize int32
	var magic1, magic2 []byte

	beforePadding := padding

	logrus.Debug(strings.Repeat(" ", depth), "[", typeName, "] x ", valueCount)

	if typeName == "StructProperty" {
		structName, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		structType, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		structSize = util.Int32(data[padding:])
		padding += 4

		// TODO Unknown
		magic1 = data[padding : padding+4]
		padding += 4

		structClassType, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		// TODO Unknown
		magic2 = data[padding : padding+17]
		padding += 17
	}

	var values []ReadOrWritable
	valueSizes := padding - beforePadding

	for i := 0; i < valueCount; i++ {
		if typeName == "StructProperty" {
			value, padded, _ := ParseStructProperty(data[padding:], &structClassType, depth+1)
			padding += padded
			valueSizes += padded
			values = append(values, value)
		} else {
			value, padded, valueSize := ResolveType(typeName, data[padding:], true, depth+1)
			padding += padded
			valueSizes += valueSize
			values = append(values, value)
		}
	}

	if typeName == "StructProperty" {
		return ArrayProperty{
			Type:            typeName,
			Values:          values,
			StructName:      structName,
			StructType:      structType,
			StructSize:      structSize,
			Magic1:          magic1,
			StructClassType: structClassType,
			Magic2:          magic2,
		}, padding, valueSizes + 4
	}

	return ArrayProperty{
		Type:   typeName,
		Values: values,
	}, padding, valueSizes + 4
}

func ProcessArrayProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, depth int) (int, int) {
	var targetArray = (*target.(*ReadOrWritable)).(*ArrayProperty)

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

func ParseStructProperty(data []byte, arrayTypeName *string, depth int) (StructProperty, int, int) {
	padding := 0

	var typeName string
	var magic []byte

	if arrayTypeName == nil {
		newTypeName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength
		typeName = newTypeName

		// TODO Unknown
		magic = data[padding : padding+17]
		padding += 17
	} else {
		typeName = *arrayTypeName
	}

	logrus.Debug(strings.Repeat(" ", depth), typeName)

	structProperty := StructProperty{
		Type:  typeName,
		Magic: &magic,
	}

	beforePadding := padding

	switch typeName {
	case "Vector":
		fallthrough
	case "Rotator":
		vec3 := util.Vec3(data[padding:])
		padding += 12

		structProperty.Value = vec3

		return structProperty, padding, padding - beforePadding
	case "Color":
		structProperty.Value = Color{
			R: data[padding],
			G: data[padding+1],
			B: data[padding+2],
			A: data[padding+3],
		}

		return structProperty, padding + 4, 4
	case "LinearColor":
		vec4 := util.Vec4(data[padding:])
		padding += 16

		structProperty.Value = LinearColor{
			R: vec4.X,
			G: vec4.Y,
			B: vec4.Z,
			A: vec4.W,
		}

		return structProperty, padding, padding - beforePadding
	case "Quat":
		vec4 := util.Vec4(data[padding:])
		padding += 16

		structProperty.Value = vec4

		return structProperty, padding, padding - beforePadding
	case "Box":
		min := util.Vec3(data[padding:])
		padding += 12

		max := util.Vec3(data[padding:])
		padding += 12

		valid := data[padding]
		padding += 1

		structProperty.Value = Box{
			Min:   min,
			Max:   max,
			Valid: valid,
		}

		return structProperty, padding, padding - beforePadding
	case "InventoryItem":
		magic, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		itemName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		levelName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		pathName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		structProperty.Value = InventoryItem{
			Magic:     magic,
			ItemName:  itemName,
			LevelName: levelName,
			PathName:  pathName,
		}

		return structProperty, padding, padding - beforePadding
	case "RailroadTrackPosition":
		world, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		entityType, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		offset := util.Float32(data[padding:])
		padding += 4

		forward := util.Float32(data[padding:])
		padding += 4

		structProperty.Value = RailroadTrackPosition{
			World:      world,
			EntityType: entityType,
			Offset:     offset,
			Forward:    forward,
		}

		return structProperty, padding, padding - beforePadding
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
		values, valueSize := ReadToNone(data[padding:], depth+1)
		padding += valueSize

		structProperty.Value = GenericStruct{
			Values: values,
		}

		return structProperty, padding, valueSize
	}

	logrus.Panicf("Unknown struct: %s - %#v\n", typeName, string(data[padding:]))

	panic(1) // Logrus will panic for us
}

func ProcessStructProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, arrayTypeName *string, depth int) (int, int) {
	var targetStruct = (*target.(*ReadOrWritable)).(*StructProperty)

	padding := 0

	var typeName string

	if arrayTypeName == nil {
		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Type, buf)
		padding += 4
		typeName = targetStruct.Type

		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+17), targetStruct.Magic, buf)
		padding += 17
	} else {
		typeName = *arrayTypeName
	}

	logrus.Debug(strings.Repeat(" ", depth), typeName)

	beforePadding := padding

	// TODO Merge reading/writing
	switch typeName {
	case "Vector":
		fallthrough
	case "Rotator":
		tempCast := targetStruct.Value.(util.Vector3)

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
		tempCast := targetStruct.Value.(util.Vector4)

		util.RoWVec4(data.From(padding), &tempCast, buf)
		padding += 16

		return padding, padding - beforePadding
	case "Box":
		tempCast := targetStruct.Value.(Box)

		util.RoWVec3(data.From(padding), &tempCast.Min, buf)
		padding += 12

		util.RoWVec3(data.From(padding), &tempCast.Max, buf)
		padding += 12

		util.RoWB(data.At(padding), &tempCast.Valid, buf)
		padding += 1

		return padding, padding - beforePadding
	case "InventoryItem":
		tempCast := targetStruct.Value.(InventoryItem)

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
		tempCast := targetStruct.Value.(RailroadTrackPosition)

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

func ParseMapProperty(data []byte, depth int) (MapProperty, int, int) {
	padding := 0

	keyType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	valueType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	beforePadding := padding

	// TODO Unknown
	magic := data[padding : padding+4]
	padding += 4

	pairCount := int(util.Int32(data[padding:]))
	padding += 4

	values := make([]MapEntry, pairCount)

	for i := 0; i < pairCount; i++ {
		key, keySize, _ := ResolveType(keyType, data[padding:], true, depth+1)
		padding += keySize

		innerValues, padded := ReadToNone(data[padding:], depth+1)
		padding += padded

		values[i] = MapEntry{
			Key:   key,
			Value: innerValues,
		}
	}

	return MapProperty{
		KeyType:   keyType,
		ValueType: valueType,
		Magic:     magic,
		Values:    values,
	}, padding, padding - beforePadding
}

func ProcessMapProperty(data util.RawHolder, target ReadOrWritable, buf *bytes.Buffer, depth int) (int, int) {
	var targetMap = (*target.(*ReadOrWritable)).(*MapProperty)

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
			fmt.Printf("%#v\n", entry.Key)
			y := InterfaceIfy(entry.Key).(ReadOrWritable)
			fmt.Printf("%#v\n", y)
			fmt.Println(targetMap.KeyType)

			padded, _ := ProcessType(targetMap.KeyType, data.FromNew(padding), &y, buf, true, depth+1)
			padding += padded

			// return padding, padding - beforePadding

			padding += RoWToNone(data.FromNew(padding), &entry.Value, buf, depth+1)
		}
	}

	return padding, padding - beforePadding
}

// Absolutely necessary hack
func InterfaceIfy(data interface{}) interface{} {
	return data
}
