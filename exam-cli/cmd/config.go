package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func readConfig() map[string]string {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return nil
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	return result
}
