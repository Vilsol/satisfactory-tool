package save

import (
	"bytes"
	"satisfactory-tool/util"
)

var specialClasses = map[string]func([]byte) (interface{}, int){
	"/Game/FactoryGame/-Shared/Blueprint/BP_CircuitSubsystem.BP_CircuitSubsystem_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magicSystem := data[padding : padding+4]
		padding += 4

		circuitCount := int(util.Int32(data[padding:]))
		padding += 4

		circuits := make([]BP_CircuitSubsystem_C_Circuit, circuitCount)

		for i := 0; i < circuitCount; i++ {
			// TODO Unknown
			magicCircuit := data[padding : padding+4]
			padding += 4

			world, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			entity, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			circuits[i] = BP_CircuitSubsystem_C_Circuit{
				Magic:  magicCircuit,
				World:  world,
				Entity: entity,
			}
		}

		return BP_CircuitSubsystem_C{
			Magic:    magicSystem,
			Circuits: circuits,
		}, padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_GameMode.BP_GameMode_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magic := data[padding : padding+4]
		padding += 4

		objectCount := int(util.Int32(data[padding:]))
		padding += 4

		objects := make([]ObjectProperty, objectCount)

		for i := 0; i < objectCount; i++ {
			world, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			class, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			objects[i] = ObjectProperty{
				World: world,
				Class: class,
			}
		}

		return BP_GameMode_C{
			Magic:   magic,
			Objects: objects,
		}, padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_GameState.BP_GameState_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magic := data[padding : padding+4]
		padding += 4

		objectCount := int(util.Int32(data[padding:]))
		padding += 4

		objects := make([]ObjectProperty, objectCount)

		for i := 0; i < objectCount; i++ {
			world, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			class, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			objects[i] = ObjectProperty{
				World: world,
				Class: class,
			}
		}

		return BP_GameState_C{
			Magic:   magic,
			Objects: objects,
		}, padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_RailroadSubsystem.BP_RailroadSubsystem_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magicRailroad := data[padding : padding+4]
		padding += 4

		trainCount := int(util.Int32(data[padding:]))
		padding += 4

		trains := make([]BP_RailroadSubsystem_C_Train, trainCount)

		for i := 0; i < trainCount; i++ {
			// TODO Unknown
			magicTrain := data[padding : padding+4]
			padding += 4

			world, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			entity, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			worldSecond, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			entitySecond, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			worldTimetable, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			entityTimetable, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			trains[i] = BP_RailroadSubsystem_C_Train{
				Magic:           magicTrain,
				World:           world,
				Entity:          entity,
				WorldSecond:     worldSecond,
				EntitySecond:    entitySecond,
				WorldTimetable:  worldTimetable,
				EntityTimetable: entityTimetable,
			}
		}

		return BP_RailroadSubsystem_C{
			Magic:  magicRailroad,
			Trains: trains,
		}, padding
	},
	"/Game/FactoryGame/Buildable/Factory/PowerLine/Build_PowerLine.Build_PowerLine_C": func(data []byte) (interface{}, int) {
		padding := 0

		values, padded := ReadToNone(data[padding:], 0)
		padding += padded

		sourceWorld, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		sourceEntity, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		targetWorld, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		targetEntity, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		return Build_PowerLine_C{
			Values:       values,
			SourceWorld:  sourceWorld,
			SourceEntity: sourceEntity,
			TargetWorld:  targetWorld,
			TargetEntity: targetEntity,
		}, padding
	},
	"/Game/FactoryGame/Character/Player/BP_PlayerState.BP_PlayerState_C": func(data []byte) (interface{}, int) {
		padding := 0

		values, padded := ReReadToZero(data, 0)
		padding += padded

		magic := data[padding:]
		padding += len(data[padding:])

		return BP_PlayerState_C{
			Values: values,
			Magic:  magic,
		}, padding
	},
	"/Game/FactoryGame/Buildable/Vehicle/Train/Wagon/BP_FreightWagon.BP_FreightWagon_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magic := data[padding : padding+8]
		padding += 8

		beforeWorld, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		beforeEntity, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		frontWorld, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		frontEntity, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		return BP_FreightWagon_C{
			Magic:        magic,
			BeforeWorld:  beforeWorld,
			BeforeEntity: beforeEntity,
			FrontWorld:   frontWorld,
			FrontEntity:  frontEntity,
		}, padding
	},
	"/Game/FactoryGame/Buildable/Vehicle/Train/Locomotive/BP_Locomotive.BP_Locomotive_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magic := data[padding : padding+8]
		padding += 8

		// TODO Might be flipped

		beforeWorld, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		beforeEntity, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		frontWorld, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		frontEntity, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		return BP_Locomotive_C{
			Magic:        magic,
			BeforeWorld:  beforeWorld,
			BeforeEntity: beforeEntity,
			FrontWorld:   frontWorld,
			FrontEntity:  frontEntity,
		}, padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_StorySubsystem.BP_StorySubsystem_C":                                     ReadGeneric,
	"/Game/FactoryGame/Buildable/Building/Foundation/Build_Foundation_8x4_01.Build_Foundation_8x4_01_C":             ReadGeneric,
	"/Game/FactoryGame/Buildable/Factory/ConveyorPole/Build_ConveyorPole.Build_ConveyorPole_C":                      ReadGeneric,
	"/Game/FactoryGame/Buildable/Factory/MinerMK1/Build_MinerMk1.Build_MinerMk1_C":                                  ReadGeneric,
	"/Game/FactoryGame/Buildable/Factory/PowerPoleMk1/Build_PowerPoleMk1.Build_PowerPoleMk1_C":                      ReadGeneric,
	"/Game/FactoryGame/Buildable/Factory/SmelterMk1/Build_SmelterMk1.Build_SmelterMk1_C":                            ReadGeneric,
	"/Game/FactoryGame/Buildable/Factory/StorageContainerMk1/Build_StorageContainerMk1.Build_StorageContainerMk1_C": ReadGeneric,
	"/Game/FactoryGame/Buildable/Factory/TradingPost/Build_TradingPost.Build_TradingPost_C":                         ReadGeneric,
	"/Game/FactoryGame/Buildable/Factory/Workshop/Build_Workshop.Build_Workshop_C":                                  ReadGeneric,
	"/Game/FactoryGame/Character/Creature/BP_CreatureSpawner.BP_CreatureSpawner_C":                                  ReadGeneric,
	"/Game/FactoryGame/Recipes/Research/BP_ResearchManager.BP_ResearchManager_C":                                    ReadGeneric,
	"/Script/FactoryGame.FGFoliageRemoval":                                                                          ReadGeneric,
	"/Game/FactoryGame/Buildable/Vehicle/Tractor/BP_Tractor.BP_Tractor_C":                                           ReadVehicle,
	"/Game/FactoryGame/Buildable/Vehicle/Truck/BP_Truck.BP_Truck_C":                                                 ReadVehicle,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk1/Build_ConveyorBeltMk1.Build_ConveyorBeltMk1_C":             ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk2/Build_ConveyorBeltMk2.Build_ConveyorBeltMk2_C":             ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk3/Build_ConveyorBeltMk3.Build_ConveyorBeltMk3_C":             ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk4/Build_ConveyorBeltMk4.Build_ConveyorBeltMk4_C":             ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk5/Build_ConveyorBeltMk5.Build_ConveyorBeltMk5_C":             ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk6/Build_ConveyorBeltMk6.Build_ConveyorBeltMk6_C":             ReadBelt,
}

