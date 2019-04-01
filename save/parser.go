package save

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"satisfactory-tool/util"
	"strconv"
)

type SatisfactorySave struct {
	SaveHeaderVersion int32             `json:"save_header_version"`
	SaveVersion       int32             `json:"save_version"`
	BuildVersion      int32             `json:"build_version"`
	LevelType         string            `json:"level_type"`
	LevelOptions      string            `json:"level_options"`
	SessionName       string            `json:"session_name"`
	PlayTimeSeconds   int32             `json:"play_time_seconds"`
	SaveDate          int64             `json:"save_date"`
	SessionVisibility byte              `json:"session_visibility"`
	WorldData         []ParsableWrapper `json:"world_data"`
	ExtraObjects      []ObjectProperty  `json:"extra_objects"`
	Extra             []byte            `json:"extra,omitempty"`
}

const SaveComponentTypeID = 0x0
const EntityTypeID = 0x1

func (save *SatisfactorySave) SaveSave() []byte {
	b := &bytes.Buffer{}
	Process(util.RawHolder{}, save, b)

	all, _ := ioutil.ReadAll(b)
	fmt.Println(len(all))
	return all
}

func ParseSave(path string) *SatisfactorySave {
	saveData, err := ioutil.ReadFile(path)

	if err != nil {
		logrus.Panic(err)
	}

	padding := 0

	saveHeaderVersion := util.Int32(saveData[padding:])
	padding += 4

	saveVersion := util.Int32(saveData[padding:])
	padding += 4

	buildVersion := util.Int32(saveData[padding:])
	padding += 4

	levelType, strLength := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLength

	levelOptions, strLength := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLength

	sessionName, strLength := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLength

	playTimeSeconds := util.Int32(saveData[padding:])
	padding += 4

	saveDate := util.Int64(saveData[padding:])
	padding += 8

	sessionVisibility := saveData[padding]
	padding += 1

	worldDataLength := int(util.Int32(saveData[padding:]))
	padding += 4

	worldData := make([]ParsableWrapper, worldDataLength)

	for i := 0; i < worldDataLength; i++ {
		dataType := int(util.Int32(saveData[padding:]))
		padding += 4

		switch dataType {
		case SaveComponentTypeID:
			saveComponentType, padded := ParseSaveComponentType(saveData[padding:])
			padding += padded
			worldData[i] = ParsableWrapper{
				Type: "save",
				Data: saveComponentType,
			}
			break
		case EntityTypeID:
			entityType, padded := ParseEntityType(saveData[padding:])
			padding += padded
			worldData[i] = ParsableWrapper{
				Type: "entity",
				Data: entityType,
			}
			break
		default:
			logrus.Panic("Unknown type: " + strconv.Itoa(dataType))
		}
	}

	extraWorldDataLength := int(util.Int32(saveData[padding:]))
	padding += 4

	for i := 0; i < extraWorldDataLength; i++ {
		length := util.Int32(saveData[padding:])
		padding += 4
		logrus.Infof("Obj: %d Pos: %#x Len: %d", i, padding, length)

		worldData[i].Length = length
		worldData[i].Data.Parse(int(length), saveData[padding:padding+int(length)])
		padding += int(length)
	}

	extraObjectCount := int(util.Int32(saveData[padding:]))
	padding += 4

	extraObjects := make([]ObjectProperty, extraObjectCount)

	for i := 0; i < extraObjectCount; i++ {
		world, strLength := util.Int32StringNull(saveData[padding:])
		padding += 4 + strLength

		class, strLength := util.Int32StringNull(saveData[padding:])
		padding += 4 + strLength

		extraObjects[i] = ObjectProperty{
			World: world,
			Class: class,
		}
	}

	var extraData []byte

	if len(saveData)-padding > 0 {
		logrus.Errorf("Extra at the end of the file: %5d\n%s", len(saveData)-padding, util.HexDump(saveData[padding:]))
		extraData = saveData[padding:]
	}

	return &SatisfactorySave{
		SaveHeaderVersion: saveHeaderVersion,
		SaveVersion:       saveVersion,
		BuildVersion:      buildVersion,
		LevelType:         levelType,
		LevelOptions:      levelOptions,
		SessionName:       sessionName,
		PlayTimeSeconds:   playTimeSeconds,
		SaveDate:          saveDate,
		SessionVisibility: sessionVisibility,
		WorldData:         worldData,
		ExtraObjects:      extraObjects,
		Extra:             extraData,
	}
}

