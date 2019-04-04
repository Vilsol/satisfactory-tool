package satisfactory_tool

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"satisfactory-tool/save"
	"testing"
)

func TestSav2Json2Sav(t *testing.T) {
	satisfactorySave := save.ParseSaveNew("test.sav")

	logrus.Infof("Generating JSON output\n")

	var marshaled []byte
	var err error

	marshaled, err = jsoniter.Marshal(satisfactorySave)

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

	bytes, _ := ioutil.ReadFile("test.sav")

	if len(result) != len(bytes) {
		logrus.Error("File sizes don't match! Original: ", len(bytes), " vs ", len(result))
	}

	minSize := len(result)

	if len(bytes) < minSize {
		minSize = len(bytes)
	}

	for i := 0; i < minSize; i++ {
		if result[i] != bytes[i] {
			logrus.Panic(fmt.Sprintf("Mismatch at: %#x!", i))
		}
	}

	logrus.Info("Test Pass!")
}
