package save

import (
	"encoding/binary"
	"fmt"
	"math"
	"satisfactory-tool/util"
	"strings"
)

func ResolveType(typeName string, data []byte, inArray bool) (interface{}, int) {
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
		return ParseArrayProperty(data)
	case "MapProperty":
		return ParseMapProperty(data)
	case "StructProperty":
		return ParseStructProperty(data)
	}

	panic("Don't know how to process: " + typeName)
}

func ParseProperty(data []byte) (string, string, int, interface{}, int, int) {
	padding := 0

	name, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	if name == "None" || name == "" {
		return "None", "", 0, "None", 0, padding
	}

	typeName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	valueSize := int(util.Int32(data[padding:]))
	padding += 4

	keyIndex := int(util.Int32(data[padding:]))
	padding += 4

	value, valuePadding := ResolveType(typeName, data[padding:], false)

	return name, typeName, keyIndex, value, valueSize, padding + valuePadding
}

func ParseIntProperty(data []byte, inArray bool) (int32, int) {
	if inArray {
		return util.Int32(data), 4
	}
	return util.Int32(data[1:]), 5
}

func ParseBoolProperty(data []byte) (bool, int) {
	return util.Int16(data) > 0, 2
}

func ParseByteProperty(data []byte) (byte, int) {
	padding := 0

	isNone, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	if isNone == "None" {
		return data[padding], padding + 1
	} else {
		// TODO Check if correct
		actualBoolean, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength
		return strings.ToLower(actualBoolean)[0], padding + 1
	}
}

func ParseEnumProperty(data []byte) (string, int) {
	padding := 0

	enumType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	enumName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	return enumType + ":" + enumName, padding
}

func ParseFloatProperty(data []byte) (float32, int) {
	return math.Float32frombits(binary.LittleEndian.Uint32(data[1:5])), 5
}

func ParseStringProperty(data []byte) (string, int) {
	str, strLength := util.Int32StringNull(data[1:])
	return str, 1 + 4 + strLength
}

func ParseTextProperty(data []byte) (string, int) {
	padding := 1

	// TODO Unknown
	padding += 13

	str, strLength := util.Int32StringNull(data[padding:])
	return str, padding + 4 + strLength
}

func ParseObjectProperty(data []byte, inArray bool) ([]string, int) {
	padding := 0

	if !inArray {
		padding += 1
	}

	string1, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	string2, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	return []string{string1, string2}, padding
}

func ParseArrayProperty(data []byte) (interface{}, int) {
	padding := 0

	typeName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip byte
	padding += 1

	valueCount := int(util.Int32(data[padding:]))
	padding += 4

	var structName, structType, structClassType string

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

		// TODO Unknown
		padding += 17
	}

	var values []interface{}

	for i := 0; i < valueCount; i++ {
		if typeName == "StructProperty" {
			_, _, _, value, _, padded := ParseProperty(data[padding:])
			padding += padded
			values = append(values, value)
		} else {
			value, padded := ResolveType(typeName, data[padding:], true)
			padding += padded
			values = append(values, value)
		}
	}

	if typeName == "StructProperty" {
		return map[string]interface{}{
			"type":            typeName,
			"values":          values,
			"structName":      structName,
			"structType":      structType,
			"structClassType": structClassType,
		}, padding
	}

	return map[string]interface{}{
		"type":   typeName,
		"values": values,
	}, padding
}

func ParseStructProperty(data []byte) (interface{}, int) {
	padding := 0

	typeName, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// Skip 4 x int32 + 1 byte TODO Unknown
	padding += 17

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
		}, padding
	case "LinearColor":
		vec4 := util.Vec4(data[padding:])
		padding += 16

		return map[string]interface{}{
			"type": typeName,
			"r":    vec4.X,
			"g":    vec4.Y,
			"b":    vec4.Z,
			"a":    vec4.W,
		}, padding
	case "Quat":
		vec4 := util.Vec4(data[padding:])
		padding += 16

		return map[string]interface{}{
			"type": typeName,
			"a":    vec4.X,
			"b":    vec4.Y,
			"c":    vec4.Z,
			"d":    vec4.W,
		}, padding
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
		}, padding
	case "InventoryItem":
		magic, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		itemName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		levelName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		pathName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		name := ""
		values := make([]interface{}, 0)

		for name != "None" {
			propName, typeName, _, value, _, padded := ParseProperty(data[padding:])
			name = propName
			padding += padded

			if propName != "None" {
				values = append(values, map[string]interface{}{
					"name":  propName,
					"type":  typeName,
					"value": value,
				})
			}
		}

		return map[string]interface{}{
			"type":      typeName,
			"magic":     magic,
			"itemName":  itemName,
			"levelName": levelName,
			"pathName":  pathName,
			"values":    values,
		}, padding
	case "InventoryStack":
		fallthrough
	case "RemovedInstanceArray":
		fallthrough
	case "Transform":
		name := ""
		values := make([]interface{}, 0)

		for name != "None" {
			propName, typeName, _, value, _, padded := ParseProperty(data[padding:])
			name = propName
			padding += padded

			if propName != "None" {
				values = append(values, map[string]interface{}{
					"name":  propName,
					"type":  typeName,
					"value": value,
				})
			}
		}

		return map[string]interface{}{
			"type":   typeName,
			"values": values,
		}, padding
	}

	fmt.Println(padding)
	fmt.Printf("%#v\n", string(data[padding:padding+100]))

	panic("Unknown struct: " + typeName)
}

func ParseMapProperty(data []byte) (interface{}, int) {
	padding := 0

	keyType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	valueType, strLength := util.Int32StringNull(data[padding:])
	padding += 4 + strLength

	// TODO Unknown
	padding += 5

	pairCount := int(util.Int32(data[padding:]))
	padding += 4

	values := make(map[string][]map[string]interface{})

	for i := 0; i < pairCount; i++ {
		key, keySize := ResolveType(keyType, data[padding:], true)
		padding += keySize

		name := ""
		innerValues := make([]map[string]interface{}, 0)

		for name != "None" {
			propName, typeName, _, value, _, padded := ParseProperty(data[padding:])
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
	}, padding
}