func Process(data util.RawHolder, save *SatisfactorySave, buf *bytes.Buffer) {
	padding := 0

	util.RoWInt32(data.From(padding), &save.SaveHeaderVersion, buf)
	padding += 4

	util.RoWInt32(data.From(padding), &save.SaveVersion, buf)
	padding += 4

	util.RoWInt32(data.From(padding), &save.BuildVersion, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &save.LevelType, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &save.LevelOptions, buf)
	padding += 4

	padding += util.RoWInt32StringNull(data.From(padding), &save.SessionName, buf)
	padding += 4

	util.RoWInt32(data.From(padding), &save.PlayTimeSeconds, buf)
	padding += 4

	util.RoWInt64(data.From(padding), &save.SaveDate, buf)
	padding += 4

	util.RoWB(data.At(padding), &save.SessionVisibility, buf)
	padding += 1

	var worldDataLength = int32(len(save.WorldData))
	util.RoWInt32(data.From(padding), &worldDataLength, buf)
	padding += 4

	if buf == nil {
		save.WorldData = make([]ParsableWrapper, worldDataLength)
	}

	for i := 0; i < int(worldDataLength); i++ {
		var dataType int32

		if buf != nil {
			switch save.WorldData[i].Type {
			case "save":
				dataType = SaveComponentTypeID
				break
			case "entity":
				dataType = EntityTypeID
				break
			}
		}

		util.RoWInt32(data.From(padding), &dataType, buf)
		padding += 4

		if buf == nil {
			save.WorldData[i] = ParsableWrapper{}

			switch dataType {
			case SaveComponentTypeID:
				save.WorldData[i].Type = "save"
				save.WorldData[i].Data = &SaveComponentType{}
				break
			case EntityTypeID:
				save.WorldData[i].Type = "entity"
				save.WorldData[i].Data = &EntityType{}
				break
			}
		}

		switch dataType {
		case SaveComponentTypeID:
			padded := ProcessSaveComponentType(data.FromNew(padding), save.WorldData[i].Data.(*SaveComponentType), buf)
			padding += padded
			break
		case EntityTypeID:
			padded := ProcessEntityType(data.FromNew(padding), save.WorldData[i].Data.(*EntityType), buf)
			padding += padded
			break
		default:
			logrus.Panicf("Unknown type: %d\n", dataType)
		}
	}

	var extraWorldDataLength = int32(len(save.WorldData))
	util.RoWInt32(data.From(padding), &extraWorldDataLength, buf)
	padding += 4

	for i := 0; i < int(extraWorldDataLength); i++ {
		util.RoWInt32(data.From(padding), &save.WorldData[i].Length, buf)
		padding += 4

		logrus.Infof("Obj: %d Pos: %#x Len: %d", i, padding, save.WorldData[i].Length)

		save.WorldData[i].Data.Process(data.FromToNew(padding, padding+int(save.WorldData[i].Length)), &save.WorldData[i].Data, buf)

		padding += int(save.WorldData[i].Length)

		if i >= 365 {
			// fmt.Printf("%#v\n", save.WorldData[i].Data)
			// return
		}
	}

	var extraObjectCount = int32(len(save.ExtraObjects))
	util.RoWInt32(data.From(padding), &extraObjectCount, buf)
	padding += 4

	for i := 0; i < int(extraObjectCount); i++ {
		padding += util.RoWInt32StringNull(data.From(padding), &save.ExtraObjects[i].World, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &save.ExtraObjects[i].Class, buf)
		padding += 4
	}

	util.RoWBytes(data.From(padding), &save.Extra, buf)
}
