package save

import (
	"fmt"
	"satisfactory-tool/util"
)

var specialClasses = map[string]func([]byte) (interface{}, int){
	"/Game/FactoryGame/-Shared/Blueprint/BP_CircuitSubsystem.BP_CircuitSubsystem_C": func(data []byte) (interface{}, int) {
		fmt.Printf("%s, %#v\n", "BP_CircuitSubsystem_C", string(data))
		// TODO
		return nil, 0
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_GameMode.BP_GameMode_C": func(data []byte) (interface{}, int) {
		fmt.Printf("%s, %#v\n", "BP_GameMode_C", string(data))
		// TODO
		return nil, 0
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_GameState.BP_GameState_C": func(data []byte) (interface{}, int) {
		fmt.Printf("%s, %#v\n", "BP_GameState_C", string(data))
		// TODO
		return nil, 0
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_RailroadSubsystem.BP_RailroadSubsystem_C": func(data []byte) (interface{}, int) {
		fmt.Printf("%s, %#v\n", "BP_RailroadSubsystem_C", string(data))
		// TODO
		return nil, 0
	},
	"/Game/FactoryGame/-Shared/Blueprint/BP_StorySubsystem.BP_StorySubsystem_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Building/Foundation/Build_Foundation_8x4_01.Build_Foundation_8x4_01_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/ConveyorBeltMk1/Build_ConveyorBeltMk1.Build_ConveyorBeltMk1_C": func(data []byte) (interface{}, int) {
		padding := 0

		values1, padded := ReReadToZero(data[padding:], 0)
		padding += padded

		// TODO

		return []interface{}{values1}, padding
	},
	"/Game/FactoryGame/Buildable/Factory/ConveyorPole/Build_ConveyorPole.Build_ConveyorPole_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/MinerMK1/Build_MinerMk1.Build_MinerMk1_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Game/FactoryGame/Buildable/Factory/PowerLine/Build_PowerLine.Build_PowerLine_C": func(data []byte) (interface{}, int) {
		padding := 0

		values1, padded := ReadToNone(data[padding:], 0)
		padding += padded

		// TODO

		return []interface{}{values1}, padding
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

		values1, padded := ReReadToZero(data, 0)
		padding += padded

		// TODO

		return []interface{}{values1}, padding
	},
	"/Game/FactoryGame/Recipes/Research/BP_ResearchManager.BP_ResearchManager_C": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
	"/Script/FactoryGame.FGFoliageRemoval": func(data []byte) (interface{}, int) {
		return ReReadToZero(data, 0)
	},
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
