package save

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"satisfactory-tool/util"
	"strconv"
)

type SatisfactorySave struct {
	SaveHeaderVersion int
	SaveVersion       int
	BuildVersion      int
	LevelType         string
	LevelOptions      string
	SessionName       string
	PlayTimeSeconds   int
	SaveDate          int
	SessionVisibility byte
	WorldData         []Parsable
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
		logrus.Info("Obj: ", i, " Pos: ", padding, " Len: ", length)

		// fmt.Println(length)

		worldData[i].Parse(length, saveData[padding:padding+length])
		padding += length
		// fmt.Println(padding)
	}

	if len(saveData)-padding > 4 {
		logrus.Errorf("Extra at the end of the file: %5d: %#v, %#v\n", len(saveData)-padding, saveData[padding:padding+4], string(saveData[padding+4:padding+4+(len(saveData)-padding)]))
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
	}
}
