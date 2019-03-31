package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"satisfactory-tool/save"
)

func main() {
	indent := flag.Bool("indent", false, "Whether to indent output JSON")
	flag.Parse()

	if flag.NArg() < 3 {
		PrintUsage()
		return
	}

	// TODO Convert to flag or env param
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableSorting:         false,
		DisableLevelTruncation: false,
		QuoteEmptyFields:       true,
	})

	if flag.Arg(0) == "sav2json" {
		satisfactorySave := save.ParseSave(flag.Arg(1))

		var marshaled []byte
		var err error

		if *indent {
			marshaled, err = json.MarshalIndent(satisfactorySave, "", "  ")
		} else {
			marshaled, err = json.Marshal(satisfactorySave)
		}

		if err != nil {
			logrus.Panic(err)
		}

		err = ioutil.WriteFile(flag.Arg(2), marshaled, 0666)

		if err != nil {
			logrus.Panic(err)
		}
	} else if flag.Arg(0) == "json2sav" {
		file, err := ioutil.ReadFile(flag.Arg(1))

		if err != nil {
			logrus.Panic(err)
		}

		target := save.SatisfactorySave{}
		err = json.Unmarshal(file, &target)

		if err != nil {
			logrus.Panic(err)
		}
	}

}

func PrintUsage() {
	fmt.Println("Usage: satisfactory-tool [action] [input] [output]")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  sav2json: Convert .sav file to .json")
	fmt.Println("  json2sav: Convert .json file to .sav")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}
