package main

import (
	"flag"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"satisfactory-tool/save"
)

var extractedClasses = map[string]bool{
	"/Game/FactoryGame/Character/Creature/BP_CreatureSpawner.BP_CreatureSpawner_C":                true,
	"/Game/FactoryGame/Prototype/WAT/BP_WAT2.BP_WAT2_C":                                           true,
	"/Game/FactoryGame/Prototype/WAT/BP_WAT1.BP_WAT1_C":                                           true,
	"/Game/FactoryGame/World/Benefit/Mushroom/BP_Shroom_01.BP_Shroom_01_C":                        true,
	"/Game/FactoryGame/World/Benefit/NutBush/BP_NutBush.BP_NutBush_C":                             true,
	"/Game/FactoryGame/World/Benefit/BerryBush/BP_BerryBush.BP_BerryBush_C":                       true,
	"/Game/FactoryGame/World/Benefit/DropPod/BP_DropPod.BP_DropPod_C":                             true,
	"/Game/FactoryGame/Resource/Environment/AnimalParts/BP_CrabEggParts.BP_CrabEggParts_C":        true,
	"/Game/FactoryGame/Resource/BP_ResourceNodeGeyser.BP_ResourceNodeGeyser_C":                    true,
	"/Game/FactoryGame/Resource/BP_ResourceDeposit.BP_ResourceDeposit_C":                          true,
	"/Game/FactoryGame/Resource/Environment/Crystal/BP_Crystal_mk3.BP_Crystal_mk3_C":              true,
	"/Game/FactoryGame/Resource/Environment/Crystal/BP_Crystal_mk2.BP_Crystal_mk2_C":              true,
	"/Game/FactoryGame/Resource/Environment/Crystal/BP_Crystal.BP_Crystal_C":                      true,
	"/Game/FactoryGame/Resource/BP_ResourceNode.BP_ResourceNode_C":                                true,
	"/Game/FactoryGame/Equipment/C4Dispenser/BP_DestructibleSmallRock.BP_DestructibleSmallRock_C": true,
	"/Game/FactoryGame/Equipment/C4Dispenser/BP_DestructibleLargeRock.BP_DestructibleLargeRock_C": true,
}

func main() {
	indent := flag.Bool("indent", false, "Whether to indent output JSON")
	flag.Parse()

	// TODO Convert to flag or env param
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableSorting:         false,
		DisableLevelTruncation: false,
		QuoteEmptyFields:       true,
	})

	if flag.NArg() == 2 {
		if flag.Arg(0) == "test" {
			cmdTest(indent)
			return
		}
	} else if flag.NArg() == 3 {
		if flag.Arg(0) == "sav2json" {
			cmdSav2json(indent)
			return
		} else if flag.Arg(0) == "json2sav" {
			cmdJson2sav()
			return
		} else if flag.Arg(0) == "extract" {
			cmdExtract(indent)
			return
		}
	}

	PrintUsage()
}

func PrintUsage() {
	flag.CommandLine.SetOutput(os.Stdout)
	fmt.Println("Usage: satisfactory-tool [flags] [action] [input] [output]")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  sav2json: Convert .sav file to .json")
	fmt.Println("    sav2json {save} {json}")
	fmt.Println("  json2sav: Convert .json file to .sav")
	fmt.Println("    json2sav {json} {save}")
	fmt.Println("  test:     Test converting .sav to .json and back to .sav and compare")
	fmt.Println("    test {save}")
	fmt.Println("  extract:  Extract map-based data to .json")
	fmt.Println("    extract {save} {json}")
	fmt.Println()
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func cmdTest(indent *bool) {
	satisfactorySave := save.ParseSave(flag.Arg(1))

	logrus.Infof("Generating JSON output\n")

	var marshaled []byte
	var err error

	if *indent {
		marshaled, err = jsoniter.MarshalIndent(satisfactorySave, "", "  ")
	} else {
		marshaled, err = jsoniter.Marshal(satisfactorySave)
	}

	if err != nil {
		logrus.Panic(err)
	}

	logrus.Infof("Parsing JSON output\n")

	target := save.SatisfactorySave{}
	err = jsoniter.Unmarshal(marshaled, &target)

	if err != nil {
		logrus.Panic(err)
	}

	logrus.Infof("Generating SAV file\n")

	result := target.SaveSave()

	logrus.Infof("Comparing SAV file to original\n")

	bytes, _ := ioutil.ReadFile(flag.Arg(1))

	if len(result) != len(bytes) {
		logrus.Panic("File sizes don't match!")
	}

	for k, v := range result {
		if v != bytes[k] {
			logrus.Panic(fmt.Sprintf("Mismatch at: %#x!", k))
		}
	}

	logrus.Info("Test Pass!")
}

func cmdSav2json(indent *bool) {
	satisfactorySave := save.ParseSave(flag.Arg(1))

	logrus.Infof("Generating JSON output\n")

	var marshaled []byte
	var err error

	if *indent {
		marshaled, err = jsoniter.MarshalIndent(satisfactorySave, "", "  ")
	} else {
		marshaled, err = jsoniter.Marshal(satisfactorySave)
	}

	if err != nil {
		logrus.Panic(err)
	}

	logrus.Infof("Saving JSON output to %s\n", flag.Arg(2))

	err = ioutil.WriteFile(flag.Arg(2), marshaled, 0666)

	if err != nil {
		logrus.Panic(err)
	}
}

func cmdJson2sav() {
	logrus.Infof("Loading JSON file: %s\n", flag.Arg(1))

	file, err := ioutil.ReadFile(flag.Arg(1))

	if err != nil {
		logrus.Panic(err)
	}

	logrus.Infof("Parsing JSON file\n")

	target := save.SatisfactorySave{}
	err = jsoniter.Unmarshal(file, &target)

	if err != nil {
		logrus.Panic(err)
	}

	logrus.Infof("Generating SAV file\n")

	result := target.SaveSave()

	logrus.Infof("Saving SAV file: %s\n", flag.Arg(2))

	err = ioutil.WriteFile(flag.Arg(2), result, 0666)

	if err != nil {
		logrus.Panic(err)
	}
}

func cmdExtract(indent *bool) {
	satisfactorySave := save.ParseSave(flag.Arg(1))

	validEntities := make([]interface{}, 0)
	for _, v := range satisfactorySave.WorldData {
		if _, ok := extractedClasses[v.Data.GetClassType()]; ok {
			entityType := v.Data.(*save.EntityType)
			fmt.Printf("%s: %#v\n", entityType.GetInstanceType(), entityType.Fields)
			validEntities = append(validEntities, entityType)
		}
	}

	var marshaled []byte
	var err error

	if *indent {
		marshaled, err = jsoniter.MarshalIndent(validEntities, "", "  ")
	} else {
		marshaled, err = jsoniter.Marshal(validEntities)
	}

	if err != nil {
		logrus.Panic(err)
	}

	logrus.Infof("Saving JSON output to %s\n", flag.Arg(2))

	err = ioutil.WriteFile(flag.Arg(2), marshaled, 0666)

	if err != nil {
		logrus.Panic(err)
	}
}
