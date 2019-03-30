package save

import (
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

		circuits := make([]map[string]interface{}, circuitCount)

		for i := 0; i < circuitCount; i++ {
			// TODO Unknown
			magicCircuit := data[padding : padding+4]
			padding += 4

			world, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			entity, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			circuits[i] = map[string]interface{}{
				"magic":  magicCircuit,
				"world":  world,
				"entity": entity,
			}
		}

		return map[string]interface{}{
			"magic":    magicSystem,
			"circuits": circuits,
		}, padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_GameMode.BP_GameMode_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magicState := data[padding : padding+4]
		padding += 4

		objectCount := int(util.Int32(data[padding:]))
		padding += 4

		objects := make([]map[string]interface{}, objectCount)

		for i := 0; i < objectCount; i++ {
			world, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			entity, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			objects[i] = map[string]interface{}{
				"world":  world,
				"entity": entity,
			}
		}

		return map[string]interface{}{
			"magic":   magicState,
			"objects": objects,
		}, padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_GameState.BP_GameState_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magicState := data[padding : padding+4]
		padding += 4

		objectCount := int(util.Int32(data[padding:]))
		padding += 4

		objects := make([]map[string]interface{}, objectCount)

		for i := 0; i < objectCount; i++ {
			world, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			entity, strLength := util.Int32StringNull(data[padding:])
			padding += 4 + strLength

			objects[i] = map[string]interface{}{
				"world":  world,
				"entity": entity,
			}
		}

		return map[string]interface{}{
			"magic":   magicState,
			"objects": objects,
		}, padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_RailroadSubsystem.BP_RailroadSubsystem_C": func(data []byte) (interface{}, int) {
		padding := 0

		// TODO Unknown
		magicRailroad := data[padding : padding+4]
		padding += 4

		trainCount := int(util.Int32(data[padding:]))
		padding += 4

		trains := make([]map[string]interface{}, trainCount)

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

			trains[i] = map[string]interface{}{
				"magic":           magicTrain,
				"world":           world,
				"entity":          entity,
				"worldSecond":     worldSecond,
				"entitySecond":    entitySecond,
				"worldTimetable":  worldTimetable,
				"entityTimetable": entityTimetable,
			}
		}

		return map[string]interface{}{
			"magic":  magicRailroad,
			"trains": trains,
		}, padding
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_StorySubsystem.BP_StorySubsystem_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Building/Foundation/Build_Foundation_8x4_01.Build_Foundation_8x4_01_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk1/Build_ConveyorBeltMk1.Build_ConveyorBeltMk1_C": ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk2/Build_ConveyorBeltMk2.Build_ConveyorBeltMk2_C": ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk3/Build_ConveyorBeltMk3.Build_ConveyorBeltMk3_C": ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk4/Build_ConveyorBeltMk4.Build_ConveyorBeltMk4_C": ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk5/Build_ConveyorBeltMk5.Build_ConveyorBeltMk5_C": ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk6/Build_ConveyorBeltMk6.Build_ConveyorBeltMk6_C": ReadBelt,
	"/Game/FactoryGame/Buildable/Factory/ConveyorPole/Build_ConveyorPole.Build_ConveyorPole_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/MinerMK1/Build_MinerMk1.Build_MinerMk1_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
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

		return map[string]interface{}{
			"values":       values,
			"sourceWorld":  sourceWorld,
			"sourceEntity": sourceEntity,
			"targetWorld":  targetWorld,
			"targetEntity": targetEntity,
		}, padding
	},
	"/Game/FactoryGame/Buildable/Factory/PowerPoleMk1/Build_PowerPoleMk1.Build_PowerPoleMk1_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/SmelterMk1/Build_SmelterMk1.Build_SmelterMk1_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/StorageContainerMk1/Build_StorageContainerMk1.Build_StorageContainerMk1_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/TradingPost/Build_TradingPost.Build_TradingPost_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/Workshop/Build_Workshop.Build_Workshop_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Character/Creature/BP_CreatureSpawner.BP_CreatureSpawner_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Character/Player/BP_PlayerState.BP_PlayerState_C": func(data []byte) (interface{}, int) {
		padding := 0

		values, padded := ReReadToZero(data, 0)
		padding += padded

		magic := data[padding:]
		padding += len(data[padding:])

		return map[string]interface{}{
			"values": values,
			"magic":  magic,
		}, padding
	},
	"/Game/FactoryGame/Recipes/Research/BP_ResearchManager.BP_ResearchManager_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Script/FactoryGame.FGFoliageRemoval": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
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

		return map[string]interface{}{
			"magic":        magic,
			"beforeWorld":  beforeWorld,
			"beforeEntity": beforeEntity,
			"frontWorld":   frontWorld,
			"frontEntity":  frontEntity,
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

		return map[string]interface{}{
			"magic":        magic,
			"beforeWorld":  beforeWorld,
			"beforeEntity": beforeEntity,
			"frontWorld":   frontWorld,
			"frontEntity":  frontEntity,
		}, padding
	},
	"/Game/FactoryGame/Buildable/Vehicle/Tractor/BP_Tractor.BP_Tractor_C": ReadVehicle,
	"/Game/FactoryGame/Buildable/Vehicle/Truck/BP_Truck.BP_Truck_C":       ReadVehicle,
}

func ReReadToZero(data []byte, depth int) ([][]map[string]interface{}, int) {
	padding := 0
	values := make([][]map[string]interface{}, 0)

	for len(data)-padding > 4 && util.Int32(data[padding:]) > 0 {
		value, padded := ReadToNone(data[padding:], depth+1)
		padding += padded
		values = append(values, value)
	}

	return values, padding
}

func ReadToNone(data []byte, depth int) ([]map[string]interface{}, int) {
	padding := 0
	values := make([]map[string]interface{}, 0)

	name := ""
	for name != "None" && len(data)-padding > 4 {
		propName, typeName, _, value, _, padded := ParseProperty(data[padding:], depth+1)
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

	return values, padding
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

	items := make([]map[string]interface{}, itemCount)

	for i := 0; i < itemCount; i++ {
		itemMagic1 := data[padding : padding+4]
		padding += 4

		itemName, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		itemMagic2 := data[padding : padding+12]
		padding += 12

		items[i] = map[string]interface{}{
			"magic1":   itemMagic1,
			"itemName": itemName,
			"magic2":   itemMagic2,
		}
	}

	return map[string]interface{}{
		"values": values,
		"magic":  magicBelt,
		"items":  items,
	}, padding
}

func ReadVehicle(data []byte) (interface{}, int) {
	padding := 0

	// TODO Unknown
	magicOuter := data[padding : padding+4]
	padding += 4

	objectCount := int(util.Int32(data[padding:]))
	padding += 4

	objects := make([]map[string]interface{}, objectCount)

	for i := 0; i < objectCount; i++ {
		name, strLength := util.Int32StringNull(data[padding:])
		padding += 4 + strLength

		// TODO Unknown
		magicInner := data[padding : padding+53]
		padding += 53

		objects[i] = map[string]interface{}{
			"name":  name,
			"magic": magicInner,
		}
	}

	return map[string]interface{}{
		"magic":   magicOuter,
		"objects": objects,
	}, padding
}
