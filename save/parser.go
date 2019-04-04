package save

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"satisfactory-tool/util"
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

	logrus.Infof("Generated SAV size: %d bytes\n", len(all))

	return all
}

func ParseSave(path string) *SatisfactorySave {
	logrus.Infof("Loading save file: %s\n", path)

	saveData, err := ioutil.ReadFile(path)

	if err != nil {
		logrus.Panic(err)
	}

	save := SatisfactorySave{}

	Process(util.RawHolder{
		Data: saveData,
	}, &save, nil)

	return &save
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
	padding += 8

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
	}

	var extraObjectCount = int32(len(save.ExtraObjects))
	util.RoWInt32(data.From(padding), &extraObjectCount, buf)
	padding += 4

	if buf == nil {
		save.ExtraObjects = make([]ObjectProperty, extraObjectCount)
	}

	for i := 0; i < int(extraObjectCount); i++ {
		padding += util.RoWInt32StringNull(data.From(padding), &save.ExtraObjects[i].World, buf)
		padding += 4

		padding += util.RoWInt32StringNull(data.From(padding), &save.ExtraObjects[i].Class, buf)
		padding += 4
	}

	util.RoWBytes(data.From(padding), &save.Extra, buf)
}
