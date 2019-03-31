package save

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"satisfactory-tool/util"
)

type ParsableWrapper struct {
	Type string   `json:"type"`
	Data Parsable `json:"data"`
}

func (wrapper *ParsableWrapper) UnmarshalJSON(b []byte) error {
	var temp map[string]json.RawMessage
	err := json.Unmarshal(b, &temp)

	if err != nil {
		return err
	}

	err = json.Unmarshal(temp["type"], &wrapper.Type)

	if err != nil {
		return err
	}

	switch wrapper.Type {
	case "save":
		data := SaveComponentType{}
		err = json.Unmarshal(temp["data"], &data)
		wrapper.Data = &data
		break
	case "entity":
		data := EntityType{}
		err = json.Unmarshal(temp["data"], &data)
		wrapper.Data = &data
		break
	default:
		logrus.Panic("Unknown Type: " + wrapper.Type)
	}

	if err != nil {
		return err
	}

	return nil
}

type Parsable interface {
	Parse(length int, data []byte)
}

type SaveComponentType struct {
	ClassType        string                   `json:"class_type"`
	EntityType       string                   `json:"entity_type"`
	InstanceType     string                   `json:"instance_type"`
	ParentEntityType string                   `json:"parent_entity_type"`
	Fields           []map[string]interface{} `json:"fields"`
}

type EntityType struct {
	ClassType        string                   `json:"class_type"`
	EntityType       string                   `json:"entity_type"`
	InstanceType     string                   `json:"instance_type"`
	MagicInt1        int32                    `json:"magic_int_1"`
	MagicInt2        int32                    `json:"magic_int_2"`
	Rotation         util.Vector4             `json:"rotation"`
	Position         util.Vector3             `json:"position"`
	Scale            util.Vector3             `json:"scale"`
	ParentObjectRoot string                   `json:"parent_object_root"`
	ParentObjectName string                   `json:"parent_object_name"`
	Components       [][]string               `json:"components"`
	Fields           []map[string]interface{} `json:"fields"`
	Extra            interface{}              `json:"extra,omitempty"`
}

func ParseSaveComponentType(saveData []byte) (*SaveComponentType, int) {
	padding := 0

	classType, strLen := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLen

	entityType, strLen := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLen

	instanceType, strLen := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLen

	parentEntityType, strLen := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLen

	return &SaveComponentType{
		ClassType:        classType,
		EntityType:       entityType,
		InstanceType:     instanceType,
		ParentEntityType: parentEntityType,
	}, padding
}

func ParseEntityType(saveData []byte) (*EntityType, int) {
	padding := 0

	classType, strLen := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLen

	entityType, strLen := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLen

	instanceType, strLen := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLen

	magicInt1 := util.Int32(saveData[padding:])
	padding += 4

	rotation := util.Vec4(saveData[padding:])
	padding += 16

	position := util.Vec3(saveData[padding:])
	padding += 12

	scale := util.Vec3(saveData[padding:])
	padding += 12

	magicInt2 := util.Int32(saveData[padding:])
	padding += 4

	return &EntityType{
		ClassType:    classType,
		EntityType:   entityType,
		InstanceType: instanceType,
		MagicInt1:    magicInt1,
		MagicInt2:    magicInt2,
		Rotation:     rotation,
		Position:     position,
		Scale:        scale,
	}, padding
}

func (saveComponentType *SaveComponentType) Parse(length int, data []byte) {
	padding := 0

	saveComponentType.Fields, padding = ReadToNone(data, 0)

	if length-padding > 4 {
		logrus.Errorf("%v has >4 bytes [%d] left and is not handled as a special case!\n", saveComponentType.ClassType, length-padding)
	}
}

func (entityType *EntityType) Parse(length int, data []byte) {
	padding := 0
	var strLen int

	entityType.ParentObjectRoot, strLen = util.Int32StringNull(data[padding:])
	padding += 4 + strLen

	entityType.ParentObjectName, strLen = util.Int32StringNull(data[padding:])
	padding += 4 + strLen

	componentCount := int(util.Int32(data[padding:]))
	padding += 4

	entityType.Components = make([][]string, componentCount)

	for i := 0; i < componentCount; i++ {
		root, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		name, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		entityType.Components[i] = []string{root, name}
	}

	var padded int
	entityType.Fields, padded = ReadToNone(data[padding:], 0)
	padding += padded

	if length-padding > 4 {
		if specialFunc, ok := specialClasses[entityType.ClassType]; ok {
			extraData, padded := specialFunc(data[padding:length])
			padding += padded

			if extraData == nil {
				logrus.Errorf("%v Did not process any data [%d]\n", entityType.ClassType, length-padding)
			} else if length-padding > 4 {
				logrus.Errorf("%v Did not read to end [%d]\n", entityType.ClassType, length-padding)
			} else {
				entityType.Extra = extraData
				return
			}
		} else {
			logrus.Errorf("%v has >4 bytes [%d] left and is not handled as a special case!\n", entityType.ClassType, length-padding)
		}
	}

	entityType.Extra = data[padding:]
}
