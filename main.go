package main

import (
	"encoding/json"
	"github.com/donopttorg/dbupdater/updater"
	"io/ioutil"
	"log"
	"os"
)

const envConfigFile = "env.json"

func main() {
	setSystemEnv()
	updater.StartUpdater()
}

func setSystemEnv() {
	if _, err := os.Stat(envConfigFile); !os.IsNotExist(err) {
		file, _  := os.Open(envConfigFile)
		raw, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		data := map[string]string{}
		err = json.Unmarshal(raw, &data)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range data {
			os.Setenv(k, v)
		}
	}
}
