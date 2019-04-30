package save

import (
	"bytes"
	"satisfactory-tool/util"
)

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

		// TODO Unknown
		util.RoWBytes(data.From(padding), &targetStruct.Magic, buf)
		padding += len(data.From(padding))

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

		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Magic, buf)
		padding += 4

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

		// TODO Unknown
		util.RoWBytes(data.FromTo(padding, padding+4), &targetStruct.Magic, buf)
		padding += 4

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
	"/Game/FactoryGame/Buildable/Vehicle/Explorer/BP_Explorer.BP_Explorer_C":                            RoWVehicle,
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