func ReReadToZero(data []byte, depth int) ([][]Property, int) {
	padding := 0
	values := make([][]Property, 0)

	for len(data)-padding > 4 && util.Int32(data[padding:]) > 0 {
		value, padded := ReadToNone(data[padding:], depth+1)
		padding += padded
		values = append(values, value)
	}

	return values, padding
}

func ReadToNone(data []byte, depth int) ([]Property, int) {
	padding := 0
	values := make([]Property, 0)

	name := ""
	for name != "None" && len(data)-padding > 4 {
		property, padded := ParseProperty(data[padding:], depth+1)
		name = property.Name
		padding += padded

		if name != "None" {
			values = append(values, property)
		}
	}

	if depth == 0 {
		// Skip 4 null bytes?
		padding += 4
	}

	return values, padding
}

func ReadGeneric(data []byte) (interface{}, int) {
	values, padding := ReReadToZero(data, 0)
	return BP_Generic{
		Values: values,
	}, padding
}

func ReadBelt(data []byte) (interface{}, int) {
	padding := 0

	values, padded := ReReadToZero(data[padding:], 0)
	padding += padded

	// TODO Unknown
	magicBelt := data[padding : padding+4]
	padding += 4

	itemCount := int(util.Int32(data[padding:]))
	padding += 4

	items := make([]BP_Belt_Item, itemCount)

	for i := 0; i < itemCount; i++ {
		itemMagic1 := data[padding : padding+4]
		padding += 4

		itemName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		itemMagic2 := data[padding : padding+12]
		padding += 12

		items[i] = BP_Belt_Item{
			Magic1: itemMagic1,
			Name:   itemName,
			Magic2: itemMagic2,
		}
	}

	return BP_Belt{
		Values: values,
		Magic:  magicBelt,
		Items:  items,
	}, padding
}

