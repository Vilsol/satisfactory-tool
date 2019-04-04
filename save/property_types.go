package save

import (
	"errors"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"satisfactory-tool/util"
)

type Property struct {
	Name  string         `json:"name"`
	Type  string         `json:"type"`
	Index int32          `json:"index"`
	Size  int32          `json:"size"`
	Value ReadOrWritable `json:"value"`
}

type ReadOrWritable interface {
}

type ArrayProperty struct {
	Type   string           `json:"type"`
	Values []ReadOrWritable `json:"values"`

	StructName      string `json:"struct_name,omitempty"`
	StructType      string `json:"struct_type,omitempty"`
	StructSize      int32  `json:"struct_size,omitempty"`
	Magic1          []byte `json:"magic_1,omitempty"`
	StructClassType string `json:"struct_class_type,omitempty"`
	Magic2          []byte `json:"magic_2,omitempty"`
}

type StructProperty struct {
	Type  string         `json:"type"`
	Magic []byte         `json:"magic,omitempty"`
	Value ReadOrWritable `json:"value"`
}

type MapEntry struct {
	Key   ReadOrWritable `json:"key"`
	Value []Property     `json:"value"`
}

type MapProperty struct {
	KeyType   string     `json:"key_type"`
	ValueType string     `json:"value_type"`
	Magic     []byte     `json:"magic"`
	Values    []MapEntry `json:"values"`
}

type ObjectProperty struct {
	World string `json:"world"`
	Class string `json:"class"`
}

type TextProperty struct {
	Magic  []byte `json:"magic"`
	String string `json:"string"`
}

type EnumProperty struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type ByteProperty struct {
	EnumType string `json:"enum_type,omitempty"`
	Byte     byte   `json:"byte"`
	EnumName string `json:"enum_name,omitempty"`
}

func (wrapper *Property) UnmarshalJSON(b []byte) error {
	var temp map[string]jsoniter.RawMessage
	err := jsoniter.Unmarshal(b, &temp)

	if err != nil {
		return err
	}

	_ = jsoniter.Unmarshal(temp["name"], &wrapper.Name)
	_ = jsoniter.Unmarshal(temp["type"], &wrapper.Type)
	_ = jsoniter.Unmarshal(temp["index"], &wrapper.Index)
	_ = jsoniter.Unmarshal(temp["size"], &wrapper.Size)

	wrapper.Value, err = UnwrapProperty(wrapper.Type, temp["value"])

	if err != nil {
		return err
	}

	return nil
}

func (wrapper *ArrayProperty) UnmarshalJSON(b []byte) error {
	var temp map[string]jsoniter.RawMessage
	err := jsoniter.Unmarshal(b, &temp)

	if err != nil {
		return err
	}

	_ = jsoniter.Unmarshal(temp["type"], &wrapper.Type)
	_ = jsoniter.Unmarshal(temp["struct_name"], &wrapper.StructName)
	_ = jsoniter.Unmarshal(temp["struct_type"], &wrapper.StructType)
	_ = jsoniter.Unmarshal(temp["struct_size"], &wrapper.StructSize)
	_ = jsoniter.Unmarshal(temp["magic_1"], &wrapper.Magic1)
	_ = jsoniter.Unmarshal(temp["struct_class_type"], &wrapper.StructClassType)
	_ = jsoniter.Unmarshal(temp["magic_2"], &wrapper.Magic2)

	var tempArray []jsoniter.RawMessage
	_ = jsoniter.Unmarshal(temp["values"], &tempArray)

	wrapper.Values = make([]ReadOrWritable, len(tempArray))

	for i := 0; i < len(wrapper.Values); i++ {
		wrapper.Values[i], err = UnwrapProperty(wrapper.Type, tempArray[i])

		if err != nil {
			return err
		}
	}

	return nil
}

func UnwrapProperty(propertyType string, data []byte) (interface{}, error) {
	switch propertyType {
	case "ArrayProperty":
		prop := ArrayProperty{}
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "StructProperty":
		prop := StructProperty{}
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "MapProperty":
		prop := MapProperty{}
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "ObjectProperty":
		prop := ObjectProperty{}
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "TextProperty":
		prop := TextProperty{}
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "EnumProperty":
		prop := EnumProperty{}
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "ByteProperty":
		prop := ByteProperty{}
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "FloatProperty":
		prop := float32(0)
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "IntProperty":
		prop := int32(0)
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "Int8Property":
		prop := int8(0)
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "BoolProperty":
		prop := false
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	case "StrProperty":
		fallthrough
	case "NameProperty":
		prop := ""
		err := jsoniter.Unmarshal(data, &prop)
		return &prop, err
	default:
		return nil, errors.New("Unknown Property Type: " + propertyType)
	}
}

type Color struct {
	R byte `json:"r"`
	G byte `json:"g"`
	B byte `json:"b"`
	A byte `json:"a"`
}

type LinearColor struct {
	R float32 `json:"r"`
	G float32 `json:"g"`
	B float32 `json:"b"`
	A float32 `json:"a"`
}

