package save

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"satisfactory-tool/util"
	"strconv"
)

type SatisfactorySave struct {
	SaveHeaderVersion int              `json:"save_header_version"`
	SaveVersion       int              `json:"save_version"`
	BuildVersion      int              `json:"build_version"`
	LevelType         string           `json:"level_type"`
	LevelOptions      string           `json:"level_options"`
	SessionName       string           `json:"session_name"`
	PlayTimeSeconds   int              `json:"play_time_seconds"`
	SaveDate          int              `json:"save_date"`
	SessionVisibility byte             `json:"session_visibility"`
	WorldData         []Parsable       `json:"world_data"`
	ExtraObjects      []ObjectProperty `json:"extra_objects"`
	Extra             *[]byte          `json:"extra,omitempty"`
}

const SaveComponentTypeID = 0x0
const EntityTypeID = 0x1

func ParseSave(path string) *SatisfactorySave {
	saveData, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	padding := 0

	saveHeaderVersion := int(util.Int32(saveData[padding:]))
	padding += 4

	saveVersion := int(util.Int32(saveData[padding:]))
	padding += 4

	buildVersion := int(util.Int32(saveData[padding:]))
	padding += 4

	levelType, strLength := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLength

	levelOptions, strLength := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLength

	sessionName, strLength := util.Int32StringNull(saveData[padding:])
	padding += 4 + strLength

	playTimeSeconds := int(util.Int32(saveData[padding:]))
	padding += 4

	saveDate := int(util.Int64(saveData[padding:]))
	padding += 8

	sessionVisibility := saveData[padding]
	padding += 1

	worldDataLength := int(util.Int32(saveData[padding:]))
	padding += 4

	worldData := make([]Parsable, worldDataLength)

	for i := 0; i < worldDataLength; i++ {
		dataType := int(util.Int32(saveData[padding:]))
		padding += 4

		switch dataType {
		case SaveComponentTypeID:
			saveComponentType, padded := ParseSaveComponentType(saveData[padding:])
			padding += padded
			worldData[i] = saveComponentType
			break
		case EntityTypeID:
			entityType, padded := ParseEntityType(saveData[padding:])
			padding += padded
			worldData[i] = entityType
			break
		default:
			panic("Unknown type: " + strconv.Itoa(dataType))
		}
	}

	extraWorldDataLength := int(util.Int32(saveData[padding:]))
	padding += 4

	for i := 0; i < extraWorldDataLength; i++ {
		length := int(util.Int32(saveData[padding:]))
		padding += 4
		logrus.Infof("Obj: %d Pos: %#x Len: %d", i, padding, length)

		// fmt.Println(length)

		worldData[i].Parse(length, saveData[padding:padding+length])
		padding += length
		// fmt.Println(padding)
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
		Extra:             &extraData,
	}
}