func ReadVehicle(data []byte) (interface{}, int) {
	padding := 0

	// TODO Unknown
	magicOuter := data[padding : padding+4]
	padding += 4

	objectCount := int(util.Int32(data[padding:]))
	padding += 4

	objects := make([]BP_Vehicle_Object, objectCount)

	for i := 0; i < objectCount; i++ {
		name, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		// TODO Unknown
		magicInner := data[padding : padding+53]
		padding += 53

		objects[i] = BP_Vehicle_Object{
			Name:  name,
			Magic: magicInner,
		}
	}

	return BP_Vehicle{
		Magic:   magicOuter,
		Objects: objects,
	}, padding
}

var specialProcessorClasses = map[string]func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int{
	"/Game/FactoryGame/-Shared/Blueprint/BP_CircuitSubsystem.BP_CircuitSubsystem_C": func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
		var targetStruct BP_CircuitSubsystem_C

		if buf == nil {
			targetStruct = BP_CircuitSubsystem_C{}
			*target.(*interface{}) = &targetStruct
		} else {
			targetStruct = (*target.(*interface{})).(BP_CircuitSubsystem_C)
		}

		padding := 0

		/*
			// TODO Unknown
			util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Magic, buf)
			padding += 4
		*/

		var circuitCount = int32(len(targetStruct.Circuits))
		util.RoWInt32(data.From(padding), &circuitCount, buf)
		padding += 4

		if buf == nil {
			targetStruct.Circuits = make([]BP_CircuitSubsystem_C_Circuit, circuitCount)
		}

		for i := 0; i < int(circuitCount); i++ {
			// TODO Unknown
			util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Circuits[i].Magic, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Circuits[i].World, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Circuits[i].Entity, buf)
			padding += 4
		}

		return padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_GameMode.BP_GameMode_C": func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
		var targetStruct BP_GameMode_C

		if buf == nil {
			targetStruct = BP_GameMode_C{}
			*target.(*interface{}) = &targetStruct
		} else {
			targetStruct = (*target.(*interface{})).(BP_GameMode_C)
		}

		padding := 0

		/*
			// TODO Unknown
			util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Magic, buf)
			padding += 4
		*/

		var objectCount = int32(len(targetStruct.Objects))
		util.RoWInt32(data.From(padding), &objectCount, buf)
		padding += 4

		if buf == nil {
			targetStruct.Objects = make([]ObjectProperty, objectCount)
		}

		for i := 0; i < int(objectCount); i++ {
			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Objects[i].World, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Objects[i].Class, buf)
			padding += 4
		}

		return padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_GameState.BP_GameState_C": func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
		var targetStruct BP_GameState_C

		if buf == nil {
			targetStruct = BP_GameState_C{}
			*target.(*interface{}) = &targetStruct
		} else {
			targetStruct = (*target.(*interface{})).(BP_GameState_C)
		}

		padding := 0

		/*
			// TODO Unknown
			util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Magic, buf)
			padding += 4
		*/

		var objectCount = int32(len(targetStruct.Objects))
		util.RoWInt32(data.From(padding), &objectCount, buf)
		padding += 4

		if buf == nil {
			targetStruct.Objects = make([]ObjectProperty, objectCount)
		}

		for i := 0; i < int(objectCount); i++ {
			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Objects[i].World, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Objects[i].Class, buf)
			padding += 4
		}

		return padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_RailroadSubsystem.BP_RailroadSubsystem_C": func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
		var targetStruct BP_RailroadSubsystem_C

		if buf == nil {
			targetStruct = BP_RailroadSubsystem_C{}
			*target.(*interface{}) = &targetStruct
		} else {
			targetStruct = (*target.(*interface{})).(BP_RailroadSubsystem_C)
		}

		padding := 0

		/*
			// TODO Unknown
			util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Magic, buf)
			padding += 4
		*/

		var trainCount = int32(len(targetStruct.Trains))
		util.RoWInt32(data.From(padding), &trainCount, buf)
		padding += 4

		if buf == nil {
			targetStruct.Trains = make([]BP_RailroadSubsystem_C_Train, trainCount)
		}

		for i := 0; i < int(trainCount); i++ {
			// TODO Unknown
			util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Trains[i].Magic, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Trains[i].World, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Trains[i].Entity, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Trains[i].WorldSecond, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Trains[i].EntitySecond, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Trains[i].WorldTimetable, buf)
			padding += 4

			padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Trains[i].EntityTimetable, buf)
			padding += 4
		}

		return padding
	},
	"/Game/FactoryGame/Buildable/Factory/PowerLine/Build_PowerLine.Build_PowerLine_C": func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
		var targetStruct Build_PowerLine_C

		if buf == nil {
			targetStruct = Build_PowerLine_C{}
			*target.(*interface{}) = &targetStruct
		} else {
			targetStruct = (*target.(*interface{})).(Build_PowerLine_C)
		}

		padding := 0

		// padding += RoWToNone(data.FromNew(padding), &targetStruct.Values, buf, 0)

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.SourceWorld, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.SourceEntity, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.TargetWorld, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.TargetEntity, buf)
		padding += 4

		return padding
	},
	"/Game/FactoryGame/Character/Player/BP_PlayerState.BP_PlayerState_C": func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
		var targetStruct BP_PlayerState_C

		if buf == nil {
			targetStruct = BP_PlayerState_C{}
			*target.(*interface{}) = &targetStruct
		} else {
			targetStruct = (*target.(*interface{})).(BP_PlayerState_C)
		}

		padding := 0

		/*
			TODO
			values, padded := ReReadToZero(data, 0)
			padding += padded
		*/

		// TODO Merge reading/writing
		if buf == nil {
			// TODO Unknown
			util.RoWBytes(data.From(padding), &targetStruct.Magic, buf)
			padding += len(data.From(padding))
		} else {
			tempArray := targetStruct.Magic
			util.RoWBytes(data.From(padding), &tempArray, buf)
			padding += len(data.From(padding))
		}

		return padding
	},
	"/Game/FactoryGame/Buildable/Vehicle/Train/Wagon/BP_FreightWagon.BP_FreightWagon_C": func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
		var targetStruct BP_FreightWagon_C

		if buf == nil {
			targetStruct = BP_FreightWagon_C{}
			*target.(*interface{}) = &targetStruct
		} else {
			targetStruct = (*target.(*interface{})).(BP_FreightWagon_C)
		}

		padding := 0

		/*
			// TODO Unknown
			magic := data[padding : padding+8]
			padding += 8
		*/

		// TODO Unknown
		if buf == nil {
			util.RoWBytes(data.FromTo(padding+4, padding+8), &targetStruct.Magic, buf)
			padding += 4
		} else {
			tempArray := targetStruct.Magic[4:]
			util.RoWBytes(data.From(padding+4), &tempArray, buf)
			padding += len(data.From(padding + 4))
		}

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.BeforeWorld, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.BeforeEntity, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.FrontWorld, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.FrontEntity, buf)
		padding += 4

		return padding
	},
	"/Game/FactoryGame/Buildable/Vehicle/Train/Locomotive/BP_Locomotive.BP_Locomotive_C": func(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
		var targetStruct BP_Locomotive_C

		if buf == nil {
			targetStruct = BP_Locomotive_C{}
			*target.(*interface{}) = &targetStruct
		} else {
			targetStruct = (*target.(*interface{})).(BP_Locomotive_C)
		}

		padding := 0

		/*
			// TODO Unknown
			magic := data[padding : padding+8]
			padding += 8
		*/

		// TODO Unknown
		if buf == nil {
			util.RoWBytes(data.FromTo(padding+4, padding+8), &targetStruct.Magic, buf)
			padding += 4
		} else {
			tempArray := targetStruct.Magic[4:]
			util.RoWBytes(data.From(padding+4), &tempArray, buf)
			padding += len(data.From(padding + 4))
		}

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.BeforeWorld, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.BeforeEntity, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.FrontWorld, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.FrontEntity, buf)
		padding += 4

		return padding
	},
	"/Game/FactoryGame/Buildable/Vehicle/Tractor/BP_Tractor.BP_Tractor_C":                               RoWVehicle,
	"/Game/FactoryGame/Buildable/Vehicle/Truck/BP_Truck.BP_Truck_C":                                     RoWVehicle,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk1/Build_ConveyorBeltMk1.Build_ConveyorBeltMk1_C": RoWBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk2/Build_ConveyorBeltMk2.Build_ConveyorBeltMk2_C": RoWBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk3/Build_ConveyorBeltMk3.Build_ConveyorBeltMk3_C": RoWBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk4/Build_ConveyorBeltMk4.Build_ConveyorBeltMk4_C": RoWBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk5/Build_ConveyorBeltMk5.Build_ConveyorBeltMk5_C": RoWBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk6/Build_ConveyorBeltMk6.Build_ConveyorBeltMk6_C": RoWBelt,
}