type Box struct {
	Min   util.Vector3 `json:"min"`
	Max   util.Vector3 `json:"max"`
	Valid byte         `json:"valid"`
}

type InventoryItem struct {
	Magic     string `json:"magic"`
	ItemName  string `json:"item_name"`
	LevelName string `json:"level_name"`
	PathName  string `json:"path_name"`
}

type RailroadTrackPosition struct {
	World      string  `json:"world"`
	EntityType string  `json:"entity_type"`
	Offset     float32 `json:"offset"`
	Forward    float32 `json:"forward"`
}

type GenericStruct struct {
	Values []Property `json:"values"`
}

func (wrapper *StructProperty) UnmarshalJSON(b []byte) error {
	var temp map[string]jsoniter.RawMessage
	err := jsoniter.Unmarshal(b, &temp)

	if err != nil {
		return err
	}

	_ = jsoniter.Unmarshal(temp["type"], &wrapper.Type)
	_ = jsoniter.Unmarshal(temp["magic"], &wrapper.Magic)

	switch wrapper.Type {
	case "Vector":
		fallthrough
	case "Rotator":
		data := util.Vector3{}
		err = jsoniter.Unmarshal(temp["value"], &data)
		wrapper.Value = data
		break
	case "Color":
		data := Color{}
		err = jsoniter.Unmarshal(temp["value"], &data)
		wrapper.Value = data
		break
	case "LinearColor":
		data := LinearColor{}
		err = jsoniter.Unmarshal(temp["value"], &data)
		wrapper.Value = data
		break
	case "Quat":
		data := util.Vector4{}
		err = jsoniter.Unmarshal(temp["value"], &data)
		wrapper.Value = data
		break
	case "Box":
		data := Box{}
		err = jsoniter.Unmarshal(temp["value"], &data)
		wrapper.Value = data
		break
	case "InventoryItem":
		data := InventoryItem{}
		err = jsoniter.Unmarshal(temp["value"], &data)
		wrapper.Value = data
		break
	case "RailroadTrackPosition":
		data := RailroadTrackPosition{}
		err = jsoniter.Unmarshal(temp["value"], &data)
		wrapper.Value = data
		break
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
		data := GenericStruct{}
		err = jsoniter.Unmarshal(temp["value"], &data)
		wrapper.Value = data
		break
	default:
		logrus.Panic("Unknown Type: " + wrapper.Type)
	}

	if err != nil {
		return err
	}

	return nil
}

type BP_CircuitSubsystem_C_Circuit struct {
	Magic  []byte `json:"magic"`
	World  string `json:"world"`
	Entity string `json:"entity"`
}

type BP_CircuitSubsystem_C struct {
	Circuits []BP_CircuitSubsystem_C_Circuit `json:"circuits"`
}

type BP_GameMode_C struct {
	Objects []ObjectProperty `json:"objects"`
}

type BP_GameState_C struct {
	Objects []ObjectProperty `json:"objects"`
}

type BP_RailroadSubsystem_C_Train struct {
	Magic           []byte `json:"magic"`
	World           string `json:"world"`
	Entity          string `json:"entity"`
	WorldSecond     string `json:"world_second"`
	EntitySecond    string `json:"entity_second"`
	WorldTimetable  string `json:"world_timetable"`
	EntityTimetable string `json:"entity_timetable"`
}

type BP_RailroadSubsystem_C struct {
	Trains []BP_RailroadSubsystem_C_Train `json:"trains"`
}

type Build_PowerLine_C struct {
	SourceWorld  string `json:"source_world"`
	SourceEntity string `json:"source_entity"`
	TargetWorld  string `json:"target_world"`
	TargetEntity string `json:"target_entity"`
}

type BP_PlayerState_C struct {
	Magic []byte `json:"magic"`
}

type BP_FreightWagon_C struct {
	Magic        []byte `json:"magic"`
	BeforeWorld  string `json:"before_world"`
	BeforeEntity string `json:"before_entity"`
	FrontWorld   string `json:"front_world"`
	FrontEntity  string `json:"front_entity"`
}

type BP_Locomotive_C struct {
	Magic        []byte `json:"magic"`
	BeforeWorld  string `json:"before_world"`
	BeforeEntity string `json:"before_entity"`
	FrontWorld   string `json:"front_world"`
	FrontEntity  string `json:"front_entity"`
}

type BP_Generic struct {
	Values [][]Property `json:"values"`
}

type BP_Vehicle_Object struct {
	Name  string `json:"name"`
	Magic []byte `json:"magic"`
}

type BP_Vehicle struct {
	Objects []BP_Vehicle_Object `json:"objects"`
}

type BP_Belt_Item struct {
	Magic1 []byte `json:"magic_1"`
	Name   string `json:"name"`
	Magic2 []byte `json:"magic_2"`
}

type BP_Belt struct {
	Values [][]Property   `json:"values"`
	Items  []BP_Belt_Item `json:"items"`
}
