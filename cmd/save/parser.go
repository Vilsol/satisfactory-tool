package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"satisfactory-tool/save"
)

func main() {
	flag.Parse()
	satisfactorySave := save.ParseSave(flag.Arg(0))

	bytes, err := json.Marshal(satisfactorySave)
	fmt.Println(err)
	_ = ioutil.WriteFile("output.json", bytes, 0666)
}
