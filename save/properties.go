package save

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"satisfactory-tool/util"
	"strings"
)

func ResolveType(typeName string, data []byte, inArray bool, depth int) (interface{}, int, int) {
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

func ParseProperty(data []byte, depth int) (string, string, int, interface{}, int, int) {
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

func ParseByteProperty(data []byte) (interface{}, int, int) {
	padding := 0

	enumType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	if enumType == "None" {
		return data[padding], padding + 1, 1
	} else {
		enumName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength
		return enumType + ":" + enumName, padding, 4 + strLength
	}
}

func ParseEnumProperty(data []byte) (string, int, int) {
	padding := 0

	enumType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	enumName, enumNameLength := util.Int32StringNull(data[padding:])
	padding += 4 + enumNameLength

	return enumType + ":" + enumName, padding, enumNameLength + 4
}

func ParseFloatProperty(data []byte) (float32, int, int) {
	return util.Float32(data[1:]), 5, 4
}

func ParseStringProperty(data []byte) (string, int, int) {
	str, strLength := util.Int32StringNull(data[1:])
	return str, 1 + 4 + strLength, 4 + strLength
}

func ParseTextProperty(data []byte) (string, int, int) {
	padding := 1

	// TODO Unknown
	padding += 13

	str, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	return str, padding, padding - 1
}

func ParseObjectProperty(data []byte, inArray bool) ([]string, int, int) {
	padding := 0
	overhead := 0

	if !inArray {
		overhead += 1
	}

	string1, strLength := util.Int32StringNull(data[padding+overhead:])
	padding += 4 + strLength

	string2, strLength := util.Int32StringNull(data[padding+overhead:])
	padding += 4 + strLength

	return []string{string1, string2}, padding + overhead, padding
}

func ParseArrayProperty(data []byte, depth int) (interface{}, int, int) {
	padding := 0

	typeName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	valueCount := int(util.Int32(data[padding:]))
	padding += 4

	var structName, structType, structClassType string

	beforePadding := padding

	logrus.Debug(strings.Repeat(" ", depth), "[", typeName, "] x ", valueCount)

	if typeName == "StructProperty" {
		structName, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		structType, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		// structSize := int(util.Int32(data[padding:]))
		padding += 4

		// TODO Unknown
		padding += 4

		structClassType, strLength = util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		// fmt.Println(structName, structType, structClassType)

		// TODO Unknown
		padding += 17
	}

	var values []interface{}
	valueSizes := padding - beforePadding
	// fmt.Println(valueSizes)

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

	// fmt.Println(valueSizes, padding, valueCount)

	if typeName == "StructProperty" {
		return map[string]interface{}{
			"type":            typeName,
			"values":          values,
			"structName":      structName,
			"structType":      structType,
			"structClassType": structClassType,
		}, padding, valueSizes + 4
	}

	return map[string]interface{}{
		"type":   typeName,
		"values": values,
	}, padding, valueSizes + 4
}

func ParseStructProperty(data []byte, arrayTypeName *string, depth int) (interface{}, int, int) {
	padding := 0

	var typeName string

	if arrayTypeName == nil {
		newTypeName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength
		typeName = newTypeName

		// Skip 4 x int32 + 1 byte TODO Unknown
		padding += 17
	} else {
		typeName = *arrayTypeName
	}

	logrus.Debug(strings.Repeat(" ", depth), typeName)

	switch typeName {
	case "Vector":
		fallthrough
	case "Rotator":
		vec3 := util.Vec3(data[padding:])
		padding += 12

		return map[string]interface{}{
			"type": typeName,
			"x":    vec3.X,
			"y":    vec3.Y,
			"z":    vec3.Z,
		}, padding, 12
	case "LinearColor":
		vec4 := util.Vec4(data[padding:])
		padding += 16

		return map[string]interface{}{
			"type": typeName,
			"r":    vec4.X,
			"g":    vec4.Y,
			"b":    vec4.Z,
			"a":    vec4.W,
		}, padding, 16
	case "Quat":
		vec4 := util.Vec4(data[padding:])
		padding += 16

		return map[string]interface{}{
			"type": typeName,
			"a":    vec4.X,
			"b":    vec4.Y,
			"c":    vec4.Z,
			"d":    vec4.W,
		}, padding, 16
	case "Box":
		min := util.Vec3(data[padding:])
		padding += 12

		max := util.Vec3(data[padding:])
		padding += 12

		valid := data[padding]
		padding += 1

		return map[string]interface{}{
			"type":  typeName,
			"min":   min,
			"max":   max,
			"valid": valid,
		}, padding, 25
	case "InventoryItem":
		beforePadding := padding

		magic, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		itemName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		levelName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		pathName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		values, _ := ReadToNone(data[padding:], depth+1)

		return map[string]interface{}{
			"type":      typeName,
			"magic":     magic,
			"itemName":  itemName,
			"levelName": levelName,
			"pathName":  pathName,
			"values":    values,
		}, padding, padding - beforePadding
	case "RailroadTrackPosition":
		beforePadding := padding

		world, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		entityType, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		offset := util.Float32(data[padding:])
		padding += 4

		forward := util.Float32(data[padding:])
		padding += 4

		return map[string]interface{}{
			"type":       typeName,
			"world":      world,
			"entityType": entityType,
			"offset":     offset,
			"forward":    forward,
		}, padding, padding - beforePadding
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

		return map[string]interface{}{
			"type":   typeName,
			"values": values,
		}, padding, valueSize
	}

	// fmt.Println(padding)
	// fmt.Printf("%#v\n", string(data[padding:padding+100]))

	logrus.Panic("Unknown struct: " + typeName)

	// fmt.Println("Unknown struct: " + typeName)
	return nil, 0, 0
}

func ParseMapProperty(data []byte, depth int) (interface{}, int, int) {
	padding := 0

	keyType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	valueType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// TODO Unknown
	padding += 1

	beforePadding := padding

	// TODO Unknown
	padding += 4

	pairCount := int(util.Int32(data[padding:]))
	padding += 4

	values := make(map[string][]map[string]interface{})

	for i := 0; i < pairCount; i++ {
		key, keySize, _ := ResolveType(keyType, data[padding:], true, depth+1)
		padding += keySize

		name := ""
		innerValues := make([]map[string]interface{}, 0)

		for name != "None" {
			propName, typeName, _, value, _, padded := ParseProperty(data[padding:], depth+1)
			name = propName
			padding += padded

			if propName != "None" {
				innerValues = append(innerValues, map[string]interface{}{
					"name":  propName,
					"type":  typeName,
					"value": value,
				})
			}
		}

		// JSON Compatibility
		values[fmt.Sprintf("%v", key)] = innerValues
	}

	return map[string]interface{}{
		"keyType":   keyType,
		"valueType": valueType,
		"values":    values,
	}, padding, padding - beforePadding
}
