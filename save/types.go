package save

import (
	"bytes"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"satisfactory-tool/util"
)

type ParsableWrapper struct {
	Type   string   `json:"type"`
	Length int32    `json:"length"`
	Data   Parsable `json:"data"`
}

func (wrapper *ParsableWrapper) UnmarshalJSON(b []byte) error {
	var temp map[string]jsoniter.RawMessage
	err := jsoniter.Unmarshal(b, &temp)

	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(temp["type"], &wrapper.Type)
	err = jsoniter.Unmarshal(temp["length"], &wrapper.Length)

	if err != nil {
		return err
	}

	switch wrapper.Type {
	case "save":
		data := SaveComponentType{}
		err = jsoniter.Unmarshal(temp["data"], &data)
		wrapper.Data = &data
		break
	case "entity":
		data := EntityType{}
		err = jsoniter.Unmarshal(temp["data"], &data)
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
	Process(data util.RawHolder, component *Parsable, buf *bytes.Buffer) int
}

type SaveComponentType struct {
	ClassType        string     `json:"class_type"`
	EntityType       string     `json:"entity_type"`
	InstanceType     string     `json:"instance_type"`
	ParentEntityType string     `json:"parent_entity_type"`
	Fields           []Property `json:"fields"`
}

type EntityType struct {
	ClassType        string           `json:"class_type"`
	EntityType       string           `json:"entity_type"`
	InstanceType     string           `json:"instance_type"`
	MagicInt1        int32            `json:"magic_int_1"`
	MagicInt2        int32            `json:"magic_int_2"`
	Rotation         util.Vector4     `json:"rotation"`
	Position         util.Vector3     `json:"position"`
	Scale            util.Vector3     `json:"scale"`
	ParentObjectRoot string           `json:"parent_object_root"`
	ParentObjectName string           `json:"parent_object_name"`
	Components       []ObjectProperty `json:"components"`
	Fields           []Property       `json:"fields"`
	ExtraObjects     interface{}      `json:"extra_objects,omitempty"`
	Extra            []byte           `json:"extra,omitempty"`
}

func ProcessSaveComponentType(data util.RawHolder, component *SaveComponentType, buf *bytes.Buffer) int {
	padding := 0

	padding += util.RoWInt32StringNull(data.From(padding), &component.ClassType, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &component.EntityType, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &component.InstanceType, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &component.ParentEntityType, buf)
	padding += 4

	return padding
}

func ProcessEntityType(data util.RawHolder, component *EntityType, buf *bytes.Buffer) int {
	padding := 0

	padding += util.RoWInt32StringNull(data.From(padding), &component.ClassType, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &component.EntityType, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &component.InstanceType, buf)
	padding += 4

	util.RoWInt32(data.From(padding), &component.MagicInt1, buf)
	padding += 4

	util.RoWVec4(data.From(padding), &component.Rotation, buf)
	padding += 16

	util.RoWVec3(data.From(padding), &component.Position, buf)
	padding += 12

	util.RoWVec3(data.From(padding), &component.Scale, buf)
	padding += 12

	util.RoWInt32(data.From(padding), &component.MagicInt2, buf)
	padding += 4

	return padding
}

func (saveComponentType *SaveComponentType) Process(data util.RawHolder, component *Parsable, buf *bytes.Buffer) int {
	padding := 0

	if buf == nil {
		saveComponentType.Fields = make([]Property, 0)
	}

	padding += RoWToNone(data.FromNew(padding), &saveComponentType.Fields, buf, 0)

	return padding
}

func (entityType *EntityType) Process(data util.RawHolder, component *Parsable, buf *bytes.Buffer) int {
	padding := 0

	padding += util.RoWInt32StringNull(data.From(padding), &entityType.ParentObjectRoot, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &entityType.ParentObjectName, buf)
	padding += 4

	var componentCount = int32(len(entityType.Components))
	util.RoWInt32(data.From(padding), &componentCount, buf)
	padding += 4

	if buf == nil {
		entityType.Components = make([]ObjectProperty, componentCount)
	}

	for i := 0; i < int(componentCount); i++ {
		padding += util.RoWInt32StringNull(data.From(padding), &entityType.Components[i].World, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &entityType.Components[i].Class, buf)
		padding += 4
	}

	if buf == nil {
		entityType.Fields = make([]Property, 0)
	}

	padding += RoWToNone(data.FromNew(padding), &entityType.Fields, buf, 0)

	if (buf == nil && data.Length()-padding > 4) || (buf != nil && entityType.ExtraObjects != nil) {
		if specialFunc, ok := specialProcessorClasses[entityType.ClassType]; ok {
			padded := specialFunc(data.FromNew(padding), &entityType.ExtraObjects, buf)
			padding += padded

			if entityType.ExtraObjects == nil {
				logrus.Errorf("%v Did not process any data [%d]\n", entityType.ClassType, data.Length()-padding)
			} else if data.Length()-padding > 4 {
				logrus.Errorf("%v Did not read to end [%d]\n", entityType.ClassType, data.Length()-padding)
			}
		} else {
			logrus.Errorf("%v has >4 bytes [%d] left and is not handled as a special case!\n", entityType.ClassType, data.Length()-padding)
		}
	}

	util.RoWBytes(data.From(padding), &entityType.Extra, buf)
	padding += len(data.From(padding))

	return padding
}

func RoWToNone(data util.RawHolder, target *[]Property, buf *bytes.Buffer, depth int) int {
	if target != nil {
		if buf == nil && !data.IsNil() {
			padding := 0
			values := make([]Property, 0)

			name := ""
			for name != "None" && data.Length()-padding > 4 {
				var tempValue Property
				padded := ProcessProperty(data.FromNew(padding), &tempValue, buf, depth)
				name = tempValue.Name
				padding += padded

				if name != "None" {
					values = append(values, tempValue)
				}
			}

			if depth == 0 {
				// Skip 4 null bytes?
				padding += 4
			}

			*target = values

			return padding
		} else if buf != nil && data.IsNil() {
			padding := 0

			for _, property := range *target {
				padding += ProcessProperty(data.FromNew(padding), &property, buf, depth)
			}

			padding += util.WriteInt32StringNull("None", buf)
			padding += 4

			if depth == 0 {
				buf.Write([]byte{0x00, 0x00, 0x00, 0x00})
				padding += 4
			}

			return padding
		}
	}

	panic("Invalid State!")
}

func (wrapper *EntityType) UnmarshalJSON(b []byte) error {
	var temp map[string]jsoniter.RawMessage
	err := jsoniter.Unmarshal(b, &temp)

	if err != nil {
		return err
	}

	_ = jsoniter.Unmarshal(temp["class_type"], &wrapper.ClassType)
	_ = jsoniter.Unmarshal(temp["entity_type"], &wrapper.EntityType)
	_ = jsoniter.Unmarshal(temp["instance_type"], &wrapper.InstanceType)
	_ = jsoniter.Unmarshal(temp["magic_int_1"], &wrapper.MagicInt1)
	_ = jsoniter.Unmarshal(temp["magic_int_2"], &wrapper.MagicInt2)
	_ = jsoniter.Unmarshal(temp["rotation"], &wrapper.Rotation)
	_ = jsoniter.Unmarshal(temp["position"], &wrapper.Position)
	_ = jsoniter.Unmarshal(temp["scale"], &wrapper.Scale)
	_ = jsoniter.Unmarshal(temp["parent_object_root"], &wrapper.ParentObjectRoot)
	_ = jsoniter.Unmarshal(temp["parent_object_name"], &wrapper.ParentObjectName)
	_ = jsoniter.Unmarshal(temp["components"], &wrapper.Components)
	_ = jsoniter.Unmarshal(temp["fields"], &wrapper.Fields)

	if _, ok := temp["extra"]; ok {
		_ = jsoniter.Unmarshal(temp["extra"], &wrapper.Extra)
	}

	if _, ok := temp["extra_objects"]; ok {
		switch wrapper.ClassType {
		case "/Game/FactoryGame/-Shared/Blueprint/BP_CircuitSubsystem.BP_CircuitSubsystem_C":
			data := BP_CircuitSubsystem_C{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/-Shared/Blueprint/BP_GameMode.BP_GameMode_C":
			data := BP_GameMode_C{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/-Shared/Blueprint/BP_GameState.BP_GameState_C":
			data := BP_GameState_C{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/-Shared/Blueprint/BP_RailroadSubsystem.BP_RailroadSubsystem_C":
			data := BP_RailroadSubsystem_C{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/Buildable/Factory/PowerLine/Build_PowerLine.Build_PowerLine_C":
			data := Build_PowerLine_C{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/Character/Player/BP_PlayerState.BP_PlayerState_C":
			data := BP_PlayerState_C{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/Buildable/Vehicle/Train/Wagon/BP_FreightWagon.BP_FreightWagon_C":
			data := BP_FreightWagon_C{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/Buildable/Vehicle/Train/Locomotive/BP_Locomotive.BP_Locomotive_C":
			data := BP_Locomotive_C{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/Buildable/Vehicle/Tractor/BP_Tractor.BP_Tractor_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Vehicle/Truck/BP_Truck.BP_Truck_C":
			data := BP_Vehicle{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk1/Build_ConveyorBeltMk1.Build_ConveyorBeltMk1_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk2/Build_ConveyorBeltMk2.Build_ConveyorBeltMk2_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk3/Build_ConveyorBeltMk3.Build_ConveyorBeltMk3_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk4/Build_ConveyorBeltMk4.Build_ConveyorBeltMk4_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk5/Build_ConveyorBeltMk5.Build_ConveyorBeltMk5_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk6/Build_ConveyorBeltMk6.Build_ConveyorBeltMk6_C":
			data := BP_Belt{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		case "/Game/FactoryGame/-Shared/Blueprint/BP_StorySubsystem.BP_StorySubsystem_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Building/Foundation/Build_Foundation_8x4_01.Build_Foundation_8x4_01_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/ConveyorPole/Build_ConveyorPole.Build_ConveyorPole_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/MinerMK1/Build_MinerMk1.Build_MinerMk1_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/PowerPoleMk1/Build_PowerPoleMk1.Build_PowerPoleMk1_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/SmelterMk1/Build_SmelterMk1.Build_SmelterMk1_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/StorageContainerMk1/Build_StorageContainerMk1.Build_StorageContainerMk1_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/TradingPost/Build_TradingPost.Build_TradingPost_C":
			fallthrough
		case "/Game/FactoryGame/Buildable/Factory/Workshop/Build_Workshop.Build_Workshop_C":
			fallthrough
		case "/Game/FactoryGame/Character/Creature/BP_CreatureSpawner.BP_CreatureSpawner_C":
			fallthrough
		case "/Game/FactoryGame/Recipes/Research/BP_ResearchManager.BP_ResearchManager_C":
			fallthrough
		case "/Script/FactoryGame.FGFoliageRemoval":
			data := BP_Generic{}
			err = jsoniter.Unmarshal(temp["extra_objects"], &data)
			wrapper.ExtraObjects = data
			break
		}
	}

	return nil
}
