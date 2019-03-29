package save

import (
	"github.com/sirupsen/logrus"
	"satisfactory-tool/util"
)

type Parsable interface {
	Parse(length int, data []byte)
}

type SaveComponentType struct {
	ClassType        string
	EntityType       string
	InstanceType     string
	ParentEntityType string
	Fields           [][]map[string]interface{}
	// SaveObjectCount  int32
}

type EntityType struct {
	ClassType        string
	EntityType       string
	InstanceType     string
	MagicInt1        int32
	MagicInt2        int32
	Rotation         util.Vector4
	Position         util.Vector3
	Scale            util.Vector3
	ParentObjectRoot string
	ParentObjectName string
	Components       [][]string
	Fields           [][]map[string]interface{}
	Extra            interface{}
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
	saveComponentType.Fields, _ = ReReadToZero(data, 0)
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

		entityType.Components = append(entityType.Components, []string{root, name})
	}

	innerValues, padded := ReadToNone(data[padding:], 0)
	padding += padded

	entityType.Fields = append(entityType.Fields, innerValues)

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
			}
		} else {
			logrus.Errorf("%v has >4 bytes [%d] left and is not handled as a special case!\n", entityType.ClassType, length-padding)
		}
	}
}
