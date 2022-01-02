package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func getStatus() map[string]string {
	ret := map[string]string{
		"Filename": file,
		"Updated":  "--/--/----",
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		return ret
	}
	ret["Updated"] = fileInfo.ModTime().Format("2/1/2006")
	return ret
}

func loadNames() []string {
	var data dataset
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return []string{} // probably no file found
	}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data.Names
}