func RoWBelt(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
	var targetStruct BP_Belt

	if buf == nil {
		targetStruct = BP_Belt{}
		*target.(*interface{}) = &targetStruct
	} else {
		targetStruct = (*target.(*interface{})).(BP_Belt)
	}

	padding := 0

	/*
		TODO
		values, padded := ReReadToZero(data[padding:], 0)
		padding += padded

		// TODO Unknown
		magicBelt := data[padding : padding+4]
		padding += 4
	*/

	var itemCount = int32(len(targetStruct.Items))
	util.RoWInt32(data.From(padding), &itemCount, buf)
	padding += 4

	if buf == nil {
		targetStruct.Items = make([]BP_Belt_Item, itemCount)
	}

	for i := 0; i < int(itemCount); i++ {
		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Items[i].Magic1, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Items[i].Name, buf)
		padding += 4

		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+12), &targetStruct.Items[i].Magic2, buf)
		padding += 12
	}

	return padding
}

func RoWVehicle(data util.RawHolder, target interface{}, buf *bytes.Buffer) int {
	var targetStruct BP_Vehicle

	if buf == nil {
		targetStruct = BP_Vehicle{}
		*target.(*interface{}) = &targetStruct
	} else {
		targetStruct = (*target.(*interface{})).(BP_Vehicle)
	}

	padding := 0

	/*
		// TODO Unknown
		magicOuter := data[padding : padding+4]
		padding += 4
	*/

	var objectCount = int32(len(targetStruct.Objects))
	util.RoWInt32(data.From(padding), &objectCount, buf)
	padding += 4

	if buf == nil {
		targetStruct.Objects = make([]BP_Vehicle_Object, objectCount)
	}

	for i := 0; i < int(objectCount); i++ {
		padding += util.RoWInt32StringNull(data.From(padding), &targetStruct.Objects[i].Name, buf)
		padding += 4

		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+53), &targetStruct.Objects[i].Magic, buf)
		padding += 53
	}

	return padding
}
