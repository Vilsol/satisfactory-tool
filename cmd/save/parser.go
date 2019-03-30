package main

import (
	"encoding/json"
	"flag"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"satisfactory-tool/save"
)

func main() {
	flag.Parse()

	// TODO Convert to flag or env param
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableSorting:         false,
		DisableLevelTruncation: false,
		QuoteEmptyFields:       true,
	})

	satisfactorySave := save.ParseSave(flag.Arg(0))

	bytes, err := json.Marshal(satisfactorySave)

	if err != nil {
		logrus.Panic(err)
	}

	err = ioutil.WriteFile("output.json", bytes, 0666)

	if err != nil {
		logrus.Panic(err)
	}
}
