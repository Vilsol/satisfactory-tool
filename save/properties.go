package save

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"satisfactory-tool/util"
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
	case "ObjectProperty":
		return ParseObjectProperty(data, inArray)
	case "ArrayProperty":
		return ParseArrayProperty(data, depth)
	case "MapProperty":
		return ParseMapProperty(data, depth)
	case "StructProperty":
		return ParseStructProperty(data, nil, depth)
	}

	panic("Don't know how to process: " + typeName)
}

func ParseProperty(data []byte, depth int) (string, string, int, ReadOrWritable, int, int) {
	padding := 0

	name, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	if name == "None" || name == "" {
		return "None", "", 0, "None", 0, padding
	}

	typeName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	logrus.Debug(strings.Repeat(" ", depth), typeName, " ", name)

	valueSize := int(util.Int32(data[padding:]))
	padding += 4

	keyIndex := int(util.Int32(data[padding:]))
	padding += 4

	value, valuePadding, actualValueSize := ResolveType(typeName, data[padding:], false, depth+1)

	if actualValueSize != valueSize {
		logrus.Errorf("%s%s read %d bytes, expected %d (%d)\n", strings.Repeat(" ", depth), typeName, actualValueSize, valueSize, valueSize-actualValueSize)
	}

	return name, typeName, keyIndex, value, valueSize, padding + valuePadding
}

func ParseIntProperty(data []byte, inArray bool) (int32, int, int) {
	if inArray {
		return util.Int32(data), 4, 4
	}

	return util.Int32(data[1:]), 5, 4
}

func ParseBoolProperty(data []byte) (bool, int, int) {
	return data[1] > 0, 2, 0
}

func ParseByteProperty(data []byte) (ByteProperty, int, int) {
	padding := 0

	enumType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	if enumType == "None" {
		return ByteProperty{
			Byte: data[padding],
		}, padding + 1, 1
	} else {
		enumName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength
		return ByteProperty{
			EnumType: &enumType,
			EnumName: &enumName,
		}, padding, 4 + strLength
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

func ParseFloatProperty(data []byte) (float32, int, int) {
	return util.Float32(data[1:]), 5, 4
}

func ParseStringProperty(data []byte) (string, int, int) {
	str, strLength := util.Int32StringNull(data[1:])
	return str, 1 + 4 + strLength, 4 + strLength
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

func ParseArrayProperty(data []byte, depth int) (ArrayProperty, int, int) {
	padding := 0

	typeName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	valueCount := int(util.Int32(data[padding:]))
	padding += 4

	var structName, structType, structClassType string
	var structSize int
	var magic1, magic2 []byte

	beforePadding := padding

	logrus.Debug(strings.Repeat(" ", depth), "[", typeName, "] x ", valueCount)

	if typeName == "StructProperty" {
		structName, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		structType, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		structSize = int(util.Int32(data[padding:]))
		padding += 4

		// TODO Unknown
		magic1 = data[padding : padding+4]
		padding += 4

		structClassType, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		// TODO Unknown
		magic1 = data[padding : padding+17]
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
			StructName:      &structName,
			StructType:      &structType,
			StructSize:      &structSize,
			Magic1:          &magic1,
			StructClassType: &structClassType,
			Magic2:          &magic2,
		}, padding, valueSizes + 4
	}

	return ArrayProperty{
		Type:   typeName,
		Values: values,
	}, padding, valueSizes + 4
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

		structProperty.Value = map[string]interface{}{
			"x": vec3.X,
			"y": vec3.Y,
			"z": vec3.Z,
		}

		return structProperty, padding, padding - beforePadding
	case "Color":
		structProperty.Value = map[string]interface{}{
			"r": data[padding],
			"g": data[padding+1],
			"b": data[padding+2],
			"a": data[padding+3],
		}

		return structProperty, padding + 4, 4
	case "LinearColor":
		vec4 := util.Vec4(data[padding:])
		padding += 16

		structProperty.Value = map[string]interface{}{
			"r": vec4.X,
			"g": vec4.Y,
			"b": vec4.Z,
			"a": vec4.W,
		}

		return structProperty, padding, padding - beforePadding
	case "Quat":
		vec4 := util.Vec4(data[padding:])
		padding += 16

		structProperty.Value = map[string]interface{}{
			"a": vec4.X,
			"b": vec4.Y,
			"c": vec4.Z,
			"d": vec4.W,
		}

		return structProperty, padding, padding - beforePadding
	case "Box":
		min := util.Vec3(data[padding:])
		padding += 12

		max := util.Vec3(data[padding:])
		padding += 12

		valid := data[padding]
		padding += 1

		structProperty.Value = map[string]interface{}{
			"min":   min,
			"max":   max,
			"valid": valid,
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

		values, _ := ReadToNone(data[padding:], depth+1)

		structProperty.Value = map[string]interface{}{
			"magic":     magic,
			"itemName":  itemName,
			"levelName": levelName,
			"pathName":  pathName,
			"values":    values,
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

		structProperty.Value = map[string]interface{}{
			"world":      world,
			"entityType": entityType,
			"offset":     offset,
			"forward":    forward,
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

		structProperty.Value = map[string]interface{}{
			"values": values,
		}

		return structProperty, padding, valueSize
	}

	logrus.Panicf("Unknown struct: %s - %#v\n", typeName, string(data[padding:]))

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

	values := make(map[string][]map[string]interface{})

	for i := 0; i < pairCount; i++ {
		key, keySize, _ := ResolveType(keyType, data[padding:], true, depth+1)
		padding += keySize

		innerValues, padded := ReadToNone(data[padding:], depth+1)
		padding += padded

		// JSON Compatibility
		values[fmt.Sprintf("%v", key)] = innerValues
	}

	return MapProperty{
		KeyType:   keyType,
		ValueType: valueType,
		Magic:     magic,
		Values:    values,
	}, padding, padding - beforePadding
}
